package member_service

import (
	"math/rand"
	"net"
	"time"

	"cs425_mp2/config"
	"cs425_mp2/member_service/protocol_buffer"
	"cs425_mp2/util/logger"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

var BytesSent int = 0
var MessageLossRate float64 = -1

func EncodeMembershipServiceMessage(serviceMessage *protocol_buffer.MembershipServiceMessage) ([]byte, error) {
	message, err := proto.Marshal(serviceMessage)

	return message, err
}

func DecodeMembershipServiceMessage(message []byte) (*protocol_buffer.MembershipServiceMessage, error) {
	list := &protocol_buffer.MembershipServiceMessage{}
	err := proto.Unmarshal(message, list)

	return list, err
}

func SendGossip(serviceMessage *protocol_buffer.MembershipServiceMessage, k int, selfID string) error {
	message, err := EncodeMembershipServiceMessage(serviceMessage)
	if err != nil {
		return err
	}

	dests := GetOtherMembershipListIPs(serviceMessage, selfID)

	if k < len(serviceMessage.MemberList) {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(dests), func(i, j int) {
			dests[i], dests[j] = dests[j], dests[i]
		})

		dests = dests[:k]
	}

	return SendAll(dests, message)
}

func SendHeartbeat(fullMessage, selfMessage *protocol_buffer.MembershipServiceMessage, selfID string) error {
	dests := GetOtherMembershipListIPs(fullMessage, selfID)

	message, err := EncodeMembershipServiceMessage(selfMessage)
	if err != nil {
		return err
	}

	return SendAll(dests, message)
}

func HeartbeatGossip(serviceMessage *protocol_buffer.MembershipServiceMessage, k int, selfID string) error {
	return SendGossip(serviceMessage, k, selfID)
}

func HeartbeatAllToAll(serviceMessage *protocol_buffer.MembershipServiceMessage, selfID string) error {
	selfMessage := &protocol_buffer.MembershipServiceMessage{
		MemberList:      make(map[string]*protocol_buffer.Member),
		Strategy:        serviceMessage.Strategy,
		StrategyCounter: serviceMessage.StrategyCounter,
	}

	selfMember := protocol_buffer.Member{
		HeartbeatCounter: serviceMessage.MemberList[selfID].HeartbeatCounter,
		LastSeen:         ptypes.TimestampNow(),
		IsLeaving:        serviceMessage.MemberList[selfID].IsLeaving,
	}

	selfMessage.MemberList[selfID] = &selfMember

	return SendHeartbeat(serviceMessage, selfMessage, selfID)
}

func SendAll(destinations []string, message []byte) error {
	for _, v := range destinations {
		err := Send(v, message)
		if err != nil {
			return err
		}
	}

	return nil
}

func Send(dest string, message []byte) error {
	if len(message) > config.BUFFER_SIZE {
		logger.WarningLogger.Println("Send: message is larger than BUFFER_SIZE")
	}

	rand.NewSource(time.Now().UnixNano())
	if rand.Float64() > MessageLossRate {
		addr, err := net.ResolveUDPAddr("udp", dest+":"+config.MemberServicePort)
		if err != nil {
			return err
		}

		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			return err
		}

		defer conn.Close()

		_, err = conn.Write(message)
		if err != nil {
			return err
		}

		BytesSent += len(message)
	}

	return nil
}

func Listen(port string, callback func(message []byte) error) error {
	port = ":" + port

	addr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	defer conn.Close()
	buffer := make([]byte, config.BUFFER_SIZE)

	for {
		n, err := conn.Read(buffer)

		if err != nil {
			return err
		}

		callback(buffer[0:n])
	}
}

func ChangeStrategy(input string) {
	if input == config.STRAT_GOSSIP {
		if localMessage.Strategy == config.STRAT_GOSSIP {
			logger.PrintError("System strategy is already gossip")
			return
		}

		localMessage.Strategy = config.STRAT_GOSSIP
	} else if input == config.STRAT_ALL {
		if localMessage.Strategy == config.STRAT_ALL {
			logger.PrintError("System strategy is already all-to-all")
			return
		}

		localMessage.Strategy = config.STRAT_ALL
	} else {
		logger.PrintError("Invalid strategy - must be gossip or all")
		return
	}

	localMessage.StrategyCounter++
	logger.PrintInfo("System strategy successfully changed to", localMessage.Strategy)
}

func sendLeaveRequest() {
	isSending = false

	mux.Lock()
	localMessage.MemberList[selfID].IsLeaving = true

	if localMessage.Strategy == config.STRAT_GOSSIP {
		HeartbeatGossip(localMessage, config.GOSSIP_FANOUT, selfID)
	} else {
		HeartbeatAllToAll(localMessage, selfID)
	}
	mux.Unlock()

	selfID = ""
	localMessage = nil
	logger.PrintInfo("Successfully left")
}



func readNewMessage(message []byte) error {
	if !isSending {
		return nil
	}

	remoteMessage, err := DecodeMembershipServiceMessage(message)
	if err != nil {
		return err
	}

	mux.Lock()

	if isJoining && remoteMessage.Type == protocol_buffer.MessageType_JOINREP {
		isJoining = false
		localMessage.Type = protocol_buffer.MessageType_STANDARD
	}

	if !isMaster && remoteMessage.Type == protocol_buffer.MessageType_JOINREQ {
		mux.Unlock()
		return nil
	}

	mergeMembershipLists(localMessage, remoteMessage, failureList)

	if isMaster && remoteMessage.Type == protocol_buffer.MessageType_JOINREQ {
		logger.PrintInfo("Received join request")
		localMessage.Type = protocol_buffer.MessageType_JOINREP
		message, err := EncodeMembershipServiceMessage(localMessage)
		localMessage.Type = protocol_buffer.MessageType_STANDARD

		if err != nil {
			return err
		}

		dests := GetOtherMembershipListIPs(remoteMessage, selfID)
		logger.PrintWarning("Sending", dests[0], " ", message)
		Send(dests[0], message)
	}

	mux.Unlock()

	return nil
}

func startHeartbeat() {
	for isSending {
		mux.Lock()

		localMessage.MemberList[selfID].LastSeen = ptypes.TimestampNow()
		localMessage.MemberList[selfID].HeartbeatCounter++
		CheckAndRemoveMembershipListFailures(localMessage, &failureList)
		logger.InfoLogger.Println("Current memlist:\n", GetMembershipListString(localMessage, failureList), "\n")

		if isJoining {
			message, _ := EncodeMembershipServiceMessage(localMessage)
			Send(masterIP, message)
		} else {
			if localMessage.Strategy == config.STRAT_GOSSIP {
				HeartbeatGossip(localMessage, config.GOSSIP_FANOUT, selfID)
			} else {
				HeartbeatAllToAll(localMessage, selfID)
			}

			for machineID := range localMessage.MemberList {
				if localMessage.MemberList[machineID].IsLeaving && !failureList[machineID] {
					logger.PrintInfo("Received leave request from machine", machineID)
					failureList[machineID] = true
				}
			}
		}

		mux.Unlock()

		time.Sleep(config.PULSE_TIME * time.Millisecond)
	}
}


package detector

import (
	"bufio"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	pb "../ProtocolBuffers/ProtoPackage"
	"../config"
	"../logger"
	"../membership"
	"../networking"
	"github.com/golang/protobuf/ptypes"
)

var (
	localMessage *pb.MembershipServiceMessage
	mux          sync.Mutex
	failureList  map[string]bool
	selfID       string
	introducerIP string
	isSending    bool
	isIntroducer bool
	isJoining    bool
)

func getLocalIPAddr() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		logger.PrintError("net.Dial")
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func changeStrategy(input string) {
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
		networking.HeartbeatGossip(localMessage, config.GOSSIP_FANOUT, selfID)
	} else {
		networking.HeartbeatAllToAll(localMessage, selfID)
	}
	mux.Unlock()

	selfID = ""
	localMessage = nil
	logger.PrintInfo("Successfully left")
}

func handleCommands(input string) {
	args := strings.Split(input, " ")
	cmd := args[0]
	param := ""

	if len(args) > 1 {
		param = args[1]
	}

	switch cmd {
	case "strat":
		changeStrategy(param)
	case "list":
		if param == "membership" {
			if localMessage != nil {
				mux.Lock()
				logger.PrintInfo("Printing membership list:\n", membership.GetMembershipListString(localMessage, failureList))
				mux.Unlock()
			} else {
				logger.PrintInfo("Membership list is nil")
			}
		} else if param == "self" {
			if selfID == "" {
				logger.PrintInfo("selfID is non-existent")
			} else {
				logger.PrintInfo(selfID)
			}
		} else {
			logger.PrintError("Invalid argument to 'list'")
		}
	case "join":
		if param == "" {
			logger.PrintInfo("Please specify introducer IP address for joining")
		} else if !isSending {
			introducerIP = param
			initMembershipList(true)
			isJoining = true
			isSending = true
			go startHeartbeat()
			logger.PrintInfo("Successfully sent join request")
		} else {
			logger.PrintError("Cannot join, already actively sending")
		}
	case "leave":
		sendLeaveRequest()
	default:
		logger.PrintError("Invalid command")
	}
}

func readNewMessage(message []byte) error {
	if !isSending {
		return nil
	}

	remoteMessage, err := networking.DecodeMembershipServiceMessage(message)
	if err != nil {
		return err
	}

	mux.Lock()

	if isJoining && remoteMessage.Type == pb.MessageType_JOINREP {
		isJoining = false
		localMessage.Type = pb.MessageType_STANDARD
	}

	if !isIntroducer && remoteMessage.Type == pb.MessageType_JOINREQ {
		mux.Unlock()
		return nil
	}

	membership.MergeMembershipLists(localMessage, remoteMessage, failureList)

	if isIntroducer && remoteMessage.Type == pb.MessageType_JOINREQ {
		logger.PrintInfo("Received join request")
		localMessage.Type = pb.MessageType_JOINREP
		message, err := networking.EncodeMembershipServiceMessage(localMessage)
		localMessage.Type = pb.MessageType_STANDARD

		if err != nil {
			return err
		}

		dests := membership.GetOtherMembershipListIPs(remoteMessage, selfID)
		networking.Send(dests[0], message)
	}

	mux.Unlock()

	return nil
}

func startHeartbeat() {
	for isSending {
		mux.Lock()

		localMessage.MemberList[selfID].LastSeen = ptypes.TimestampNow()
		localMessage.MemberList[selfID].HeartbeatCounter++
		membership.CheckAndRemoveMembershipListFailures(localMessage, &failureList)
		logger.InfoLogger.Println("Current memlist:\n", membership.GetMembershipListString(localMessage, failureList), "\n")

		if isJoining {
			message, _ := networking.EncodeMembershipServiceMessage(localMessage)
			networking.Send(introducerIP, message)
		} else {
			if localMessage.Strategy == config.STRAT_GOSSIP {
				networking.HeartbeatGossip(localMessage, config.GOSSIP_FANOUT, selfID)
			} else {
				networking.HeartbeatAllToAll(localMessage, selfID)
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

func initMembershipList(isGossip bool) {
	selfMember := pb.Member{
		HeartbeatCounter: 1,
		LastSeen:         ptypes.TimestampNow(),
	}

	strat := config.STRAT_GOSSIP

	if !isGossip {
		strat = config.STRAT_ALL
	}

	localMessage = &pb.MembershipServiceMessage{
		MemberList:      make(map[string]*pb.Member),
		Strategy:        strat,
		StrategyCounter: 1,
	}

	if isIntroducer {
		localMessage.Type = pb.MessageType_STANDARD
	} else {
		localMessage.Type = pb.MessageType_JOINREQ
	}

	localIP := getLocalIPAddr().String()
	selfID = localIP + ":" + ptypes.TimestampString(selfMember.LastSeen)

	membership.AddMemberToMembershipList(localMessage, selfID, &selfMember)
}

func Run(isIntro bool, isGossip bool, introIP string) {
	logger.PrintInfo("Starting detector\nIs introducer:", isIntro, "\nintroducerIP:", introIP, "\nIs gossip:", isGossip)
	isIntroducer = isIntro
	introducerIP = introIP

	isSending = true
	isJoining = !isIntroducer

	initMembershipList(isGossip)
	failureList = make(map[string]bool)

	logger.PrintInfo("Starting server with id", selfID, "on port", config.PORT)
	go networking.Listen(config.PORT, readNewMessage)
	go startHeartbeat()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()

		logger.InfoLogger.Println("Commandline input:", input)

		handleCommands(input)
	}

	if scanner.Err() != nil {
		logger.PrintError("Error reading input from commandline")
	}
}

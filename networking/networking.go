package networking

import (
	"math/rand"
	"net"
	"time"

	pb "../ProtocolBuffers/ProtoPackage"
	"../config"
	"../logger"
	"../membership"
	proto "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

var BytesSent int = 0
var MessageLossRate float64 = -1

func EncodeMembershipServiceMessage(serviceMessage *pb.MembershipServiceMessage) ([]byte, error) {
	message, err := proto.Marshal(serviceMessage)

	return message, err
}

func DecodeMembershipServiceMessage(message []byte) (*pb.MembershipServiceMessage, error) {
	list := &pb.MembershipServiceMessage{}
	err := proto.Unmarshal(message, list)

	return list, err
}

func SendGossip(serviceMessage *pb.MembershipServiceMessage, k int, selfID string) error {
	message, err := EncodeMembershipServiceMessage(serviceMessage)
	if err != nil {
		return err
	}

	dests := membership.GetOtherMembershipListIPs(serviceMessage, selfID)

	if k < len(serviceMessage.MemberList) {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(dests), func(i, j int) {
			dests[i], dests[j] = dests[j], dests[i]
		})

		dests = dests[:k]
	}

	return SendAll(dests, message)
}

func SendHeartbeat(fullMessage, selfMessage *pb.MembershipServiceMessage, selfID string) error {
	dests := membership.GetOtherMembershipListIPs(fullMessage, selfID)

	message, err := EncodeMembershipServiceMessage(selfMessage)
	if err != nil {
		return err
	}

	return SendAll(dests, message)
}

func HeartbeatGossip(serviceMessage *pb.MembershipServiceMessage, k int, selfID string) error {
	return SendGossip(serviceMessage, k, selfID)
}

func HeartbeatAllToAll(serviceMessage *pb.MembershipServiceMessage, selfID string) error {
	selfMessage := &pb.MembershipServiceMessage{
		MemberList:      make(map[string]*pb.Member),
		Strategy:        serviceMessage.Strategy,
		StrategyCounter: serviceMessage.StrategyCounter,
	}

	selfMember := pb.Member{
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
		addr, err := net.ResolveUDPAddr("udp", dest+":"+config.PORT)
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

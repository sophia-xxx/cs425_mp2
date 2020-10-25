package member_service

import (
	"sort"
	"strings"
	"sync"

	"github.com/golang/protobuf/ptypes"

	"cs425_mp2/command_util"
	"cs425_mp2/config"
	"cs425_mp2/member_service/protocol_buffer"
	"cs425_mp2/util"
	"cs425_mp2/util/logger"
)

var (
	localMessage *protocol_buffer.MembershipServiceMessage
	mux          sync.Mutex
	isSending    bool
	isJoining    bool
)

// vars that are accessible for other packages
var (
	failureList  map[string]bool
	selfID       string
	masterIP 	string
	isMaster 	bool
)

func RunService(isMaster bool, isGossip bool, MasterIP string) {
	logger.PrintInfo(
		"Starting detector\n",
		"Is introducer:", isMaster,
		"introducerIP:", MasterIP,
		"Is gossip:", isGossip)

	isMaster = isMaster
	MasterIP = MasterIP

	isSending = true
	isJoining = !isMaster

	initMembershipList(isGossip)
	failureList = make(map[string]bool)

	logger.PrintInfo("Member service is now running with id", selfID, "on port", config.MemberServicePort)
	go Listen(config.MemberServicePort, readNewMessage)
	go startHeartbeat()
}

func HandleCommand(command command_util.Command) {
	cmd := command.Method
	var param string
	if len(command.Params) >= 1 {
		param = command.Params[0]
	} else {
		param = ""
	}

	switch cmd {
	case "strat":
		ChangeStrategy(param)
	case command_util.CommandDisplay:
		if param == "membership" {
			if localMessage != nil {
				mux.Lock()
				logger.PrintInfo("Printing membership list:\n", GetMembershipListString(localMessage, failureList))
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
			logger.PrintError("Invalid argument to 'list':", param)
		}
	case command_util.CommandJoin:
		if param == "" {
			logger.PrintInfo("Please specify introducer IP address for joining")
		} else if !isSending {
			masterIP = param
			initMembershipList(true)
			isJoining = true
			isSending = true
			go startHeartbeat()
			logger.PrintInfo("Successfully sent join request")
		} else {
			logger.PrintError("Cannot join, already actively sending")
		}
	case command_util.CommandLeave:
		sendLeaveRequest()
	default:
		logger.PrintError("Invalid command_util")
	}
}



func initMembershipList(isGossip bool) {
	selfMember := protocol_buffer.Member{
		HeartbeatCounter: 1,
		LastSeen:         ptypes.TimestampNow(),
	}

	strat := config.STRAT_GOSSIP

	if !isGossip {
		strat = config.STRAT_ALL
	}

	localMessage = &protocol_buffer.MembershipServiceMessage{
		MemberList:      make(map[string]*protocol_buffer.Member),
		Strategy:        strat,
		StrategyCounter: 1,
	}

	if isMaster {
		localMessage.Type = protocol_buffer.MessageType_STANDARD
	} else {
		localMessage.Type = protocol_buffer.MessageType_JOINREQ
	}

	localIP := util.GetLocalIPAddr().String()
	selfID = localIP + ":" + ptypes.TimestampString(selfMember.LastSeen)

	AddMemberToMembershipList(localMessage, selfID, &selfMember)
}



// interface:


func GetOtherAliveMemberIPList() []string {
	ipList := make([]string, 0)
	for machineID, member := range localMessage.MemberList {
		if !failureList[machineID] && !member.IsLeaving {
			if machineID == selfID {
				continue
			}
			ip := strings.Split(machineID, ":")[0]
			ipList = append(ipList, ip)
		}
	}
	sort.Strings(ipList)
	return ipList
}

func GetFailNodeList() []string {
	failNodes := make([]string, 0)
	for k := range failureList {
		if failureList[k] {
			failNodes = append(failNodes, k)
		}
	}
	return failNodes
}

func GetMasterIP() string {
	return masterIP
}

func GetSelfID() string {
	return selfID
}

func IsMaster() bool {
	return isMaster
}


/*
This package provides member service, including:
	1. membership list
	2. failure detector
	3. simple master election

Credit: This package is adapted from CS425 Fall Recommended MP1 Solutions.
*/
package member_service

import (
	"sort"
	"strings"
	"sync"

	"cs425_mp2/command_util"
	"cs425_mp2/config"
	"cs425_mp2/member_service/protocol_buffer"
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
	failureList  	map[string]bool
	selfID      	string
	masterIP 		string
	isMaster 		bool
)

// used for election
var MasterChanged = make(chan int)

// the entry point of the package, run the member service
func RunService(isMasterBool bool, isGossipBool bool, MasterIPString string) {
	isMaster = isMasterBool
	masterIP = MasterIPString

	isSending = true
	isJoining = !isMaster

	initMembershipList(isGossipBool)
	failureList = make(map[string]bool)

	go Listen(config.MemberServicePort, readNewMessage)
	go startHeartbeat()

	logger.PrintInfo(
		"Member Service is now running\n",
		"\tPort:", config.MemberServicePort,
		"\tIs Master:", isMaster,
		"\tMasterIP:", masterIP,
		"\tIs gossip:", isGossipBool,
		"\n",
		"\tMember Self ID:", selfID)
}

// input a user command and let the service handle it
func HandleCommand(command command_util.Command) {
	cmd := command.Method
	var param string
	if len(command.Params) >= 1 {
		param = command.Params[0]
	} else {
		param = ""
	}

	switch cmd {
	case command_util.CommandSwitch:
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
	}
}


/*
	Following methods are exported for other packages so that they can access the membership list
 */

func GetAliveMemberIPList() []string {
	ipList := make([]string, 0)
	for machineID, member := range localMessage.MemberList {
		if !failureList[machineID] && !member.IsLeaving {
			//if machineID == selfID {
			//	continue
			//}
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

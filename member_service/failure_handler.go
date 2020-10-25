package member_service

import (
	"cs425_mp2/util/logger"
	"sort"
	"strings"
)

// MachineID to be in format IP:timestamp
func HandleMemberFailure(machineID string) {
	ip := strings.Split(machineID, ":")[0]
	if ip == masterIP {
		logger.PrintInfo("Master has failed. Waiting for electing a new Master:")
		Election()
	}
}

func Election() {

}

func getOldestServer() {
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
}
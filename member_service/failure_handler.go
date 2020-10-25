package member_service

import (
	"cs425_mp2/config"
	"cs425_mp2/util/logger"
	"strings"
	"time"
)

// MachineID to be in format IP:timestamp
func HandleMemberFailure(machineID string) {
	ip := strings.Split(machineID, ":")[0]
	if ip == masterIP {
		logger.PrintInfo("Master is down. Please waiting for electing a new Master...")
		go Election()
	}
}

func Election() {
	time.Sleep(config.WaitTimeForElection * time.Second)
	logger.PrintInfo("Begin electing a new master:")
	newMaster := getLargestAliveServer()
	logger.PrintInfo("New master is selected:", newMaster)
	masterIP = strings.Split(newMaster, ":")[0]
	MasterChanged <- 1
}

func getLargestAliveServer() string {
	largestID := selfID
	for machineID, member := range localMessage.MemberList {
		if !failureList[machineID] && !member.IsLeaving {
			if strings.Compare(machineID, largestID) > 0 {
				largestID = machineID
			}
		}
	}
	return largestID
}
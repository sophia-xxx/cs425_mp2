package member_service

/*
This file is used to handle a failure of node, especially, the failure of the master.
Once the master is detected to be failed, the member service would wait and then begin an election.
The election is implemented simply by choosing the node with biggest unique ID.
 */

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
	newMasterID := getLargestAliveServer()
	if selfID == newMasterID {
		logger.PrintInfo("This server has been elected as the new master.")
	} else {
		logger.PrintInfo("New master is selected:", newMasterID)
	}
	masterIP = strings.Split(newMasterID, ":")[0]
	// notify the file service
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
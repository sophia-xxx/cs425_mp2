package member_service

import (
	"cs425_mp2/util/logger"
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
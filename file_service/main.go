/*
This package provides file service, including:
	1. file transfer
	2. file record: the master would keep a file record about how files are replicated across nodes
	3. failure-related handle

SubPackages:
	- command_handler: handle commands from user
	- failure_handler: handle the failure of master
	- file_manager: manage local files
	- file_record: manage file record
	- message_handler: handle the incoming message from other nodes
	- networking: networking-related function
	- protocol_buffer: protocol buffer files for file service
*/
package file_service

import (
	"cs425_mp2/command_util"
	"cs425_mp2/config"
	"cs425_mp2/file_service/command_handler"
	"cs425_mp2/file_service/failure_handler"
	"cs425_mp2/file_service/file_manager"
	"cs425_mp2/file_service/file_record"
	"cs425_mp2/file_service/message_handler"
	"cs425_mp2/file_service/networking"
	"cs425_mp2/member_service"
	"cs425_mp2/util"
	"cs425_mp2/util/logger"
	"time"
)

func HandleCommand(command command_util.Command) {
	cmd := command.Method

	switch cmd {
	case command_util.CommandPut:
		localFilename := command.Params[0]
		sdfsFilename := command.Params[1]
		command_handler.HandlePutCommand(localFilename, sdfsFilename)
	case command_util.CommandGet:
		localFilename := command.Params[0]
		sdfsFilename := command.Params[1]
		command_handler.HandleGetCommand(sdfsFilename, localFilename)
	case command_util.CommandDelete:
		sdfsFilename := command.Params[0]
		command_handler.HandleDeleteCommand(sdfsFilename)
	case command_util.CommandList:
		sdfsFilename := command.Params[0]
		command_handler.HandleListCommand(sdfsFilename)
	case command_util.CommandStore:
		command_handler.HandleStoreCommand()
	}
}

func RunService() {
	file_manager.RemoveAllSDFSFile()

	// listen TCP message
	go message_handler.HandleFileMessage()

	logger.PrintInfo(
		"File service is now running:\n",
		"\tIPv4:", util.GetLocalIPAddr().String(),
		"\tFile Service Port:", config.FileServicePort,
		"\tFile Transfer Port:", config.FileTransferPort,
	)

	// loop check
	for {
		time.Sleep(config.FileCheckGapSeconds)
		// master node maintain file-node list
		if member_service.IsMaster() {
			file_record.RemoveFailedNodes()
			networking.CheckReplicate()
		}
		failure_handler.HandleMasterFailure()
	}
}

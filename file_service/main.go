package file_service

import (
	"time"

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

	logger.PrintInfo("Starting file_service on", util.GetLocalIPAddr().String() + ":" + config.FileServicePort)

	// master node maintain file-node list
	if member_service.IsMaster() {
		go func(){
			time.Sleep(config.FileCheckGapSeconds)
			file_record.RemoveFailNode()
		}()
		go func(){
			time.Sleep(config.FileCheckGapSeconds)
			networking.CheckReplicate()
		}()
		go func(){
			time.Sleep(config.FileCheckGapSeconds)
			failure_handler.HandleMasterFailure()
		}()
	}
	
	// listen TCP message
	go message_handler.HandleFileMessage()
}

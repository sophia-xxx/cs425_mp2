package command_handler

import (
	"cs425_mp2/config"
	"cs425_mp2/member_service"
	"io/ioutil"
	"os"
	"strconv"

	"cs425_mp2/util"
	"cs425_mp2/util/logger"

	"cs425_mp2/file_service/networking"
	pbm "cs425_mp2/file_service/protocl_buffer"
)

// deal with "put" command
func HandlePutCommand(localFileName string, sdfsFileName string) {
	fileInfo, err := os.Stat(config.LOCAL_DIR + localFileName)
	if err != nil {
		logger.PrintWarning("No such file in local file directory!")
		return
	}

	logger.PrintInfo("\nLocal file size is " + strconv.FormatInt(fileInfo.Size(), 10))
	fileMessage := &pbm.TCPMessage{
		Type:      pbm.MsgType_PUT_MASTER,
		SenderIP:  util.GetLocalIPAddr().String(),
		FileName:  sdfsFileName,
		LocalPath: config.LOCAL_DIR + localFileName,
		FileSize:  int32(fileInfo.Size()),
	}
	message, err := networking.EncodeTCPMessage(fileMessage)
	if err != nil {
		logger.PrintWarning("Encode error!")
	}
	logger.PrintInfo("Sent message to master!")
	networking.SendMessageViaTCP(member_service.GetMasterIP(), message)
}

// deal with "get" command
func HandleGetCommand(sdfsFileName string, localFileName string) {
	fileMessage := &pbm.TCPMessage{
		Type:     pbm.MsgType_GET_MASTER,
		SenderIP: util.GetLocalIPAddr().String(),
		FileName: sdfsFileName,
	}
	message, _ := networking.EncodeTCPMessage(fileMessage)
	networking.SendMessageViaTCP(member_service.GetMasterIP(), message)
}

//deal with "delete" command
func HandleDeleteCommand(sdfsFileName string) {
	fileMessage := &pbm.TCPMessage{
		Type:     pbm.MsgType_DELETE_MASTER,
		SenderIP: util.GetLocalIPAddr().String(),
		FileName: sdfsFileName,
	}
	message, _ := networking.EncodeTCPMessage(fileMessage)
	networking.SendMessageViaTCP(member_service.GetMasterIP(), message)
}

// deal with "list" command
func HandleListCommand(sdfsFileName string) {
	fileMessage := &pbm.TCPMessage{
		Type:     pbm.MsgType_LIST,
		SenderIP: util.GetLocalIPAddr().String(),
		FileName: sdfsFileName,
	}
	message, _ := networking.EncodeTCPMessage(fileMessage)
	networking.SendMessageViaTCP(member_service.GetMasterIP(), message)
}

// deal with "store" command
func HandleStoreCommand() {
	logger.PrintToConsole("Local file directory:")
	localFile, _ := ioutil.ReadDir(config.LOCAL_DIR)
	// change file to string
	for _, file := range localFile {
		logger.PrintToConsole("\t", file.Name(), ":", strconv.FormatInt(file.Size(), 10), "bytes")
	}
	logger.PrintToConsole("\nSDFS file directory:")
	files, _ := ioutil.ReadDir(config.SDFS_DIR)
	for _, file := range files {
		logger.PrintToConsole("\t", file.Name(), "\t:\t", strconv.FormatInt(file.Size(), 10), "bytes")
	}
}
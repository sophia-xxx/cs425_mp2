package detector

import (
	pbm "cs425_mp2/ProtocolBuffers/MessagePackage"
	"cs425_mp2/config"
	"cs425_mp2/logger"
	"io/ioutil"
	"os"
	"strconv"
)

// deal with "put" command
func putFileCommand(localFileName string, sdfsFileName string) {
	fileInfo, err := os.Stat(config.LOCAL_DIR + localFileName)
	logger.PrintInfo("Local file size is " + strconv.FormatInt(fileInfo.Size(), 10))
	if err != nil {
		logger.PrintInfo("No such file!")
	}
	fileMessage := &pbm.TCPMessage{
		Type:      pbm.MsgType_PUT_MASTER,
		SenderIP:  GetLocalIPAddr().String(),
		FileName:  sdfsFileName,
		LocalPath: config.LOCAL_DIR + localFileName,
		FileSize:  int32(fileInfo.Size()),
	}
	message, err := EncodeTCPMessage(fileMessage)
	if err != nil {
		logger.PrintInfo("Encode error!")
	}
	logger.PrintInfo("Send message to master!")
	SendMessage(introducerIp, message)

}

// deal with "get" command
func getFileCommand(sdfsFileName string, localFileName string) {
	fileMessage := &pbm.TCPMessage{
		Type:     pbm.MsgType_GET_MASTER,
		SenderIP: GetLocalIPAddr().String(),
		FileName: sdfsFileName,
	}
	message, _ := EncodeTCPMessage(fileMessage)
	SendMessage(introducerIp, message)

}

//deal with "delete" command
func deleteFileCommand(sdfsFileName string) {
	fileMessage := &pbm.TCPMessage{
		Type:     pbm.MsgType_DELETE_MASTER,
		SenderIP: GetLocalIPAddr().String(),
		FileName: sdfsFileName,
	}
	message, _ := EncodeTCPMessage(fileMessage)
	SendMessage(introducerIp, message)
}

// deal with "list" command
func listFileCommand(sdfsFileName string) {
	fileMessage := &pbm.TCPMessage{
		Type:     pbm.MsgType_LIST,
		SenderIP: GetLocalIPAddr().String(),
		FileName: sdfsFileName,
	}
	message, _ := EncodeTCPMessage(fileMessage)
	SendMessage(introducerIp, message)
}

// deal with "store" command
func StoreCommand() {
	files, _ := ioutil.ReadDir(config.SDFS_DIR)
	for _, f := range files {
		logger.PrintInfo(f.Name())
	}

}

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
	if err != nil {
		logger.PrintInfo("\n No such file in local file directory!")
		return
	}

	logger.PrintInfo("\nLocal file size is " + strconv.FormatInt(fileInfo.Size(), 10))
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
	logger.PrintInfo("\nSend message to master!")
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
	logger.PrintInfo("\nLocal file directory: \n")
	localFile, _ := ioutil.ReadDir(config.LOCAL_DIR)
	// change file to string
	for _, file := range localFile {
		logger.PrintInfo(file.Name() + ":  " + strconv.FormatInt(file.Size(), 10) + "B")
	}
	logger.PrintInfo("\nSDFS file directory: \n")
	files, _ := ioutil.ReadDir(config.SDFS_DIR)
	for _, f := range files {
		logger.PrintInfo(f.Name() + ":  " + strconv.FormatInt(f.Size(), 10) + "B")
	}

}

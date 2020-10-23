package detector

import (
	pbm "../ProtocolBuffers/MessagePackage"
	"../config"
	"../connection"
	//"fmt"
	"../logger"
	"io/ioutil"
)

var introducerIp string

// deal with "put" command
func putFileCommand(localFileName string, sdfsFileName string) {
	fileMessage := &pbm.TCPMessage{
		Type:      pbm.MsgType_PUT_MASTER,
		SenderIP:  GetLocalIPAddr().String(),
		FileName:  sdfsFileName,
		LocalPath: config.LOCAL_DIR + localFileName,
	}
	message, _ := connection.EncodeTCPMessage(fileMessage)
	connection.SendMessage(introducerIp, message)

}

// deal with "get" command
func getFileCommand(sdfsFileName string, localFileName string) {
	fileMessage := &pbm.TCPMessage{
		Type:     pbm.MsgType_GET_MASTER,
		SenderIP: GetLocalIPAddr().String(),
		FileName: sdfsFileName,
	}
	message, _ := connection.EncodeTCPMessage(fileMessage)
	connection.SendMessage(introducerIp, message)

}

//deal with "delete" command
func deleteFileCommand(sdfsFileName string) {

}

// deal with "list" command
func listFileCommand() {

}

// deal with "store" command
func StoreCommand() {
	files, _ := ioutil.ReadDir(config.SDFS_DIR)
	for _, f := range files {
		logger.PrintInfo(f.Name())
	}

}

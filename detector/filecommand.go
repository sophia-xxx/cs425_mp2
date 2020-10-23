package detector

import (
	pbm "../ProtocolBuffers/MessagePackage"
	"../config"
	"../connection"
	"os"

	//"fmt"
	"../logger"
	"io/ioutil"
)

var introducerIp string

// deal with "put" command
func putFileCommand(localFileName string, sdfsFileName string) {
	fileInfo, _ := os.Stat(localFileName)
	fileMessage := &pbm.TCPMessage{
		Type:      pbm.MsgType_PUT_MASTER,
		SenderIP:  GetLocalIPAddr().String(),
		FileName:  sdfsFileName,
		LocalPath: config.LOCAL_DIR + localFileName,
		FileSize:  int32(fileInfo.Size()),
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
	fileMessage := &pbm.TCPMessage{
		Type:     pbm.MsgType_DELETE_MASTER,
		SenderIP: GetLocalIPAddr().String(),
		FileName: sdfsFileName,
	}
	message, _ := connection.EncodeTCPMessage(fileMessage)
	connection.SendMessage(introducerIp, message)
}

// deal with "list" command
func listFileCommand(sdfsFileName string) {
	fileMessage := &pbm.TCPMessage{
		Type:     pbm.MsgType_LIST,
		SenderIP: GetLocalIPAddr().String(),
		FileName: sdfsFileName,
	}
	message, _ := connection.EncodeTCPMessage(fileMessage)
	connection.SendMessage(introducerIp, message)
}

// deal with "store" command
func StoreCommand() {
	files, _ := ioutil.ReadDir(config.SDFS_DIR)
	for _, f := range files {
		logger.PrintInfo(f.Name())
	}

}

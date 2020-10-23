package connection

import (
	pbm "../ProtocolBuffers/MessagePackage"
	"../detector"
	"../logger"
	"../master"
	"os"
)

func deleteMessageHandle(remoteMsg *pbm.TCPMessage) {
	// master send DELETE message to target nodes
	if isMaster && remoteMsg.Type == pbm.MsgType_DELETE_MASTER {
		master.DeleteMessage(remoteMsg.FileName)
	}
	// master get delete ACK then update file-node list
	if isMaster && remoteMsg.Type == pbm.MsgType_DELETE_ACK {
		master.DeleteFileRecord(remoteMsg.FileName, remoteMsg.SenderIP)
	}

	if remoteMsg.Type == pbm.MsgType_DELETE {
		deleteFile(remoteMsg.FileName)
	}

}

func deleteFile(filename string) {
	os.Remove(filename)
	sendDeleteACK(filename)
}

// Send delete success message to master
func sendDeleteACK(sdfsFileName string) {
	fileMessage := &pbm.TCPMessage{
		Type:     pbm.MsgType_DELETE_ACK,
		SenderIP: detector.GetLocalIPAddr().String(),
		FileName: sdfsFileName,
	}
	message, _ := EncodeTCPMessage(fileMessage)
	SendMessage(introducerIp, message)
}

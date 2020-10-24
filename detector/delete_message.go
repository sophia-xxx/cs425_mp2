package detector

import (
	pbm "../ProtocolBuffers/MessagePackage"
	"os"
)

func deleteMessageHandle(remoteMsg *pbm.TCPMessage) {
	// master send DELETE message to target nodes
	if isIntroducer && remoteMsg.Type == pbm.MsgType_DELETE_MASTER {
		DeleteMessage(remoteMsg.FileName)
	}
	// master get delete ACK then update file-node list
	if isIntroducer && remoteMsg.Type == pbm.MsgType_DELETE_ACK {
		DeleteFileRecord(remoteMsg.FileName, remoteMsg.SenderIP)
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
		SenderIP: GetLocalIPAddr().String(),
		FileName: sdfsFileName,
	}
	message, _ := EncodeTCPMessage(fileMessage)
	SendMessage(introducerIp, message)
}

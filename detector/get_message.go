package detector

import (

	//"strings"
	pbm "cs425_mp2/ProtocolBuffers/MessagePackage"
	"cs425_mp2/logger"
	"os"

	//"fmt"
	"cs425_mp2/config"
)

func getMessageHandler(remoteMsg *pbm.TCPMessage) {
	// master return target node to read
	if isIntroducer && remoteMsg.Type == pbm.MsgType_GET_MASTER {
		GetReplyMessage(remoteMsg.FileName, remoteMsg.SenderIP)
	}

	if remoteMsg.Type == pbm.MsgType_GET_MASTER_REP {
		// receive file from target nodes
		targetList := remoteMsg.PayLoad

		if targetList == nil {
			logger.PrintInfo(remoteMsg.FileName + "  has no record!")
		} else {
			sendReadReq(targetList[0], remoteMsg.FileName)
		}

	}
	// server reply to get request and send file to client
	if remoteMsg.Type == pbm.MsgType_GET_P2P {
		sendReadReply(remoteMsg.SenderIP, remoteMsg.FileName)
		sendFile(config.SDFS_DIR+remoteMsg.FileName, remoteMsg.SenderIP, remoteMsg.FileName)
	}
	// when get ACK, client start receiving file
	if remoteMsg.Type == pbm.MsgType_GET_P2P_ACK {
		ListenFile(config.LOCAL_DIR+remoteMsg.FileName, remoteMsg.FileSize, false)
	}

}

// client send read request to target node
func sendReadReq(targetIp string, sdfsFileName string) {
	fileMessage := &pbm.TCPMessage{
		Type:     pbm.MsgType_GET_P2P,
		FileName: sdfsFileName,
		SenderIP: GetLocalIPAddr().String(),
	}

	message, _ := EncodeTCPMessage(fileMessage)
	SendMessage(targetIp, message)
}
func sendReadReply(targetIp string, sdfsFileName string) {
	fileInfo, _ := os.Stat(config.SDFS_DIR + sdfsFileName)
	fileMessage := &pbm.TCPMessage{
		Type:     pbm.MsgType_GET_P2P_ACK,
		FileName: sdfsFileName,
		SenderIP: GetLocalIPAddr().String(),
		FileSize: int32(fileInfo.Size()),
	}
	message, _ := EncodeTCPMessage(fileMessage)
	SendMessage(targetIp, message)
}

package detector

import (
	pbm "cs425_mp2/ProtocolBuffers/MessagePackage"
	"cs425_mp2/config"
	"cs425_mp2/logger"
)

func putMessageHandler(remoteMsg *pbm.TCPMessage) {
	// master return target node to write
	if isIntroducer && remoteMsg.Type == pbm.MsgType_PUT_MASTER {
		PutReplyMessage(remoteMsg.FileName, remoteMsg.SenderIP, remoteMsg.FileSize)
		logger.PrintInfo("Master reply")
	}
	// client send write file request to target nodes
	if remoteMsg.Type == pbm.MsgType_PUT_MASTER_REP {
		logger.PrintInfo("Got reply from master! ")
		targetList := remoteMsg.PayLoad
		for _, target := range targetList {
			sendWriteReq(target, remoteMsg.FileName, remoteMsg.FileSize)
			logger.PrintInfo("Send write request to target  " + target)
		}
	}
	// server send ACK to put request and start file socket
	if remoteMsg.Type == pbm.MsgType_PUT_P2P {
		sendWriteReply(remoteMsg.SenderIP, remoteMsg.FileName)
		logger.PrintInfo("Got ACK from target  ")
		ListenFile(config.SDFS_DIR+remoteMsg.FileName, remoteMsg.FileSize, true)
		logger.PrintInfo("Finish receiving file  ")
	}
	// client start sending file
	if remoteMsg.Type == pbm.MsgType_PUT_P2P_ACK {
		logger.PrintInfo("Start sending file  ")
		sendFile(config.LOCAL_DIR+remoteMsg.FileName, remoteMsg.SenderIP, remoteMsg.FileName)
		logger.PrintInfo("Finish sending file  ")

	}
	// when write finish, master will receive write ACK to maintain file-node list
	if isIntroducer && remoteMsg.Type == pbm.MsgType_WRITE_ACK {
		// quorum determine whether the write is succeed
		logger.PrintInfo("Master got ACK from file node  ")
		ipList := make([]string, 0)
		ipList = append(ipList, remoteMsg.SenderIP)
		UpdateFileNode(remoteMsg.FileName, ipList)
	}
}

// client send write request to target nodes
func sendWriteReq(targetIp string, sdfsFileName string, fileSize int32) {
	fileMessage := &pbm.TCPMessage{
		Type:     pbm.MsgType_PUT_P2P,
		FileName: sdfsFileName,
		SenderIP: GetLocalIPAddr().String(),
		FileSize: fileSize,
	}
	message, _ := EncodeTCPMessage(fileMessage)
	SendMessage(targetIp, message)
}

func sendWriteReply(targetIp string, sdfsFileName string) {
	fileMessage := &pbm.TCPMessage{
		Type:     pbm.MsgType_PUT_P2P_ACK,
		FileName: sdfsFileName,
		SenderIP: GetLocalIPAddr().String(),
	}
	message, _ := EncodeTCPMessage(fileMessage)
	SendMessage(targetIp, message)
}

// when server finish put file, send ACK to master
func SendWriteACK(targetIp string, sdfsFileName string) {
	fileMessage := &pbm.TCPMessage{
		Type:     pbm.MsgType_WRITE_ACK,
		FileName: sdfsFileName,
		SenderIP: GetLocalIPAddr().String(),
	}

	message, _ := EncodeTCPMessage(fileMessage)
	SendMessage(targetIp, message)
}

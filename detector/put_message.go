package detector

import (
	pbm "cs425_mp2/ProtocolBuffers/MessagePackage"
	"cs425_mp2/config"
	"cs425_mp2/logger"
)

func putMessageHandler(remoteMsg *pbm.TCPMessage) {
	// master return target node to write
	if isIntroducer && remoteMsg.Type == pbm.MsgType_PUT_MASTER {
		logger.PrintInfo("Master receive put request")
		PutReplyMessage(remoteMsg)
		logger.PrintInfo("Master reply")
	}
	// client send write file request to target nodes
	if remoteMsg.Type == pbm.MsgType_PUT_MASTER_REP {
		logger.PrintInfo("Got  " + pbm.MsgType_name[int32(remoteMsg.Type)] + "  from master with filename: " + remoteMsg.FileName)
		targetList := remoteMsg.PayLoad
		// for _, target := range targetList {
		// 	sendWriteReq(target, remoteMsg.FileName, remoteMsg.FileSize)
		// 	logger.PrintInfo("Send write request to target  " + target)
		// }
		target := targetList[0]
		sendWriteReq(target, remoteMsg)
		logger.PrintInfo("Send write request to target  " + target)
	}
	// server send ACK to put request and start file socket
	if remoteMsg.Type == pbm.MsgType_PUT_P2P {
		sendWriteReply(remoteMsg)
		logger.PrintInfo("Got put request from client  ")
		logger.PrintInfo("**" + introducerIp + "**")
		ListenFile(config.SDFS_DIR+remoteMsg.FileName, remoteMsg.FileSize, true)
		logger.PrintInfo("Finish receiving file  ")
	}
	// client start sending file
	if remoteMsg.Type == pbm.MsgType_PUT_P2P_ACK {
		logger.PrintInfo("Start sending file whose filename is: " + remoteMsg.FileName)
		sendFile(remoteMsg.LocalPath, remoteMsg.SenderIP, remoteMsg.FileName)
		logger.PrintInfo("Finish sending file  " + remoteMsg.FileName)

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
func sendWriteReq(targetIp string, remoteMsg *pbm.TCPMessage) {
	fileMessage := &pbm.TCPMessage{
		Type:      pbm.MsgType_PUT_P2P,
		FileName:  remoteMsg.FileName,
		SenderIP:  GetLocalIPAddr().String(),
		FileSize:  remoteMsg.FileSize,
		LocalPath: remoteMsg.LocalPath,
	}
	message, _ := EncodeTCPMessage(fileMessage)
	logger.PrintInfo("Send putp2p mes with filename:" + fileMessage.FileName)
	SendMessage(targetIp, message)
}

func sendWriteReply(remoteMsg *pbm.TCPMessage) {
	fileMessage := &pbm.TCPMessage{
		Type:      pbm.MsgType_PUT_P2P_ACK,
		FileName:  remoteMsg.FileName,
		SenderIP:  GetLocalIPAddr().String(),
		FileSize:  remoteMsg.FileSize,
		LocalPath: remoteMsg.LocalPath,
	}
	message, _ := EncodeTCPMessage(fileMessage)
	SendMessage(remoteMsg.SenderIP, message)
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

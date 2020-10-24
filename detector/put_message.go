package detector

import (
	pbm "cs425_mp2/ProtocolBuffers/MessagePackage"
	"cs425_mp2/config"
)

func putMessageHandler(remoteMsg *pbm.TCPMessage) {
	// master return target node to write
	if isIntroducer && remoteMsg.Type == pbm.MsgType_PUT_MASTER {
		PutReplyMessage(remoteMsg.FileName, remoteMsg.SenderIP, remoteMsg.FileSize)
	}
	// client send write file request to target nodes
	if remoteMsg.Type == pbm.MsgType_PUT_MASTER_REP {
		targetList := remoteMsg.PayLoad
		for _, target := range targetList {
			sendWriteReq(target, remoteMsg.FileName, remoteMsg.FileSize)

		}
	}
	// server send ACK to put request and start file socket
	if remoteMsg.Type == pbm.MsgType_PUT_P2P {
		sendWriteReply(remoteMsg.SenderIP, remoteMsg.FileName)
		ListenFile(config.SDFS_DIR+remoteMsg.FileName, remoteMsg.FileSize, true)
	}
	// client start sending file
	if remoteMsg.Type == pbm.MsgType_PUT_P2P_ACK {
		sendFile(config.LOCAL_DIR+remoteMsg.FileName, remoteMsg.SenderIP, remoteMsg.FileName)
	}
	// when write finish, master will receive write ACK to maintain file-node list
	if isIntroducer && remoteMsg.Type == pbm.MsgType_WRITE_ACK {
		// quorum determine whether the write is succeed
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
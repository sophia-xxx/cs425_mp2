package message_handler

import (
	"cs425_mp2/config"
	"cs425_mp2/member_service"
	"cs425_mp2/util"
	"net"

	"cs425_mp2/util/logger"

	"cs425_mp2/file_service/file_manager"
	"cs425_mp2/file_service/file_record"
	"cs425_mp2/file_service/networking"
	pbm "cs425_mp2/file_service/protocl_buffer"
)

// socket to listen TCP message
func HandleFileMessage() {
	addressString := ":" + config.FileServicePort
	localAddr, err := net.ResolveTCPAddr("tcp4", addressString)
	if err != nil {
		logger.PrintWarning("Cannot resolve TCP address!  " + addressString)
	}
	listener, err := net.ListenTCP("tcp4", localAddr)
	if err != nil {
		logger.PrintWarning("Cannot listen TCP!")
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			logger.PrintError("Cannot open TCP connection!")
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn *net.TCPConn) {
	defer conn.Close()

	// read message data
	buf := make([]byte, config.BUFFER_SIZE)
	n, err := conn.Read(buf)
	if err != nil {
		logger.PrintInfo("Unable to read data!")
	}
	messageBytes := buf[0:n]
	remoteMsg, err := networking.DecodeTCPMessage(messageBytes)
	if err != nil || remoteMsg == nil {
		logger.PrintInfo("Cannot decode message!")
		return
	}

	// deal with all PUT relevant message
	if remoteMsg.Type <= config.PUT {
		logger.PrintInfo("Received Put message about:", remoteMsg.FileName, "from", remoteMsg.SenderIP)
		putMessageHandler(remoteMsg)
	}
	// deal with all GET relevant message
	if remoteMsg.Type > config.PUT && remoteMsg.Type <= config.GET {
		logger.PrintInfo("Received Get message about:", remoteMsg.FileName, "from", remoteMsg.SenderIP)
		getMessageHandler(remoteMsg)
	}
	// deal with all DELETE relevant message
	if remoteMsg.Type > config.GET && remoteMsg.Type <= config.DELETE {
		logger.PrintInfo("Received Delete message about:", remoteMsg.FileName, "from", remoteMsg.SenderIP)
		deleteMessageHandler(remoteMsg)
	}
	// deal with other message
	if remoteMsg.Type > config.DELETE {
		if member_service.IsMaster() && remoteMsg.Type == pbm.MsgType_LIST {
			networking.ListReplyMessage(remoteMsg.FileName, remoteMsg.SenderIP)
		}
		if remoteMsg.Type == pbm.MsgType_LIST_REP {
			listMessageHandler(remoteMsg)
		}
	}
	// deal with restore
	if remoteMsg.Type == pbm.MsgType_RESTORE {
		logger.PrintInfo("Received Restore message from", remoteMsg.SenderIP)
		restoreMessageHandler(remoteMsg)
	}
}

func getMessageHandler(remoteMsg *pbm.TCPMessage) {
	// master return target node to read
	if member_service.IsMaster() && remoteMsg.Type == pbm.MsgType_GET_MASTER {
		networking.GetReplyMessage(remoteMsg.FileName, remoteMsg.SenderIP)
	}

	if remoteMsg.Type == pbm.MsgType_GET_MASTER_REP {
		// receive file from target nodes
		targetList := remoteMsg.PayLoad
		if targetList == nil {
			logger.PrintInfo(remoteMsg.FileName + "  has no record!")
		} else {
			networking.SendReadReq(targetList[0], remoteMsg.FileName)
		}

	}
	// server reply to get request and send file to client
	if remoteMsg.Type == pbm.MsgType_GET_P2P {
		networking.SendReadReply(remoteMsg.SenderIP, remoteMsg.FileName)
		networking.SendFile(config.SDFS_DIR+remoteMsg.FileName, remoteMsg.SenderIP, remoteMsg.FileName)
	}
	// when get ACK, client start receiving file
	if remoteMsg.Type == pbm.MsgType_GET_P2P_ACK {
		networking.ListenFile(config.LOCAL_DIR+remoteMsg.FileName, remoteMsg.FileSize, false)
	}

}

func putMessageHandler(remoteMsg *pbm.TCPMessage) {
	// master return target node to write
	if member_service.IsMaster() && remoteMsg.Type == pbm.MsgType_PUT_MASTER {
		logger.PrintInfo("Master received a Put message.")
		networking.PutReplyMessage(remoteMsg)
	}
	// client send write file request to target nodes
	if remoteMsg.Type == pbm.MsgType_PUT_MASTER_REP {
		logger.PrintInfo("Got  " + pbm.MsgType_name[int32(remoteMsg.Type)] + "  from master with filename: " + remoteMsg.FileName)
		targetList := remoteMsg.PayLoad
		for _, target := range targetList {
			networking.SendWriteReq(target, remoteMsg)
			logger.PrintInfo("Send write request to target  " + target)
		}
	}
	// server send ACK to put request and start file socket
	if remoteMsg.Type == pbm.MsgType_PUT_P2P {
		networking.SendWriteReply(remoteMsg)
		logger.PrintInfo("Got put request from client  ")
		networking.ListenFile(config.SDFS_DIR+remoteMsg.FileName, remoteMsg.FileSize, true)
		logger.PrintInfo("Finish receiving file  ")
	}
	// client start sending file
	if remoteMsg.Type == pbm.MsgType_PUT_P2P_ACK {
		logger.PrintInfo("Start sending file whose filename is: " + remoteMsg.FileName)
		networking.SendFile(remoteMsg.LocalPath, remoteMsg.SenderIP, remoteMsg.FileName)
		logger.PrintInfo("Finish sending file  " + remoteMsg.FileName)
	}
	// when write finish, master will receive write ACK to maintain file-node list
	if member_service.IsMaster() && remoteMsg.Type == pbm.MsgType_WRITE_ACK {
		// quorum determine whether the write is succeed
		logger.PrintInfo("Master got ACK from file node  ")
		ipList := make([]string, 0)
		ipList = append(ipList, remoteMsg.SenderIP)
		file_record.UpdateFileNode(remoteMsg.FileName, ipList)
	}
}

func listMessageHandler(remoteMsg *pbm.TCPMessage) {
	nodeList := remoteMsg.PayLoad
	if nodeList == nil {
		logger.PrintInfo("No such file!")
	} else {
		fileString := util.ListToString(nodeList)
		logger.PrintInfo(remoteMsg.FileName + " is stored in machine : " + fileString)
	}
}

func deleteMessageHandler(remoteMsg *pbm.TCPMessage) {
	if member_service.IsMaster() {
		// master send DELETE message to target nodes
		if remoteMsg.Type == pbm.MsgType_DELETE_MASTER {
			networking.DeleteMessage(remoteMsg.FileName)
			file_record.DeleteFileInAllNodes(remoteMsg.FileName)
		}
	} else {
		if remoteMsg.Type == pbm.MsgType_DELETE {
			file_manager.RemoveSDFSFile(remoteMsg.FileName)
			networking.SendDeleteACK(remoteMsg.FileName)
		}
	}
}

func restoreMessageHandler(remoteMsg *pbm.TCPMessage) {
	if !member_service.IsMaster() {
		return
	}

	nodeIP := remoteMsg.SenderIP
	files := remoteMsg.PayLoad
	file_record.RestoreFileNode(nodeIP, files)
	logger.PrintInfo("Restored file records from node", nodeIP, files, ".")
}

package connection

import (

	//"strings"
	pbm "../ProtocolBuffers/MessagePackage"
	"github.com/golang/protobuf/ptypes"
	"os"

	//"fmt"
	"../config"
	"../detector"
	"../logger"
	"../master"
)

func getMessageHandler(remoteMsg *pbm.TCPMessage) {
	// master return target node to read
	if isMaster && remoteMsg.Type == pbm.MsgType_GET_MASTER {
		master.GetReplyMessage(remoteMsg.FileName, remoteMsg.SenderIP)
	}

	/*todo: read timeout*/
	// client send read request to target nodes
	if remoteMsg.Type == pbm.MsgType_GET_MASTER_REP {
		// receive file from target nodes
		targetList := remoteMsg.PayLoad
		for _, target := range targetList {
			get_ack_received = false
			sendReadReq(target, remoteMsg.FileName)
			startTime := float64(ptypes.TimestampNow().GetSeconds())
			for {
				if get_ack_received {
					break
				}
				curTime := float64(ptypes.TimestampNow().GetSeconds())
				if curTime-startTime > config.ACK_TIMEOUT {
					break
				} else {
					continue
				}
			}
			if !get_ack_received {
				continue
			}
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
		SenderIP: detector.GetLocalIPAddr().String(),
	}

	message, _ := EncodeTCPMessage(fileMessage)
	SendMessage(targetIp, message)
}
func sendReadReply(targetIp string, sdfsFileName string) {
	fileInfo, _ := os.Stat(config.SDFS_DIR + sdfsFileName)
	fileMessage := &pbm.TCPMessage{
		Type:     pbm.MsgType_GET_P2P_ACK,
		FileName: sdfsFileName,
		SenderIP: detector.GetLocalIPAddr().String(),
		FileSize: int32(fileInfo.Size()),
	}
	message, _ := EncodeTCPMessage(fileMessage)
	SendMessage(targetIp, message)
}

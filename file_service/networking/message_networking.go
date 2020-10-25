package networking

import (
	"cs425_mp2/config"
	"cs425_mp2/member_service"
	"github.com/golang/protobuf/proto"
	"net"
	"os"

	"cs425_mp2/util"
	"cs425_mp2/util/logger"

	"cs425_mp2/file_service/protocl_buffer"
)

// send TCP message
func SendMessageViaTCP(dest string, message []byte) {
	remoteAddress, _ := net.ResolveTCPAddr("tcp4", dest+":"+config.FileServicePort)
	conn, err := net.DialTCP("tcp4", nil, remoteAddress)
	//logger.PrintInfo("Set connection!")
	if err != nil || conn == nil {
		logger.PrintInfo("Cannot dial remote address!")
		return
	}
	_, err = conn.Write(message)
	if err != nil {
		logger.PrintInfo("Cannot send message!")
	}
}

func EncodeTCPMessage(fileMessage *protocl_buffer.TCPMessage) ([]byte, error) {
	message, err := proto.Marshal(fileMessage)
	if err != nil {
		logger.PrintInfo("Serialize error!")
	}
	return message, err
}

func DecodeTCPMessage(message []byte) (*protocl_buffer.TCPMessage, error) {
	list := &protocl_buffer.TCPMessage{}
	err := proto.Unmarshal(message, list)

	return list, err
}

// Send delete success message to master
func SendDeleteACK(sdfsFileName string) {
	fileMessage := &protocl_buffer.TCPMessage{
		Type:     protocl_buffer.MsgType_DELETE_ACK,
		SenderIP: util.GetLocalIPAddr().String(),
		FileName: sdfsFileName,
	}
	message, _ := EncodeTCPMessage(fileMessage)
	SendMessageViaTCP(member_service.GetMasterIP(), message)
}

// client send read request to target node
func SendReadReq(targetIp string, sdfsFileName string) {
	fileMessage := &protocl_buffer.TCPMessage{
		Type:     protocl_buffer.MsgType_GET_P2P,
		FileName: sdfsFileName,
		SenderIP: util.GetLocalIPAddr().String(),
	}

	message, _ := EncodeTCPMessage(fileMessage)
	SendMessageViaTCP(targetIp, message)
}
func SendReadReply(targetIp string, sdfsFileName string) {
	fileInfo, _ := os.Stat(config.SDFS_DIR + sdfsFileName)
	fileMessage := &protocl_buffer.TCPMessage{
		Type:     protocl_buffer.MsgType_GET_P2P_ACK,
		FileName: sdfsFileName,
		SenderIP: util.GetLocalIPAddr().String(),
		FileSize: int32(fileInfo.Size()),
	}
	message, _ := EncodeTCPMessage(fileMessage)
	SendMessageViaTCP(targetIp, message)
}

func SendWriteReply(remoteMsg *protocl_buffer.TCPMessage) {
	fileMessage := &protocl_buffer.TCPMessage{
		Type:      protocl_buffer.MsgType_PUT_P2P_ACK,
		FileName:  remoteMsg.FileName,
		SenderIP:  util.GetLocalIPAddr().String(),
		FileSize:  remoteMsg.FileSize,
		LocalPath: remoteMsg.LocalPath,
	}
	message, _ := EncodeTCPMessage(fileMessage)
	SendMessageViaTCP(remoteMsg.SenderIP, message)
}

// when server finish put file, send ACK to master
func SendWriteACK(targetIp string, sdfsFileName string) {
	fileMessage := &protocl_buffer.TCPMessage{
		Type:     protocl_buffer.MsgType_WRITE_ACK,
		FileName: sdfsFileName,
		SenderIP: util.GetLocalIPAddr().String(),
	}

	message, _ := EncodeTCPMessage(fileMessage)
	SendMessageViaTCP(targetIp, message)
}

// client send write request to target nodes
func SendWriteReq(targetIp string, remoteMsg *protocl_buffer.TCPMessage) {
	fileMessage := &protocl_buffer.TCPMessage{
		Type:      protocl_buffer.MsgType_PUT_P2P,
		FileName:  remoteMsg.FileName,
		SenderIP:  util.GetLocalIPAddr().String(),
		FileSize:  remoteMsg.FileSize,
		LocalPath: remoteMsg.LocalPath,
	}
	message, _ := EncodeTCPMessage(fileMessage)
	logger.PrintInfo("Send putp2p mes with filename:" + fileMessage.FileName)
	SendMessageViaTCP(targetIp, message)
}

// send file list
// send connection by TCP connection (send filename-->get ACK-->send connection)
func RestoreFileListToMaster(fileList []string, dest string) {
	fileMessage := &protocl_buffer.TCPMessage{
		Type:      	protocl_buffer.MsgType_RESTORE,
		SenderIP:  	util.GetLocalIPAddr().String(),
		PayLoad: 	fileList,
	}
	message, _ := EncodeTCPMessage(fileMessage)
	logger.PrintInfo("Send file list:", fileList)
	SendMessageViaTCP(dest, message)
}
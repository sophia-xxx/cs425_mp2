package networking

import (
	"cs425_mp2/config"
	"cs425_mp2/file_service/file_record"
	"cs425_mp2/file_service/protocl_buffer"
	"cs425_mp2/member_service"
	"cs425_mp2/util"
	"cs425_mp2/util/logger"
	"math/rand"
	"strings"
	"time"
)

// send the replicate request to one existed file node
func ReplicateFile(storeList []string, newList []string, filename string) {
	// decide which node is the good file
	// randomly choose a source node
	rand.Seed(time.Now().Unix())
	sourceNode := storeList[rand.Intn(len(storeList))]
	// ask the source to put file
	repMessage := &protocl_buffer.TCPMessage{
		Type:     protocl_buffer.MsgType_PUT_MASTER_REP,
		SenderIP: util.GetLocalIPAddr().String(),
		PayLoad:  newList,
		FileName: filename,
	}
	msgBytes, _ := EncodeTCPMessage(repMessage)

	SendMessageViaTCP(sourceNode, msgBytes)
}

// master return target node to write
func PutReplyMessage(remoteMsg *protocl_buffer.TCPMessage) {
	// check if key exist in map
	writeList := make([]string, 0)
	if file_record.FileNodeList[remoteMsg.FileName] == nil {
		logger.PrintInfo("Finding nodes to store the file.")
		writeList = file_record.FindNewNode(remoteMsg.FileName, remoteMsg.SenderIP)
	} else {
		writeList = file_record.FileNodeList[remoteMsg.FileName]
	}
	if writeList == nil {
		logger.PrintInfo("No valid node to store the file")
		return
	}
	repMessage := &protocl_buffer.TCPMessage{
		Type:      protocl_buffer.MsgType_PUT_MASTER_REP,
		SenderIP:  util.GetLocalIPAddr().String(),
		PayLoad:   writeList,
		FileName:  remoteMsg.FileName,
		FileSize:  remoteMsg.FileSize,
		LocalPath: remoteMsg.LocalPath,
	}
	msgBytes, _ := EncodeTCPMessage(repMessage)
	SendMessageViaTCP(remoteMsg.SenderIP, msgBytes)
}

func GetReplyMessage(filename string, sender string) {
	readList := file_record.FileNodeList[filename]
	validList := make([]string, 0)
	for _, node := range readList {
		if strings.Compare(node, sender) != 0 {
			validList = append(validList, node)
		}
	}
	repMessage := &protocl_buffer.TCPMessage{
		Type:     protocl_buffer.MsgType_GET_MASTER_REP,
		SenderIP: util.GetLocalIPAddr().String(),
		PayLoad:  validList,
		FileName: filename,
	}
	msgBytes, _ := EncodeTCPMessage(repMessage)
	SendMessageViaTCP(sender, msgBytes)
}

// master return target node with VM ip list that store the file
func ListReplyMessage(filename string, targetIp string) {
	repMessage := &protocl_buffer.TCPMessage{
		Type:     protocl_buffer.MsgType_LIST_REP,
		FileName: filename,
		SenderIP: util.GetLocalIPAddr().String(),
		PayLoad:  file_record.FileNodeList[filename],
	}
	msgBytes, _ := EncodeTCPMessage(repMessage)
	SendMessageViaTCP(targetIp, msgBytes)
}

//master send delete request to file node
func DeleteMessage(filename string) {
	ipList := file_record.FileNodeList[filename]
	if ipList == nil {
		logger.PrintInfo("No such file in SDFS")
		return
	}

	fileMessage := &protocl_buffer.TCPMessage{
		Type:     protocl_buffer.MsgType_DELETE,
		SenderIP: util.GetLocalIPAddr().String(),
		FileName: filename,
	}
	message, _ := EncodeTCPMessage(fileMessage)
	for _, target := range ipList {
		SendMessageViaTCP(target, message)
	}

}

// master check whether to replicate files or not---should run continuously
func CheckReplicate() {
	for file, nodeList := range file_record.FileNodeList {
		if len(nodeList) < config.REPLICA {
			logger.PrintInfo(file, "is detected to have only", len(nodeList), "replicas. Trying to add a new replica for it...")
			storeList := file_record.FileNodeList[file]
			if len(storeList) == 0 {
				logger.PrintWarning(file, "is detected to not have any replicas!")
				return
			}
			ipList := file_record.FindNewNode(file, member_service.GetMasterIP())
			if ipList == nil {
				logger.PrintInfo("No valid node to replicate")
			} else {
				ReplicateFile(storeList, ipList, file)
			}

		}
	}
}

package detector

import (
	pbm "cs425_mp2/ProtocolBuffers/MessagePackage"
	"cs425_mp2/config"
	"cs425_mp2/logger"
	"strconv"
	"strings"
)

// add or update record in file-node map
func UpdateFileNode(sdfsFileName string, newNodeList []string) {
	// update existed record
	if _, exist := fileNodeList[sdfsFileName]; exist {
		oldFileNodeList := fileNodeList[sdfsFileName]
		for _, node := range oldFileNodeList {
			for _, newNode := range newNodeList {
				if newNode != node {
					oldFileNodeList = append(oldFileNodeList, newNode)
				}
			}
		}
		fileNodeList[sdfsFileName] = oldFileNodeList
	} else {
		// add new record
		fileNodeList[sdfsFileName] = newNodeList
	}
}

// should run continuously
func RemoveFailNode() {
	failNodes := GetFailNodeList()
	for _, node := range failNodes {
		if getAllFile(node) == nil {
			continue
		}
		for _, file := range getAllFile(node) {
			newNodeList := make([]string, 0)
			for _, n := range fileNodeList[file] {
				if n != node {
					newNodeList = append(newNodeList, n)
				}
			}
			fileNodeList[file] = newNodeList
		}
	}
}

// find all files for a given node IP
func getAllFile(nodeIp string) []string {
	files := make([]string, 0)
	for file, nodeList := range fileNodeList {
		for _, node := range nodeList {
			if node == nodeIp {
				files = append(files, file)
			}
		}
	}
	return files
}

// delete record in file-node map
func DeleteFileRecord(sdfsFileName string, nodeIP string) {
	/*if _, exist := fileNodeList[sdfsFileName]; exist {
		delete(fileNodeList, sdfsFileName)
	}*/
	nodeList := fileNodeList[sdfsFileName]
	if nodeList == nil {
		delete(fileNodeList, sdfsFileName)
		logger.PrintInfo("File " + sdfsFileName + " has been deleted!")
	} else {
		for index, node := range nodeList {
			if node == nodeIP {
				nodeList = append(nodeList[:index], nodeList[index+1:]...)
				break
			}
		}
		fileNodeList[sdfsFileName] = nodeList
	}
}

// find nodes to write to or read from
func FindNewNode(sdfsFileName string, sender string) []string {
	// if key not exist in map, it will get nil
	storeList := fileNodeList[sdfsFileName]
	logger.PrintInfo(listToString(storeList) + "   has stored file  " + sdfsFileName)
	nodeNum := config.REPLICA - len(storeList)
	memberIdList := GetMemberIPList()

	ipList := make([]string, 0)
	validIdList := memberIdList

	//logger.PrintInfo(listToString(validIdList) + "   are valid list  ")
	//logger.PrintInfo("senderIP:" + sender)
	//logger.PrintInfo("Length of the initial validlist" + strconv.Itoa(len(validIdList)))
	for index, id := range validIdList {

		if strings.Compare(id, sender) == 0 {
			validIdList = append(validIdList[:index], validIdList[index+1:]...)
			logger.PrintInfo("Length of the modified validlist" + strconv.Itoa(len(validIdList)))
		}
		// if len(storeList) == 0 {
		// 	break
		// }
		for i, n := range storeList {
			if id == n {
				validIdList = append(validIdList[:i], validIdList[i+1:]...)
			}
		}
	}
	// when member node is less than replica
	if len(validIdList) < nodeNum {
		nodeNum = len(validIdList)
	}
	// randomly pick servers in valid nodes to store the connection
	count := 0
	for len(ipList) != nodeNum {
		valid := true
		num := int(config.Hash(sdfsFileName+string(('a'+rune(count))))) % len(validIdList)
		ip := validIdList[num]
		for _, i := range ipList {
			if ip == i {
				valid = false
			}
		}
		if valid {
			ipList = append(ipList, ip)
			logger.PrintInfo("New target has been chosen  " + ip)
		}
		count++
	}

	logger.PrintInfo("Target nodes are  " + listToString(ipList))
	return ipList
}

func listToString(list []string) string {
	var targetString strings.Builder
	for _, e := range list {
		targetString.WriteString(e + " ;  ")
	}
	return targetString.String()
}

// master check whether to replicate files or not---should run continuously
func CheckReplicate() {
	for file, nodeList := range fileNodeList {
		if len(nodeList) < config.REPLICA {
			storeList := fileNodeList[file]
			ipList := FindNewNode(file, introducerIp)
			replicateFile(storeList, ipList, file)
		}
	}
}

// send the replicate request to one existed file node
func replicateFile(storeList []string, newList []string, filename string) {
	// decide which node is the good file
	sourceNode := storeList[0]
	repMessage := &pbm.TCPMessage{
		Type:     pbm.MsgType_PUT_MASTER_REP,
		SenderIP: GetLocalIPAddr().String(),
		PayLoad:  newList,
		FileName: filename,
	}
	msgBytes, _ := EncodeTCPMessage(repMessage)

	SendMessage(sourceNode, msgBytes)

}

// master return target node to write
func PutReplyMessage(remoteMsg *pbm.TCPMessage) {
	// check if key exist in map
	writeList := make([]string, 0)
	if fileNodeList[remoteMsg.FileName] == nil {
		logger.PrintInfo("Find new node")
		writeList = FindNewNode(remoteMsg.FileName, remoteMsg.SenderIP)
	} else {
		writeList = fileNodeList[remoteMsg.FileName]
	}
	repMessage := &pbm.TCPMessage{
		Type:      pbm.MsgType_PUT_MASTER_REP,
		SenderIP:  GetLocalIPAddr().String(),
		PayLoad:   writeList,
		FileName:  remoteMsg.FileName,
		FileSize:  remoteMsg.FileSize,
		LocalPath: remoteMsg.LocalPath,
	}
	msgBytes, _ := EncodeTCPMessage(repMessage)
	SendMessage(remoteMsg.SenderIP, msgBytes)
}

func GetReplyMessage(filename string, sender string) {
	readList := fileNodeList[filename]
	if readList == nil {
		/*todo: deal with non-existed file*/
	}
	repMessage := &pbm.TCPMessage{
		Type:     pbm.MsgType_GET_MASTER_REP,
		SenderIP: GetLocalIPAddr().String(),
		PayLoad:  readList,
		FileName: filename,
	}
	msgBytes, _ := EncodeTCPMessage(repMessage)
	SendMessage(sender, msgBytes)
}

// master return target node with VM ip list that store the file
func ListReplyMessage(filename string, targetIp string) {
	repMessage := &pbm.TCPMessage{
		Type:     pbm.MsgType_LIST_REP,
		FileName: filename,
		SenderIP: GetLocalIPAddr().String(),
		PayLoad:  fileNodeList[filename],
	}
	msgBytes, _ := EncodeTCPMessage(repMessage)
	SendMessage(targetIp, msgBytes)
}

//master send delete request to file node
func DeleteMessage(filename string) {
	ipList := fileNodeList[filename]
	if ipList == nil {
		logger.PrintInfo("No such file in SDFS")
		return
	}

	fileMessage := &pbm.TCPMessage{
		Type:     pbm.MsgType_DELETE,
		SenderIP: GetLocalIPAddr().String(),
		FileName: filename,
	}
	message, _ := EncodeTCPMessage(fileMessage)
	for _, target := range ipList {
		SendMessage(target, message)
	}

}

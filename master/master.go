package master

import (
	pbm "../ProtocolBuffers/MessagePackage"
	"../config"
	"../connection"
	"../detector"
)

var (
	introducerIp string
	//localIp string
	fileList     []string
	fileNodeList map[string][]string
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
	failNodes := detector.GetFailNodeList()
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
func deleteFileRecord(sdfsFileName string) {
	if _, exist := fileNodeList[sdfsFileName]; exist {
		delete(fileNodeList, sdfsFileName)
	}
}

// find nodes to write to or read from
func FindNewNode(sdfsFileName string) []string {
	// if key not exist in map, it will get nil
	storeList := fileNodeList[sdfsFileName]
	nodeNum := config.REPLICA - len(storeList)
	memberIdList := detector.GetMemberIDList()

	ipList := make([]string, 0)
	validIdList := make([]string, 0)
	for _, id := range memberIdList {
		if id == detector.GetLocalIPAddr().String() {
			continue
		}
		for _, n := range storeList {
			if id != n {
				validIdList = append(validIdList, id)
			}
		}
	}
	// randomly pick servers in valid nodes to store the connection
	count := 0
	valid := true
	for len(ipList) != nodeNum {
		num := int(config.Hash(sdfsFileName+string(('a'+rune(count))))) % len(validIdList)
		ip := validIdList[num]
		for _, i := range ipList {
			if ip == i {
				valid = false
			}
		}
		if valid {
			ipList = append(ipList, ip)
		}
		count++
	}
	return ipList
}

// master check whether to replicate files or not---should run continuously
func CheckReplicate() {
	for file, nodeList := range fileNodeList {
		if len(nodeList) < config.REPLICA {
			storeList := fileNodeList[file]
			ipList := FindNewNode(file)
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
		SenderIP: detector.GetLocalIPAddr().String(),
		PayLoad:  newList,
	}
	msgBytes, _ := connection.EncodeTCPMessage(repMessage)
	connection.SendMessage(sourceNode, msgBytes)

}

// master return target node to write
func PutReplyMessage(fileName string, sender string) {
	// check if key exist in map
	writeList := make([]string, 0)
	if fileNodeList[fileName] == nil {
		writeList = FindNewNode(fileName)
	} else {
		writeList = fileNodeList[fileName]
	}
	repMessage := &pbm.TCPMessage{
		Type:     pbm.MsgType_PUT_MASTER_REP,
		SenderIP: detector.GetLocalIPAddr().String(),
		PayLoad:  writeList,
	}
	msgBytes, _ := connection.EncodeTCPMessage(repMessage)
	connection.SendMessage(sender, msgBytes)
}

/*todo: whether to reply one certain node to read? or reply multiple nodes then only connect one?*/
func GetReplyMessage(filename string, sender string) {
	readList := fileNodeList[filename]
	if readList == nil {
		/*todo: deal with non-existed file*/
	}
	repMessage := &pbm.TCPMessage{
		Type:     pbm.MsgType_GET_MASTER_REP,
		SenderIP: detector.GetLocalIPAddr().String(),
		PayLoad:  readList,
	}
	msgBytes, _ := connection.EncodeTCPMessage(repMessage)
	connection.SendMessage(sender, msgBytes)
}

// master return target node with VM ip list that store the file
func ListReplyMessage(filename string, storeipList []string, sender string) {
	repMessage := &pbm.TCPMessage{
		Type:     pbm.MsgType_LIST_REP,
		FileName: filename,
		SenderIP: detector.GetLocalIPAddr().String(),
		PayLoad:  storeipList,
	}
	msgBytes, _ := connection.EncodeTCPMessage(repMessage)
	connection.SendMessage(sender, msgBytes)
}

//master send delete request to file node
func DeleteMessage(filename string, targetIp string){
	fileMessage := &pbm.TCPMessage{
		Type:	  pbm.MsgType_DELETE,
		SenderIP: GetLocalIPAddr().String(),
		FileName: sdfsFileName,
	}
	message, _ := connection.EncodeTCPMessage(fileMessage)
	connection.SendMessage(targetIp, message)
}
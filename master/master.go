package master

import (
	"../config"
	"../detector"

	"../failure"
)

var (
	introducerIp string
	//localIp string
	fileList     []string
	fileNodeList map[string][]string
)

/*todo: apart from the file-node map, send message to nodes*/

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

// remove node from  file-node map
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
func findNewNode(sdfsFileName string) []string {
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

// master check whether to replicate files or not
func CheckReplicate() {
	for file, nodeList := range fileNodeList {
		if len(nodeList) < config.REPLICA {
			failure.ReplicateFile(file)
		}
	}
}

// find the node list for a certain file
func findStoreNode(sdfsFileName string) []string {
	return fileNodeList[sdfsFileName]
}

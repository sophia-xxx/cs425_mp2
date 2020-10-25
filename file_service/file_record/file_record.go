package file_record

import (
	"strconv"
	"strings"

	"cs425_mp2/config"
	"cs425_mp2/member_service"
	"cs425_mp2/util"
	"cs425_mp2/util/logger"
)

type FileName = string
type NodeIP = string

var FileNodeList = make(map[FileName][]NodeIP)

// add or update record in file-node map
func UpdateFileNode(sdfsFileName string, newNodeList []string) {
	// update existed record
	if _, exist := FileNodeList[sdfsFileName]; exist {
		oldFileNodeList := FileNodeList[sdfsFileName]
		for _, node := range oldFileNodeList {
			for _, newNode := range newNodeList {
				if newNode != node {
					oldFileNodeList = append(oldFileNodeList, newNode)
				}
			}
		}
		FileNodeList[sdfsFileName] = oldFileNodeList
	} else {
		// add new record
		FileNodeList[sdfsFileName] = newNodeList
	}
}

// should run continuously
func RemoveFailNode() {
	failNodes := member_service.GetFailNodeList()
	for _, node := range failNodes {
		if FindAllFilesInNode(node) == nil {
			continue
		}
		for _, file := range FindAllFilesInNode(node) {
			newNodeList := make([]string, 0)
			for _, n := range FileNodeList[file] {
				if n != node {
					newNodeList = append(newNodeList, n)
				}
			}
			FileNodeList[file] = newNodeList
		}
	}
}

// find all files for a given node IP
func FindAllFilesInNode(nodeIp NodeIP) []FileName {
	files := make([]string, 0)
	for file, nodeList := range FileNodeList {
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
	nodeList := FileNodeList[sdfsFileName]
	if nodeList == nil {
		delete(FileNodeList, sdfsFileName)
		logger.PrintInfo("File " + sdfsFileName + " has been deleted!")
	} else {
		for index, node := range nodeList {
			if node == nodeIP {
				nodeList = append(nodeList[:index], nodeList[index+1:]...)
				break
			}
		}
		FileNodeList[sdfsFileName] = nodeList
	}
}

// find nodes to write to or read from
func FindNewNode(sdfsFileName string, sender string) []string {
	// if key not exist in map, it will get nil
	storeList := FileNodeList[sdfsFileName]
	if len(storeList) > 0 {
		logger.PrintInfo(util.ListToString(storeList), "has stored file", sdfsFileName)
	}
	nodeNum := config.REPLICA - len(storeList)
	memberIdList := member_service.GetAliveMemberIPList()

	ipList := make([]string, 0)
	validIdList := memberIdList

	//logger.PrintInfo(listToString(validIdList) + "   are valid list  ")
	//logger.PrintInfo("senderIP:" + sender)
	//logger.PrintInfo("Length of the initial validlist" + strconv.Itoa(len(validIdList)))
	for index, id := range validIdList {

		if strings.Compare(id, sender) == 0 {
			validIdList = append(validIdList[:index], validIdList[index+1:]...)
			logger.PrintDebug("Length of the modified validlist", strconv.Itoa(len(validIdList)))
		}

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
			logger.PrintDebug("New target has been chosen", ip)
		}
		count++
	}

	logger.PrintInfo("Chosen hosts to store the file are", util.ListToString(ipList))
	return ipList
}



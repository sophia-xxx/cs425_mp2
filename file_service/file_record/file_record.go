package file_record

import (
	"cs425_mp2/file_service/file_manager"
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

func NewMasterInit() {
	logger.PrintInfo("Initializing file record...")
	FileNodeList = make(map[FileName][]NodeIP)
	RestoreFileNode(
		util.GetLocalIPAddr().String(),
		file_manager.GetLocalSDFSFileList(),
	)
}

// add or update record in file-node map
func UpdateFileNode(sdfsFileName FileName, newNodeList []NodeIP) {
	// update existed record
	oldNodeList, exist := FileNodeList[sdfsFileName]

	if !exist {
		FileNodeList[sdfsFileName] = newNodeList
	} else {
		FileNodeList[sdfsFileName] = util.Merge(oldNodeList, newNodeList)
	}
}

func RestoreFileNode(nodeIP NodeIP, filenames []FileName) {
	for _, filename := range filenames {
		if _, exist := FileNodeList[filename]; !exist {
			FileNodeList[filename] = make([]string, 0)
		}
		FileNodeList[filename] = append(FileNodeList[filename], nodeIP)
	}
}

// should run continuously
func RemoveFailedNodes() {
	failNodes := member_service.GetFailNodeList()
	for _, failedNode := range failNodes {
		filesInNode := FindAllFilesInNode(failedNode)
		// if the failed does not store any file
		if filesInNode == nil {
			continue
		}
		// otherwise, remove the node from its fileNode
		for _, file := range filesInNode {
			logger.PrintInfo("Clearing file record for a failed node...")
			newNodeList := make([]string, 0)
			for _, n := range FileNodeList[file] {
				if strings.Compare(n, failedNode) != 0 {
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
			if strings.Compare(node, nodeIp) == 0 {
				files = append(files, file)
			}
		}
	}
	return files
}

// delete record in file-node map
func DeleteFileInNodeRecord(sdfsFileName string, nodeIP string) {
	nodeList := FileNodeList[sdfsFileName]
	if nodeList == nil {
		delete(FileNodeList, sdfsFileName)
		logger.PrintInfo("File " + sdfsFileName + " has been deleted!")
	} else {
		for index, node := range nodeList {
			if strings.Compare(node, nodeIP) == 0 {
				nodeList = append(nodeList[:index], nodeList[index+1:]...)
				break
			}
		}
		FileNodeList[sdfsFileName] = nodeList
	}
}

// delete a file completely
func DeleteFileInAllNodes(sdfsFileName string) {
	delete(FileNodeList, sdfsFileName)
}

// find new nodes to store replica
func FindNewNode(sdfsFileName string, senderIP string) []string {
	// if key not exist in map, it will get nil
	currStoringNodes := FileNodeList[sdfsFileName]
	if len(currStoringNodes) > 0 {
		logger.PrintInfo(util.ListToString(currStoringNodes), "has stored file", sdfsFileName)
	}
	numNodesToPut := config.REPLICA - len(currStoringNodes)
	memberIdList := member_service.GetAliveMemberIPList()

	ipList := make([]string, 0)
	validIPList := memberIdList

	for index, id := range validIPList {

		if strings.Compare(id, senderIP) == 0 {
			if index==len(validIPList)-1{
				validIPList = validIPList[:index]
			}else{
				validIPList = append(validIPList[:index], validIPList[index+1:]...)
			}
			logger.PrintDebug("Length of the modified validlist", strconv.Itoa(len(validIPList)))
		}

		for i, n := range currStoringNodes {
			if strings.Compare(id, n) == 0 {
				if i == len(validIPList) - 1 {
					validIPList = validIPList[:i]
				}else{
					validIPList = append(validIPList[:i], validIPList[i+1:]...)
				}
			}
		}
	}

	// when member node is less than replica
	if len(validIPList) < numNodesToPut {
		numNodesToPut = len(validIPList)
	}
	// randomly pick servers in valid nodes to store the connection
	count := 0
	for len(ipList) != numNodesToPut {
		valid := true
		num := int(config.Hash(sdfsFileName+string(('a'+rune(count))))) % len(validIPList)
		ip := validIPList[num]
		for _, i := range ipList {
			if strings.Compare(ip, i) == 0 {
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

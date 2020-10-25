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
	FileNodeList = make(map[FileName][]NodeIP)
	RestoreFileNode(
		util.GetLocalIPAddr().String(),
		file_manager.GetLocalSDFSFileList(),
	)
}

// add or update record in file-node map
func UpdateFileNode(sdfsFileName string, newNodeList []string) {
	// update existed record
	if _, exist := FileNodeList[sdfsFileName]; exist {
		oldFileNodeList := FileNodeList[sdfsFileName]
		for _, node := range oldFileNodeList {
			for _, newNode := range newNodeList {
				if strings.Compare(newNode, node) != 0 {
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

func RestoreFileNode(nodeIP NodeIP, filenames []FileName) {
	for _, filename := range filenames {
		if nodelist, exist := FileNodeList[filename]; !exist {
			FileNodeList[filename] = make([]string, 0)
		} else {
			nodelist = append(nodelist, nodeIP)
		}
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
			if strings.Compare(node, nodeIP) == 0 {
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
			if strings.Compare(id, n) == 0 {
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

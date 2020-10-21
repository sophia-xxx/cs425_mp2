package master

import (
	"../config"
	"../detector"
	"hash/fnv"
)

var (
	introducerIp string
	//localIp string
	fileList     []string
	fileNodeList map[string][]string
)

// add or update record in file-node map
func updateFileRecord(sdfsFileName string, newNodeList []string) {
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

// delete record in file-node map
func deleteFileReord(sdfsFileName string) {
	if _, exist := fileNodeList[sdfsFileName]; exist {
		delete(fileNodeList, sdfsFileName)
	}
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
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
	// randomly pick servers in valid nodes to store the file
	count := 0
	valid := true
	for len(ipList) != nodeNum {
		num := int(hash(sdfsFileName+string(('a'+rune(count))))) % len(validIdList)
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

// find the stored list for a certain file
func findStoreNode(sdfsFileName string) []string {
	return fileNodeList[sdfsFileName]
}

package file

import (
	"encoding/json"
	"hash/fnv"
	//"fmt"
	"io"
	"net"
	"os"
	//"fmt"
	"../config"
	"../detector"
	"../logger"
	//"sync"
	//"../networking"
)

var (
	introducerIp string
	//localIp string
	fileList     []string
	fileNodeList map[string][]string
)

type FileMessage struct {
	messageType string
	senderAddr  string
	payload     []byte
}

const (
	MSG_SEARCH = "search"
	MSG_GET    = "get"
	MSG_PUT    = "put"
	MSG_DELTE  = "delete"
	MSG_ACK    = "ack"
)

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// find nodes to write to or read from
func findNode(sdfsFileName string) []string {
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

// server listen to TCP connection request
func listenTCP() {
	addressString := detector.GetLocalIPAddr().String() + config.PORT
	localAddr, err := net.ResolveTCPAddr("tcp4", addressString)
	if err != nil {
		logger.ErrorLogger.Println("Cannot resolve TCP address!")
	}
	listener, err := net.ListenTCP("tcp4", localAddr)
	if err != nil {
		logger.ErrorLogger.Println("Cannot open TCP listener!")
	}

	conn, err := listener.AcceptTCP()
	if err != nil {
		logger.ErrorLogger.Println("Cannot open TCP connection!")
	}
	/* ?????deal with multiple reads, but will have problem with multiple writes*/
	go handleConnection(conn)

}

/*todo: message handler*/
func handleConnection(conn *net.TCPConn) {
	// read message data
	buf := make([]byte, config.BUFFER_SIZE)
	n, err := conn.Read(buf)
	if err != nil {
		logger.ErrorLogger.Println("Unable to read data!")
	}
	messageBytes := buf[0:n]
	if len(messageBytes) == 0 {

	}

}

// deal with "get file" command
func getFileCommand(sdfsFileName string, localFileName string) {
	//send TCP message to master server
	localMessage := FileMessage{
		messageType: "search",
		senderAddr:  detector.GetLocalIPAddr().String(),
	}
	var msgeBytes []byte
	var err error
	if msgeBytes, err = json.Marshal(localMessage); err != nil {
		logger.ErrorLogger.Println("JSON marshal error:", err)
	}
	sendMessage(introducerIp, msgeBytes)

}

// deal with "put file" command
func putFileCommand(localFileName string, sdfsFileName string) {
	/*todo: send message to master server */
	/*todo : send message to data server*/

}

//deal with "delete file" command
func deleteFileCommand(sdfsFileName string) {
	/*todo: send message to master*/
}

// send TCP message
func sendMessage(dest string, message []byte) {
	remoteAddress, _ := net.ResolveTCPAddr("tcp4", dest+config.PORT)
	conn, err := net.DialTCP("tcp4", nil, remoteAddress)
	if err != nil {
		logger.ErrorLogger.Println("Cannot dial remote address!")
	}
	_, err = conn.Write(message)
}

// send file by TCP connection (send filename-->get ACK-->send file)
func sendFile(localFilePath string, dest string, filename string) {
	remoteAddress, _ := net.ResolveTCPAddr("tcp4", dest+config.PORT)
	conn, err := net.DialTCP("tcp4", nil, remoteAddress)
	if err != nil {
		logger.ErrorLogger.Println("Cannot dial remote address!")
	}
	// send filename and wait for reply
	_, err = conn.Write([]byte(filename))
	responseBuf := make([]byte, config.BUFFER_SIZE)
	n, err := conn.Read(responseBuf)
	if err != nil {
		logger.ErrorLogger.Println("Cannot read response")
	}
	if string(responseBuf[:n]) != "ok" {
		logger.ErrorLogger.Println("Cannot set up file transfer connection")
		return
	}

	defer conn.Close()
	// set directory and send file
	fs, err := os.Open(localFilePath)
	defer fs.Close()
	if err != nil {
		logger.ErrorLogger.Println("File path error!")
	}
	buf := make([]byte, config.BUFFER_SIZE)
	for {
		// open file
		n, err := fs.Read(buf)
		if err == io.EOF {
			logger.InfoLogger.Println("Compete file reading!")
			break
		}

		//  send file
		conn.Write(buf[:n])
	}
	return

}

/*todo: how to combine receive file into receive message*/
// receive file by TCP connection
/*func receiveFile(filepath string, conn *net.TCPConn) {
	defer conn.Close()
	// set directory and read file
	//os.Mkdir("./sdfs",0777)
	file, err := os.Create(filepath)
	defer file.Close()
	if err != nil {
		logger.ErrorLogger.Println("Cannot create file!")
	}
	// read data from connection
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			logger.InfoLogger.Println("Compete file reading!")
			break
		}
		file.Write(buf[:n])
	}
	return
}*/

func deleteFile(filename string) {
	os.Remove(filename)

}

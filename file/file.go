package file

import (
	"encoding/json"

	"io"
	"net"
	"os"
	//"fmt"
	"../config"
	"../detector"
	"../logger"
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

// socket to listen TCP message
func ListenMessage() {
	addressString := detector.GetLocalIPAddr().String() + config.TCPPORT
	localAddr, err := net.ResolveTCPAddr("tcp4", addressString)
	if err != nil {
		logger.ErrorLogger.Println("Cannot resolve TCP address!")
	}
	listener, err := net.ListenTCP("tcp4", localAddr)
	if err != nil {
		logger.ErrorLogger.Println("Cannot open TCP listener!")
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			logger.ErrorLogger.Println("Cannot open TCP connection!")
		}

		go handleConnection(conn)
	}

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

// socket to read filename and file
func ListenFile() {
	// open file socket
	addressString := detector.GetLocalIPAddr().String() + config.FILEPORT
	localAddr, err := net.ResolveTCPAddr("tcp4", addressString)
	if err != nil {
		logger.ErrorLogger.Println("Cannot resolve file TCP address!")
	}
	listener, err := net.ListenTCP("tcp4", localAddr)
	if err != nil {
		logger.ErrorLogger.Println("Cannot open TCP listener!")
	}
	conn, err := listener.AcceptTCP()
	if err != nil {
		logger.ErrorLogger.Println("Cannot open TCP connection!")
	}
	defer conn.Close()
	// receive filename and create file
	nameBuf := make([]byte, config.BUFFER_SIZE)
	n, err := conn.Read(nameBuf)
	if err != nil {
		logger.ErrorLogger.Println("Cannot receive filename")
	}
	filename := string(nameBuf[:n])
	logger.InfoLogger.Println("Receive filename")
	if filename != "" {
		_, err = conn.Write([]byte("ACK"))
		if err != nil {
			logger.ErrorLogger.Println("Cannot send ACK")
		}
	}
	// create sdfsfile
	file, err := os.Create("./sdfsFile" + filename)
	defer file.Close()
	if err != nil {
		logger.ErrorLogger.Println("Cannot create file!")
	}
	// read data from connection
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			logger.InfoLogger.Println("Complete file reading!")
			break
		}
		file.Write(buf[:n])
	}
	return
}

// send TCP message
func sendMessage(dest string, message []byte) {
	remoteAddress, _ := net.ResolveTCPAddr("tcp4", dest+config.TCPPORT)
	conn, err := net.DialTCP("tcp4", nil, remoteAddress)
	if err != nil {
		logger.ErrorLogger.Println("Cannot dial remote address!")
	}
	_, err = conn.Write(message)
}

// send file by TCP connection (send filename-->get ACK-->send file)
func sendFile(localFilePath string, dest string, filename string) {
	remoteAddress, _ := net.ResolveTCPAddr("tcp4", dest+config.FILEPORT)
	conn, err := net.DialTCP("tcp4", nil, remoteAddress)
	if err != nil {
		logger.ErrorLogger.Println("Cannot dial remote file socket!")
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

/*
func deleteFile(filename string) {
	os.Remove(filename)

}*/

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

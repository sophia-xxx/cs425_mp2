package detector

import (
	"cs425_mp2/config"
	"cs425_mp2/logger"
	"io"
	"net"
	"os"
	"strconv"
)

// socket to read filename and connection
func ListenFile(filePath string, fileSize int32, isPut bool) {
	// open connection socket
	addressString := GetLocalIPAddr().String() + ":" + config.FILEPORT
	localAddr, err := net.ResolveTCPAddr("tcp4", addressString)
	if err != nil {
		logger.PrintInfo("Cannot resolve connection file address!")
	}
	listener, err := net.ListenTCP("tcp4", localAddr)
	if err != nil {
		logger.PrintInfo("Cannot listen file port!")
	}

	conn, err := listener.AcceptTCP()
	if err != nil {
		logger.PrintInfo("Cannot open file connection!")
	}
	defer conn.Close()

	// receive filename and create connection
	nameBuf := make([]byte, config.BUFFER_SIZE)
	n, err := conn.Read(nameBuf)
	if err != nil {
		logger.PrintInfo("Cannot receive filename")
	}
	filename := string(nameBuf[:n])

	if filename != "" {
		_, err = conn.Write([]byte("ACK"))
		if err != nil {
			logger.PrintInfo("Cannot send ACK")
		}
	}
	logger.PrintInfo("Received filename as: " + filename)
	// create sdfsfile
	file, err := os.Create(filePath)
	defer file.Close()
	if err != nil {
		logger.PrintInfo("Cannot create file!")
	}

	// read data from connection
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		logger.PrintInfo("This time we read:" + strconv.Itoa(n) + " bytes")
		if err == io.EOF {
			logger.PrintInfo("Complete connection reading!")
			break
		}
		file.Write(buf[:n])
	}
	logger.PrintInfo("Finish receiving file!")
	if isPut {
		// finish reading file and check file size, then send ACK
		fileInfo, _ := os.Stat(filePath)
		if int32(fileInfo.Size()) == fileSize {
			SendWriteACK(introducerIp, filename)
		} else {
			logger.PrintInfo("File is broken")
			os.Remove(filePath)
		}
	}

	return

}

// send connection by TCP connection (send filename-->get ACK-->send connection)
func sendFile(localFilePath string, dest string, filename string) {
	remoteAddress, _ := net.ResolveTCPAddr("tcp4", dest+":"+config.FILEPORT)
	conn, err := net.DialTCP("tcp4", nil, remoteAddress)
	if err != nil {
		logger.ErrorLogger.Println("Cannot dial remote connection socket!")
	}
	defer conn.Close()
	// send filename and wait for reply
	logger.PrintInfo("filename is " + filename)
	sendlen, err := conn.Write([]byte(filename))
	if err != nil {
		logger.PrintInfo(err)
	}
	logger.PrintInfo("Send length of " + strconv.Itoa(sendlen) + " filename")

	responseBuf := make([]byte, config.BUFFER_SIZE)
	n, err := conn.Read(responseBuf)
	if err != nil {
		// logger.ErrorLogger.Println("Cannot read response")
		logger.PrintInfo("Cannot read response")
	}
	logger.PrintInfo("Received file ack with length" + strconv.Itoa(n))
	if string(responseBuf[:n]) != "ok" {
		// logger.ErrorLogger.Println("Cannot set up connection transfer connection")
		logger.PrintInfo("Cannot set up connection transfer connection")
		return
	}

	// set directory and send connection
	fs, err := os.Open(localFilePath)
	defer fs.Close()
	if err != nil {
		logger.PrintInfo("File path error!")
	}
	buf := make([]byte, config.BUFFER_SIZE)
	for {
		// open connection
		n, err := fs.Read(buf)
		logger.PrintInfo("This time we write " + strconv.Itoa(n) + " bytes into buffer")
		if err == io.EOF {
			logger.InfoLogger.Println("Compete connection reading!")
			break
		}

		//  send connection
		conn.Write(buf[:n])
	}
	logger.PrintInfo("Finish sending file!")
	return

}

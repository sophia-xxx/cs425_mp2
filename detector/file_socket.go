package detector

import (
	"cs425_mp2/config"
	"cs425_mp2/logger"
	"io"
	"net"
	"os"
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
	logger.PrintInfo("Receive filename  " + filename)
	if filename != "" {
		_, err = conn.Write([]byte("ACK"))
		if err != nil {
			logger.PrintInfo("Cannot send ACK")
		}
	}
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
		if err == io.EOF {
			logger.PrintInfo("Complete connection reading!")
			break
		}
		file.Write(buf[:n])
	}
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
	// send filename and wait for reply
	_, err = conn.Write([]byte(filename))
	responseBuf := make([]byte, config.BUFFER_SIZE)
	n, err := conn.Read(responseBuf)
	if err != nil {
		logger.ErrorLogger.Println("Cannot read response")
	}
	if string(responseBuf[:n]) != "ok" {
		logger.ErrorLogger.Println("Cannot set up connection transfer connection")
		return
	}

	defer conn.Close()
	// set directory and send connection
	fs, err := os.Open(localFilePath)
	defer fs.Close()
	if err != nil {
		logger.ErrorLogger.Println("File path error!")
	}
	buf := make([]byte, config.BUFFER_SIZE)
	for {
		// open connection
		n, err := fs.Read(buf)
		if err == io.EOF {
			logger.InfoLogger.Println("Compete connection reading!")
			break
		}

		//  send connection
		conn.Write(buf[:n])
	}
	return

}

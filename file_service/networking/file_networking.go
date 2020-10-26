package networking

import (
	"cs425_mp2/config"
	"cs425_mp2/member_service"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/gogf/greuse"

	"cs425_mp2/util"
	"cs425_mp2/util/logger"

	"cs425_mp2/file_service/file_record"
)

// socket to read filename and connection
func ListenFile(filePath string, fileSize int32, isPut bool) {
	// open connection socket
	listenAddr := ":" + config.FileTransferPort
	listener, err := greuse.Listen("tcp4", listenAddr)
	if err != nil {
		logger.PrintInfo("Cannot listen file port!")
	}

	conn, err := listener.Accept()
	if err != nil {
		logger.PrintError("Listen file failed because:", err)
		return
	}
	defer conn.Close()

	// receive filename and create connection
	nameBuf := make([]byte, config.BUFFER_SIZE)
	n, err := conn.Read(nameBuf)
	if err != nil {
		logger.PrintError("Cannot receive filename")
	}
	filename := string(nameBuf[:n])

	if filename != "" {
		_, err = conn.Write([]byte("ACK"))
		if err != nil {
			logger.PrintError("Cannot send ACK")
		}
	}
	//logger.PrintInfo("Received filename as: " + filename)
	// create sdfsfile
	file, err := os.Create(filePath)
	defer file.Close()
	if err != nil {
		logger.PrintError("Cannot create file!")
	}

	// read data from connection
	buf := make([]byte, config.BUFFER_SIZE)
	for {
		n, err := conn.Read(buf)
		logger.PrintDebug("This time we read:" + strconv.Itoa(n) + " bytes")
		if err == io.EOF {
			logger.PrintInfo("Complete connection reading!")
			break
		}
		file.Write(buf[:n])
	}
	logger.InfoLogger.Println("Finish receiving file!")
	// write operation will send ACK to client to guarantee quorum write
	if isPut {
		// finish reading file and check file size, then send ACK
		//fileInfo, _ := os.Stat(filePath)

		if strings.Compare(util.GetLocalIPAddr().String(), member_service.GetMasterIP()) == 0 {
			//logger.PrintInfo("Master write file")
			file_record.UpdateFileNode(filename, []string{member_service.GetMasterIP()})
			return
		}
		//if int32(fileInfo.Size()) == fileSize {
		//	SendWriteACK(member_service.GetMasterIP(), filename)
		//} else {
		//	logger.PrintInfo("File is broken")
		//	os.Remove(filePath)
		//}
	}
}

// send connection by TCP connection (send filename-->get ACK-->send connection)
func SendFile(localFilePath string, dest string, filename string) {
	remoteAddress := dest + ":" + config.FileTransferPort
	//localAddr := util.GetLocalIPAddr().String() + ":" + config.MemberServicePort
	localAddr := ":0"
	conn, err := greuse.Dial("tcp4", localAddr, remoteAddress)
	if err != nil {
		logger.PrintError("Send file", filename, "failed because:", err)
		return
	}
	defer conn.Close()
	// send filename and wait for reply
	sendlen, err := conn.Write([]byte(filename))
	if err != nil {
		logger.PrintError("Send file", filename, "failed because:", err)
		return
	}
	logger.PrintDebug("Send length of " + strconv.Itoa(sendlen) + " filename")

	responseBuf := make([]byte, config.BUFFER_SIZE)
	n, err := conn.Read(responseBuf)
	if err != nil {
		logger.PrintError("Cannot read response when sending file", filename)
		return
	}
	if string(responseBuf[:n]) != "ACK" {
		logger.PrintError("Cannot set up connection transfer connection")
		return
	}

	// set directory and send connection
	fs, err := os.Open(localFilePath)
	if err != nil {
		logger.PrintError("File path error!    " + localFilePath)
	}
	defer fs.Close()

	buf := make([]byte, config.BUFFER_SIZE)
	for {
		// open connection
		n, err := fs.Read(buf)
		logger.PrintDebug("This time we write " + strconv.Itoa(n) + " bytes into buffer")
		if err == io.EOF || n == 0 {
			logger.PrintDebug("Complete connection reading!")
			break
		}

		//  send connection
		conn.Write(buf[:n])
	}

	logger.PrintInfo("Successfully sent file", filename, "to", dest, ".")
}
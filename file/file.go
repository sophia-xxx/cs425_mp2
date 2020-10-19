package file

import (
	"io"
	"net"
	"os"
	//"fmt"
	"../config"
	"../logger"
)

func getFileNode(filename string, isMaster bool) {

}

// deal with "get file" command
func getFileCommand(filename string, sdfsFileName string) {
	//filePath:="./sdfs"+sdfsFileName
	//fileNodeList:=getFileNode(sdfsFileName,masterAddr)
	//for _,fileNodeAddr :=range fileNodeList{
	//	receiveFile()
	//}
}
func putFileCommand() {

}

func deleteFileCommand() {

}

// send file by TCP connection
func sendFile(filepath string, dest string) {
	remoteAddress, _ := net.ResolveTCPAddr("tcp4", dest+config.PORT)
	conn, err := net.DialTCP("tcp4", nil, remoteAddress)
	defer conn.Close()
	// set directory and send file
	fs, err := os.Open(filepath)
	defer fs.Close()
	if err != nil {
		logger.ErrorLogger.Println("File path error!")
	}
	buf := make([]byte, 4096)
	for {
		// open file
		n, err1 := fs.Read(buf)
		if err == io.EOF {
			logger.InfoLogger.Println("Compete file reading!")
		}
		if err1 != nil {
			logger.ErrorLogger.Println("Cannot read file!")
		}
		//  send file
		conn.Write(buf[:n])
	}

}

// server process
func receiveFile(filepath string, port string) {
	localAddr, err := net.ResolveTCPAddr("tcp4", port)
	if err != nil {
		logger.ErrorLogger.Println("Cannot resolve TCP address!")
	}
	listener, err := net.ListenTCP("tcp4", localAddr)
	if err != nil {
		logger.ErrorLogger.Println("Cannot open TCP listener!")
	}
	conn, err := listener.Accept()
	if err != nil {
		logger.ErrorLogger.Println("Cannot open TCP connection!")
	}
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
		}
		if err != nil {
			logger.ErrorLogger.Println("Cannot read from buffer!")
		}

		file.Write(buf[:n])

	}
}

func getWriteServer() {

}

func deleteFile(filename string, dirpath string) {
	os.Remove(filename)
	os.RemoveAll(dirpath)
}

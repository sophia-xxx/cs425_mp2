package file

import (
	"io"
	"net"
	"os"
	//"fmt"
	"../config"
	"../detector"
	"../logger"
	//"sync"
)

/*todo: masterNode and dataNode message handler*/
/*todo: master allocate dataNode/ find dataNode*/
/*todo: dataNode read and write sdfsfile*/

// deal with "get file" command
func getFileCommand(sdfsFileName string, localFileName string) {
	/*todo: send message to master server */
	/*todo : send message to data server*/
	localFilePath := "./localfile/" + localFileName
	receiveFile(localFilePath)
}

// deal with "put file" command
func putFileCommand(localFileName string, sdfsFileName string) {
	/*todo: send message to master server */
	/*todo : send message to data server*/
	localFilePath := "./localfile/" + localFileName
	remoteAddr := ""
	sendFile(localFilePath, remoteAddr)

}

//deal with "delete file" command
func deleteFileCommand(sdfsFileName string) {
	/*todo: send message to master*/
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
func receiveFile(filepath string) {
	addressString := detector.GetLocalIPAddr().String() + config.PORT
	localAddr, err := net.ResolveTCPAddr("tcp4", addressString)
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

func deleteFile(filename string) {
	os.Remove(filename)

}

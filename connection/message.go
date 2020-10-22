package connection

import (
	"net"

	//"fmt"
	"../config"
	"../detector"
	"../logger"
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
		logger.ErrorLogger.Println("Cannot open TCP connection!")
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

	}else{
		args := strings.Split(input, " ")
		cmd := args[0]
		param1 := ""
		param2 := ""
		switch cmd{
			case "request_for_put_target":
				
		}
	}

}

// send TCP message
func SendMessage(dest string, message []byte) {
	remoteAddress, _ := net.ResolveTCPAddr("tcp4", dest+config.TCPPORT)
	conn, err := net.DialTCP("tcp4", nil, remoteAddress)
	if err != nil {
		logger.ErrorLogger.Println("Cannot dial remote address!")
	}
	_, err = conn.Write(message)
}

func EncodeFileCommandMessage(fileMessage string) ([]byte, error) {
	message, err := proto.Marshal(fileMessage)
	return message, err
}
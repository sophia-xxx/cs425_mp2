package connection

import (
	"github.com/golang/protobuf/proto"
	"net"
	//"strings"
	mg "../ProtocolBuffers/MessagePackage"
	//"fmt"
	"../config"
	"../detector"
	"../logger"
	"../master"
)

var isMaster bool

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

func handleConnection(conn *net.TCPConn) {
	// read message data
	buf := make([]byte, config.BUFFER_SIZE)
	n, err := conn.Read(buf)
	if err != nil {
		logger.ErrorLogger.Println("Unable to read data!")
	}
	messageBytes := buf[0:n]
	remoteMsg, _ := DecodeFileMessage(messageBytes)
	// master return target node to write and read
	if isMaster && remoteMsg.Type == mg.MsgType_SEARCH {
		master.HandleSearchMessage(remoteMsg.FileInfo)
	}
	// client receive node list of search
	if remoteMsg.Type == mg.MsgType_SEARCHREP {

	}
	// client get node ACK of write (up to 4 ACK, then write sucess)
	if remoteMsg.Type == mg.MsgType_WRITEACK {

	}
	// 	server get replicate message from master, then replicate a certain file
	if remoteMsg.Type == mg.MsgType_REPLICA {

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

func EncodeFileMessage(fileMessage *mg.TCPMessage) ([]byte, error) {
	message, err := proto.Marshal(fileMessage)
	return message, err
}
func DecodeFileMessage(message []byte) (*mg.TCPMessage, error) {
	list := &mg.TCPMessage{}
	err := proto.Unmarshal(message, list)

	return list, err
}
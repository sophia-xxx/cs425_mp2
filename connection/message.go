package connection

import (
	"github.com/golang/protobuf/proto"
	"net"
	//"strings"
	pbm "../ProtocolBuffers/MessagePackage"
	//"fmt"
	"../config"
	"../detector"
	"../logger"
	"../master"
)

var isMaster bool
var quorum int

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
	defer conn.Close()
	// read message data
	buf := make([]byte, config.BUFFER_SIZE)
	n, err := conn.Read(buf)
	if err != nil {
		logger.ErrorLogger.Println("Unable to read data!")
	}
	messageBytes := buf[0:n]
	remoteMsg, _ := DecodeTCPMessage(messageBytes)
	// master return target node to write
	if isMaster && remoteMsg.Type == pbm.MsgType_PUT_MASTER {
		master.PutReplyMessage(remoteMsg.FileName, remoteMsg.SenderIP)
	}
	// master return target node to read
	if isMaster && remoteMsg.Type == pbm.MsgType_GET_MASTER {
		master.GetReplyMessage(remoteMsg.FileName, remoteMsg.SenderIP)
	}
	// client send write file request to target nodes
	if remoteMsg.Type == pbm.MsgType_PUT_MASTER_REP {
		targetList := remoteMsg.PayLoad
		for _, target := range targetList {
			sendWriteReq(target, remoteMsg.FileName)
			sendFile(remoteMsg.LocalPath, target, remoteMsg.FileName)
		}
	}
	// server receive file and save it
	if remoteMsg.Type == pbm.MsgType_PUT_P2P {
		ListenFile(config.SDFS_DIR + remoteMsg.FileName)
	}
	// when write finish, client will receive write ACK to determine write success
	if remoteMsg.Type == pbm.MsgType_WRITE_ACK {
		// quorum determine whether the write is succeed
		if quorum == 4 {
			/*todo: if there is a failed write, how to deal with that?*/
			logger.InfoLogger.Println("Write " + remoteMsg.FileName + " successfully!")
		} else {
			quorum++
		}

	}
	// client send read request to target nodes
	if remoteMsg.Type == pbm.MsgType_GET_MASTER_REP {
		// receive file from target nodes
		targetList := remoteMsg.PayLoad
		for _, target := range targetList {
			sendReadReq(target, remoteMsg.FileName)
			ListenFile(config.LOCAL_DIR + remoteMsg.FileName)
		}
	}
	// server reply to get request and send file to client
	/*todo: how to decide quorum for read??*/
	if remoteMsg.Type == pbm.MsgType_GET_P2P {
		sendFile(config.SDFS_DIR+remoteMsg.FileName, remoteMsg.SenderIP, remoteMsg.FileName)
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

func EncodeTCPMessage(fileMessage *pbm.TCPMessage) ([]byte, error) {
	message, err := proto.Marshal(fileMessage)
	return message, err
}
func DecodeTCPMessage(message []byte) (*pbm.TCPMessage, error) {
	list := &pbm.TCPMessage{}
	err := proto.Unmarshal(message, list)

	return list, err
}

// client send write request to target nodes
func sendWriteReq(targetIp string, sdfsFileName string) {
	fileMessage := &pbm.TCPMessage{
		Type:     pbm.MsgType_PUT_P2P,
		FileName: sdfsFileName,
		SenderIP: detector.GetLocalIPAddr().String(),
	}
	message, _ := EncodeTCPMessage(fileMessage)
	SendMessage(targetIp, message)
}

// client send read request to target node
func sendReadReq(targetIp string, sdfsFileName string) {
	fileMessage := &pbm.TCPMessage{
		Type:     pbm.MsgType_GET_P2P,
		FileName: sdfsFileName,
		SenderIP: detector.GetLocalIPAddr().String(),
	}

	message, _ := EncodeTCPMessage(fileMessage)
	SendMessage(targetIp, message)
}

/*// target nodes reply "ACK" to client's write request
func putFileCommandNodeACK(targetIp string, sdfsFileName string) {
	var fileMessage pbm.TCPMessage      //sch?
	fileMessage.MsgType = "PUT_P2P_ACK" //sch?
	fileMessage.FileInfo = sdfsFileName
	fileMessage.senderIP = detector.GetLocalIPAddr().String()
	//the payload is empty, should we do something to initial it?
	message, _ := connection.EncodeFileMessage(fileMessage)
	connection.SendMessage(targetIp, message)
}

// target nodes reply "ACK" to client's read request
func getFileCommandNodeACK(targetIp string, sdfsFileName string, file_size int) {
	var fileMessage pbm.TCPMessage      //sch?
	fileMessage.MsgType = "GET_P2P_ACK" //sch?
	fileMessage.FileInfo = sdfsFileName
	fileMessage.senderIP = detector.GetLocalIPAddr().String()
	//the payload is empty, should we do something to initial it?
	fileMessage.fileSize = file_size
	message, _ := connection.EncodeFileMessage(fileMessage)
	connection.SendMessage(targetIp, message)
}*/

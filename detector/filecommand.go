package detector

import (
	"../connection"
	pbm "../ProtocolBuffers/MessagePackage" //sch?
)

var introducerIp string

// deal with "put connection to master" command
func putFileCommandMaster(localFileName string, sdfsFileName string) {
	var fileMessage pbm.TCPMessage  //sch?
	fileMessage.MsgType = "PUT_MASTER" //sch?
	fileMessage.FileInfo = sdfsFileName
	fileMessage.senderIP = detector.GetLocalIPAddr().String()
	//the payload is empty, should we do something to initial it?
	message, _ := connection.EncodeFileMessage(fileMessage)
	connection.SendMessage(introducerIp, message)
}

// deal with "put connection to target ip" command
func putFileCommandNode(targetIp string, sdfsFileName string){
	var fileMessage pbm.TCPMessage  //sch?
	fileMessage.MsgType = "PUT_P2P" //sch?
	fileMessage.FileInfo = sdfsFileName
	fileMessage.senderIP = detector.GetLocalIPAddr().String()
	//the payload is empty, should we do something to initial it?
	message, _ := connection.EncodeFileMessage(fileMessage)
	connection.SendMessage(targetIp, message)
}

// deal with "master reply target ip for put query"
func putFileMasterRep(RepIP string, targetList []string){
	var fileMessage pbm.TCPMessage  //sch?
	fileMessage.MsgType = "PUT_MASTER_REP" //sch?
	fileMessage.FileInfo = sdfsFileName
	fileMessage.senderIP = detector.GetLocalIPAddr().String()
	fileMessage.payload = targetList //sch?
	message, _ := connection.EncodeFileMessage(fileMessage)
	connection.SendMessage(RepIp, message)
}

// deal with "target ip tell the put initiator that it received the file header successfully"
func putFileCommandNodeACK(targetIp string, sdfsFileName string){
	var fileMessage pbm.TCPMessage  //sch?
	fileMessage.MsgType = "PUT_P2P_ACK" //sch?
	fileMessage.FileInfo = sdfsFileName
	fileMessage.senderIP = detector.GetLocalIPAddr().String()
	//the payload is empty, should we do something to initial it?
	message, _ := connection.EncodeFileMessage(fileMessage)
	connection.SendMessage(targetIp, message)
}

// deal with "get connection to master" command
func getFileCommandMaster(sdfsFileName string, localFileName string) {
	//send TCP message to master server
	//localMessage := FileMessage{
	//	messageType: "search",
	//	senderAddr:  detector.GetLocalIPAddr().String(),
	//}
	//var msgeBytes []byte
	//var err error
	//if msgeBytes, err = json.Marshal(localMessage); err != nil {
	//	logger.ErrorLogger.Println("JSON marshal error:", err)
	//}
	//sendMessage(introducerIp, msgeBytes)
	var fileMessage pbm.TCPMessage  //sch?
	fileMessage.MsgType = "GET_MASTER" //sch?
	fileMessage.FileInfo = sdfsFileName
	fileMessage.senderIP = detector.GetLocalIPAddr().String()
	//the payload is empty, should we do something to initial it?
	message, _ := connection.EncodeFileMessage(fileMessage)
	connection.SendMessage(targetIp, message)
}

// deal with "master reply target ip for get query"
func getFileMasterRep(RepIP string, targetList []string){
	var fileMessage pbm.TCPMessage  //sch?
	fileMessage.MsgType = "GET_MASTER_REP" //sch?
	fileMessage.FileInfo = sdfsFileName
	fileMessage.senderIP = detector.GetLocalIPAddr().String()
	fileMessage.payload = targetList //sch? we might need only one ip to get the file but here lets just keep it a []string
	message, _ := connection.EncodeFileMessage(fileMessage)
	connection.SendMessage(RepIp, message)
}

// deal with "get connection to target ip" command
func getFileCommandNode(targetIp string, sdfsFileName string){
	var fileMessage pbm.TCPMessage  //sch?
	fileMessage.MsgType = "GET_P2P" //sch?
	fileMessage.FileInfo = sdfsFileName
	fileMessage.senderIP = detector.GetLocalIPAddr().String()
	//the payload is empty, should we do something to initial it?
	message, _ := connection.EncodeFileMessage(fileMessage)
	connection.SendMessage(targetIp, message)
}

// deal with "target ip tell the get initiator that it received the file header successfully and send back the file size"
func getFileCommandNodeACK(targetIp string, sdfsFileName string, file_size int){
	var fileMessage pbm.TCPMessage  //sch?
	fileMessage.MsgType = "GET_P2P_ACK" //sch?
	fileMessage.FileInfo = sdfsFileName
	fileMessage.senderIP = detector.GetLocalIPAddr().String()
	//the payload is empty, should we do something to initial it?
	fileMessage.fileSize = file_size
	message, _ := connection.EncodeFileMessage(fileMessage)
	connection.SendMessage(targetIp, message)
}

// deal with "get initiator tell get source that it received the file size successfully and you may send file now"
func getFileCommandNodeACK(targetIp string, sdfsFileName string){
	var fileMessage pbm.TCPMessage  //sch?
	fileMessage.MsgType = "GET_P2P_SIZE_ACK" //sch?
	fileMessage.FileInfo = sdfsFileName
	fileMessage.senderIP = detector.GetLocalIPAddr().String()
	//the payload is empty, should we do something to initial it?
	message, _ := connection.EncodeFileMessage(fileMessage)
	connection.SendMessage(targetIp, message)
}


//deal with "delete connection" command
func deleteFileCommand(sdfsFileName string) {
	/*todo: send message to master*/
}

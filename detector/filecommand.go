package detector

/// Do I need to import connection here?

// deal with "get connection" command
func getFileCommand(sdfsFileName string, localFileName string) {
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

}

// deal with "put connection" command
func putFileCommand(localFileName string, sdfsFileName string) {
	/*todo: send message to master server */
	/*todo : send message to data server*/
	cmd := "request_for_put_target " + localFileName + " " + sdfsFileName  //Is this ok for go?
	message, _ := connection.EncodeFileCommandMessage(cmd)
	connection.sendMessage(introducerIp, message)
}

//deal with "delete connection" command
func deleteFileCommand(sdfsFileName string) {
	/*todo: send message to master*/
}

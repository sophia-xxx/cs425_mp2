package detector

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

}

//deal with "delete connection" command
func deleteFileCommand(sdfsFileName string) {
	/*todo: send message to master*/
}

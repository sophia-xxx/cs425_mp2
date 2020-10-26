package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"

	"cs425_mp2/command_util"
	"cs425_mp2/config"
	"cs425_mp2/file_service"
	"cs425_mp2/member_service"
	"cs425_mp2/util/logger"
)

// Steps to run:
// go build main.go
// ./main -gossip -introIp=123.123.10.1 (gossip, not introducer)
// ./main -introIp=123.123.10.1 (all to all, not introducer)
// ./main -intro -gossip (gossip, introducer)
func main() {
	isMaster := flag.Bool("master", false, "flag for whether this machine is the master")
	isGossip := flag.Bool("gossip", false, "flag for whether this machine uses gossip heartbeating for dissemination")
	masterIP := flag.String("masterIp", "", "the ip of master to connect to")
	port := flag.Int("port", 0, "the port for the server")
	debugMode := flag.Bool("debug", false, "debug mode")
	normalPrompt := flag.Bool("normalPrompt", false, "use normal prompt instead of go prompt")
	flag.Parse()

	config.DebugMode = *debugMode
	if *port != 0 {
		config.MemberServicePort 	= strconv.Itoa(*port)
		config.FileServicePort 		= strconv.Itoa(*port + 1)
		config.FileTransferPort		= strconv.Itoa(*port + 2)
	}
	if (!*isMaster && *masterIP == "") || (*isMaster && *masterIP != "") {
		logger.PrintError("Machine must either be introducer or have IP address of the introducer to connect to, but not both.\nUse the following flags: -gossip -intro -introIp=<ip>")
		os.Exit(1)
	}

	// run services
	logger.InfoLogger.Println("Starting the application...")
	go member_service.RunService(*isMaster, *isGossip, *masterIP)
	go file_service.RunService()

	// handle user input
	if *normalPrompt {
		handleInputViaNormalPrompt()
	} else {
		handleInputViaGoPrompt()
	}

}

func handleInputViaNormalPrompt() {
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Println("> Please enter a new command in the format of: 'Command Additional_info'")
	for {
		input, _ := inputReader.ReadString('\n')
		inputDispatcher(input)
	}
}

func handleInputViaGoPrompt() {
	completer := func(d prompt.Document) []prompt.Suggest {
		s := []prompt.Suggest{
			{Text: "join"},
			{Text: "display"},
			{Text: "switch"},
			{Text: "leave"},

			{Text: "put"},
			{Text: "get"},
			{Text: "delete"},
			{Text: "ls"},
			{Text: "store"},

			{Text: "exit"},
		}
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), false)
	}

	for {
		input := prompt.Input(">>> ", completer)
		inputDispatcher(input)
	}

}

func inputDispatcher(input string) {
	input = strings.TrimSpace(input)
	if input == "exit" {
		os.Exit(0)
	}
	inputs := strings.Split(input, " ")
	command := command_util.Command{
		Method: inputs[0],
		Params: inputs[1:],
	}
	if command_util.IsFileCommand(command) {
		file_service.HandleCommand(command)
	} else if command_util.IsMemberCommand(command) {
		member_service.HandleCommand(command)
	} else {
		fmt.Println("invalid command_util:", input)
	}

}
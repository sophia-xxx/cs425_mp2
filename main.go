package main

import (
	"./detector"
	"./logger"
	"flag"
	"os"
)

// Steps to run:
// go build main.go
// ./main -gossip -introIp=123.123.10.1 (gossip, not introducer)
// ./main -introIp=123.123.10.1 (all to all, not introducer)
// ./main -intro -gossip (gossip, introducer)
func main() {
	isIntroducer := flag.Bool("intro", false, "flag for whether this machine is the introducer")
	isGossip := flag.Bool("gossip", false, "flag for whether this machine uses gossip heartbeating for dissemination")
	introducerIP := flag.String("introIp", "", "the string of the introducer to connect to")
	flag.Parse()

	if (!*isIntroducer && *introducerIP == "") || (*isIntroducer && *introducerIP != "") {
		logger.PrintError("Machine must either be introducer or have IP address of the introducer to connect to, but not both.\nUse the following flags: -gossip -intro -introIp=<ip>")
		os.Exit(1)
	}

	logger.InfoLogger.Println("Starting the application...")
	detector.Run(*isIntroducer, *isGossip, *introducerIP)
}

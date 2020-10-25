package util

import (
	"cs425_mp2/util/logger"
	"net"
	"strings"
)

func ListToString(list []string) string {
	var targetString strings.Builder
	for _, e := range list {
		targetString.WriteString(e + " ;  ")
	}
	return targetString.String()
}

func GetLocalIPAddr() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		logger.PrintError("net.Dial")
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}


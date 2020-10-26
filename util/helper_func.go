/*
This package contains some general utilities
*/
package util

import (
	"cs425_mp2/util/logger"
	"net"
	"strings"
)

func Contains(element string, list []string) bool {
	for _, item := range list {
		if item == element {
			return true
		}
	}
	return false
}

func Merge(list1, list2 []string) []string {
	result := make([]string, len(list1))
	copy(result, list1)
	for _, item := range list2 {
		if !Contains(item, result) {
			result = append(result, item)
		}
	}
	return result
}


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


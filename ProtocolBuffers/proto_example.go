package main

import (
	"fmt"

	"github.com/golang/protobuf/ptypes"

	pb "./ProtoPackage"
)

func main() {
	member := pb.Member{
		HeartbeatCounter: 1,
		LastSeen:         ptypes.TimestampNow(),
	}

	memberList := pb.MemberList{}
	memberList.Dict = make(map[string]*pb.Member)

	machineID := "example_ip___timestamp"
	memberList.Dict[machineID] = &member

	member.HeartbeatCounter++

	fmt.Println(memberList.Dict[machineID])
}

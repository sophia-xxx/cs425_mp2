package membership_test

import (
	"../config"
	"strconv"
	"testing"

	pb "../ProtocolBuffers/ProtoPackage"
	"../membership"
	"github.com/golang/protobuf/ptypes"
	"github.com/jinzhu/copier"
)

var (
	memberA = pb.Member{
		HeartbeatCounter: 1,
		LastSeen:         ptypes.TimestampNow(),
	}

	memberB = pb.Member{
		HeartbeatCounter: 4,
		LastSeen:         ptypes.TimestampNow(),
	}

	memberC = pb.Member{
		HeartbeatCounter: 10,
		LastSeen:         ptypes.TimestampNow(),
	}

	machineIDA = "example_ipA:timestamp"
	machineIDB = "example_ipB:timestamp"
	machineIDC = "example_ipC:timestamp"
)

func TestAddAndRemoveMember(t *testing.T) {
	memberListA := pb.MemberList{}
	memberListA.Dict = make(map[string]*pb.Member)

	membership.AddMemberToMembershipList(&memberListA, machineIDA, &memberA)

	if memberListA.Dict[machineIDA] != &memberA {
		t.Errorf("Expected machine with id \"%s\" to be successfully added to membership list", machineIDA)
	}

	membership.RemoveMemberFromMembershipList(&memberListA, machineIDA)

	if _, ok := memberListA.Dict[machineIDA]; ok {
		t.Errorf("Expected machine with id \"%s\" to have been removed from member list", machineIDA)
	}
}

func TestMergeWithNewRow(t *testing.T) {
	memberListA := pb.MemberList{}
	memberListA.Dict = make(map[string]*pb.Member)

	membership.AddMemberToMembershipList(&memberListA, machineIDA, &memberA)
	membership.AddMemberToMembershipList(&memberListA, machineIDC, &memberC)

	memberListB := pb.MemberList{}
	memberListB.Dict = make(map[string]*pb.Member)

	membership.AddMemberToMembershipList(&memberListB, machineIDA, &memberA)
	membership.AddMemberToMembershipList(&memberListB, machineIDB, &memberB)
	membership.AddMemberToMembershipList(&memberListB, machineIDC, &memberC)

	membership.MergeMembershipLists(&memberListA, &memberListB, nil)

	if _, ok := memberListA.Dict[machineIDB]; !ok {
		t.Errorf("Expected machine with id \"%s\" to exist in member list", machineIDB)
	}
}

func TestMergeWithUpdatedRow(t *testing.T) {
	memberListA := pb.MemberList{}
	memberListA.Dict = make(map[string]*pb.Member)

	membership.AddMemberToMembershipList(&memberListA, machineIDA, &memberA)

	memberListB := pb.MemberList{}
	memberListB.Dict = make(map[string]*pb.Member)

	updatedMemberA := pb.Member{}
	copier.Copy(&updatedMemberA, &memberA)
	updatedMemberA.HeartbeatCounter += 3
	updatedMemberA.LastSeen = ptypes.TimestampNow()

	membership.AddMemberToMembershipList(&memberListB, machineIDA, &updatedMemberA)

	membership.MergeMembershipLists(&memberListA, &memberListB, nil)

	if _, ok := memberListA.Dict[machineIDA]; !ok {
		t.Errorf("Expected machine with id \"%s\" to exist in member list", machineIDA)
	}

	if memberListA.Dict[machineIDA].HeartbeatCounter != memberListB.Dict[machineIDA].HeartbeatCounter {
		t.Errorf("Expected local list to be updated with new heartbeat count of value %d", updatedMemberA.HeartbeatCounter)
	}

	localTime, err1 := ptypes.Timestamp(memberListA.Dict[machineIDA].LastSeen)
	remoteTime, err2 := ptypes.Timestamp(memberListB.Dict[machineIDA].LastSeen)

	if err1 != nil || err2 != nil || !localTime.After(remoteTime) {
		t.Errorf("Expected local list to have higher LastSeen timestamp than remote after merge")
	}
}

func TestGetMemberShipListIPs(t *testing.T) {
	memberList := pb.MemberList{}
	memberList.Dict = make(map[string]*pb.Member)

	membership.AddMemberToMembershipList(&memberList, machineIDA, &memberA)
	membership.AddMemberToMembershipList(&memberList, machineIDB, &memberB)
	membership.AddMemberToMembershipList(&memberList, machineIDC, &memberC)

	ips := membership.GetOtherMembershipListIPs(&memberList, "")
	expectedIps := map[string]bool{"example_ipA": true, "example_ipB": true, "example_ipC": true}

	for _, ip := range ips {
		if _, ok := expectedIps[ip]; !ok {
			t.Errorf("Expected IP: %s to have been extracted from membership list.", ip)
		}
	}
}

func CheckAndRemoveMembershipListFailures(t *testing.T) {
	memberList := pb.MemberList{}
	memberList.Dict = make(map[string]*pb.Member)

	failedList := map[string]bool{"example_ipA": false, "example_ipB": false, "example_ipC": false}

	// Should not be marked failed
	machineIDA := "example_ipA:" + ptypes.TimestampNow().String()

	// Should be marked as failed
	machineIDB := "example_ipB:" +
		strconv.FormatInt(ptypes.TimestampNow().GetSeconds()-config.T_TIMEOUT, 10)

	// Should be marked as failed and removed from both membership list and failed list
	machineIDC := "example_ipC:" +
		strconv.FormatInt(ptypes.TimestampNow().GetSeconds()-config.T_TIMEOUT-config.T_CLEANUP, 10)

	membership.AddMemberToMembershipList(&memberList, machineIDA, &memberA)
	membership.AddMemberToMembershipList(&memberList, machineIDB, &memberB)
	membership.AddMemberToMembershipList(&memberList, machineIDC, &memberC)

	membership.CheckAndRemoveMembershipListFailures(&memberList, &failedList)

	// Machine A
	isFailedA, ok := failedList[machineIDA]
	if !ok {
		t.Errorf("Expected machine with id \"%s\" to exist in failed list", machineIDA)
	} else if isFailedA {
		t.Errorf("Expected machine with id \"%s\" to not be marked failed in failed list", machineIDA)
	}

	if _, ok := memberList.Dict[machineIDA]; !ok {
		t.Errorf("Expected machine with id \"%s\" to exist in member list", machineIDA)
	}

	// Machine B
	isFailedB, ok := failedList[machineIDB]

	if !ok {
		t.Errorf("Expected machine with id \"%s\" to exist in failed list", machineIDB)
	} else if !isFailedB {
		t.Errorf("Expected machine with id \"%s\" to be marked failed in failed list", machineIDB)
	}

	if _, ok := memberList.Dict[machineIDB]; !ok {
		t.Errorf("Expected machine with id \"%s\" to exist in member list", machineIDB)
	}

	// Machine C
	if _, ok := failedList[machineIDC]; ok {
		t.Errorf("Expected machine with id \"%s\" to not exist in failed list", machineIDC)
	}

	if _, ok := memberList.Dict[machineIDA]; ok {
		t.Errorf("Expected machine with id \"%s\" to not exist in member list", machineIDC)
	}
}

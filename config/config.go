package config

import "hash/fnv"

var DebugMode = false

const (
	MemberServicePort 	= "7234"
	FileTransferPort 	= "7235"
	FileServicePort 	= "7236"
)


const BUFFER_SIZE int = 32768

const (
	T_TIMEOUT = 5
	T_CLEANUP = 5
	WaitTimeForElection = 5
)

const STRAT_GOSSIP = "gossip"

const STRAT_ALL = "all"

const PULSE_TIME = 500

const GOSSIP_FANOUT = 4

const REPLICA = 4

const SDFS_DIR = "./sdfsFiles/"

const LOCAL_DIR = "./localFiles/"

const PERM_MODE = 0777

const ACK_TIMEOUT = 1

const PUT = 4

const GET = 9

const DELETE = 12

func Hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

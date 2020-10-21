package config

import "hash/fnv"

const PORT string = "6789"

const FILEPORT string = "8001"

const TCPPORT string = "8002"

const BUFFER_SIZE int = 32768

const T_TIMEOUT = 2

const T_CLEANUP = 2

const STRAT_GOSSIP = "gossip"

const STRAT_ALL = "all"

const PULSE_TIME = 500

const GOSSIP_FANOUT = 4

const REPLICA = 4

func Hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

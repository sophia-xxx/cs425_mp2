# CS 425 MP1 - Distributed Group Membership

### Prerequisites for compilation
- The Go programming language: https://golang.org/
- Google Protocol Buffers (for Go): https://developers.google.com/protocol-buffers/docs/gotutorial

### Compilation & Running
1. Clone the repo
2. Copy the MP2 folder out to `/home/username/go/src` (by mkdir go and mkdir src)
3. Change setup.sh content into your username, in MP2 folder run `bash setup.sh`
4. Compile the code with `go build main.go`, which will produce the executable `./main`
5. Run the code with the following commands:
    - `./main -gossip -introIp=123.123.10.1` (gossip strategy, join group with introducer 123.123.10.1)
    - `./main -introIp=123.123.10.1` (all-to-all strategy, join group with introducer 123.123.10.1)
    - `./main -intro -gossip` (gossip strategy, this machine is the introducer)

6. Then use put, get, ls, store... command, with the local file folder.

Alternatively, you can run the code without building an exectuable using `go run main.go [args]`, such as `go run main.go -gossip -introIp=123.123.10.1`

We have included an executable file (named `./vm_main`) for Linux AMD64 machines, which can be run out of the box without installing any of the prerequisites.
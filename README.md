# CS 425 MP1 - Distributed Group Membership

### Prerequisites for compilation
- The Go programming language: https://golang.org/
- Google Protocol Buffers (for Go): https://developers.google.com/protocol-buffers/docs/gotutorial

### Compilation & Running
1. Clone the repo 
2. Copy and paste mp2 folder into /home/$USER/go/src (by mkdir go and src)
3. run `bash setup.sh` in mp2 folder
4. Compile the code with `go build main.go`, which will produce the executable `./main`
5. Run the code with the following commands:
    - `./main -gossip -introIp=123.123.10.1` (gossip strategy, join group with introducer 123.123.10.1)
    - `./main -introIp=123.123.10.1` (all-to-all strategy, join group with introducer 123.123.10.1)
    - `./main -intro -gossip` (gossip strategy, this machine is the introducer)
6. Run all mp2 command

Alternatively, you can run the code without building an exectuable using `go run main.go [args]`, such as `go run main.go -gossip -introIp=123.123.10.1`

We have included an executable file (named `./vm_main`) for Linux AMD64 machines, which can be run out of the box without installing any of the prerequisites.
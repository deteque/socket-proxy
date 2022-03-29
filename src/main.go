package main

import (
	"flag"
	"os"
	"fmt"
	"os/signal"
	"syscall"
)

const (
	VERSION = "socket-relay VERSION: 20220308.1"
	RETRY_DELAY = 5 //seconds
	KILOBYTE = 1024
	BUFFER_SIZE = 128
)

type arrayFlags []string
var destSocket arrayFlags
var sourceSocketDec *string

func (i *arrayFlags) String() string {
        return "destination socket"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}


func main() {
	sourceSocketDec = flag.String("s", "", "Location of the socket to read dnstap messages from")
	flag.Var(&destSocket, "d", "Destination socket locations. Declare this option for each destination socket. There is no hard limit to the number of destination sockets that can be set.")
	version := flag.Bool("v", false, "Check the version of the program")
	flag.Parse()

	if *version {
		fmt.Println(VERSION)
		return
	}

	if *sourceSocketDec == "" {
		fmt.Println("Source Socket not set")
		return
	}
	if len(destSocket) == 0 {
		fmt.Println("Destination Socket not set")
		return
	}

	signal.Ignore(os.Signal(syscall.SIGHUP))

	run()

}

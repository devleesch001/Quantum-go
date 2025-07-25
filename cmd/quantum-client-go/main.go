package main

import (
	"flag"
	"fmt"
	"github.com/devleesch001/Quantum-go/game/client"
	"github.com/devleesch001/Quantum-go/tools"
	"log/slog"
	"net"
)

var flagPort int
var flagDebug bool

func main() {
	debugFlag := flag.Bool("debug", false, "Activer le mode debug")
	flag.Parse()

	if *debugFlag {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	args := flag.Args()
	if len(args) == 2 {
		name := args[0]
		address := args[1]

		host, port, err := tools.ParseAddressWithDefault(address, "127.0.0.1", tools.DefaultPort)
		if err != nil {
			panic(err)
		}

		client.Run(net.JoinHostPort(host, port), name)

		panic("Client not implemented yet")
	}

	fmt.Println("Usage: ./quantum -server <port> (default 18467)")
}

package main

import (
	"flag"
	"fmt"
	"github.com/devleesch001/Quantum-go/game/server"
	"github.com/devleesch001/Quantum-go/tools"
	"log/slog"
	"net"
)

var flagPort int

func main() {
	debugFlag := flag.Bool("debug", false, "Activer le mode debug")
	serverFlag := flag.Bool("server", false, "Présence du flag serveur")
	flag.Parse()

	if *debugFlag {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	if *serverFlag {
		// Regarder dans flag.Args() s'il y a un argument après le flag -server
		var host = "0.0.0.0"
		var port = "18467"
		var err error

		args := flag.Args()
		if len(args) > 0 {
			host, port, err = tools.ParseAddress(args[0])
			if err != nil {
				slog.Error("Impsible de parser l'addresse", "error", err)
				panic(err)
			}
		}

		server.Run(net.JoinHostPort(host, port))
	}

	fmt.Println("Usage: ./quantum -server <port> (default 18467)")
}

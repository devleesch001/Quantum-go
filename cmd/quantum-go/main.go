package main

import (
	"flag"
	"fmt"
	"github.com/devleesch001/Quantum-go/game/client"
	"github.com/devleesch001/Quantum-go/game/server"
	"github.com/devleesch001/Quantum-go/tools"
	"log/slog"
	"net"
)

func main() {
	address := ""
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
	} else {
		args := flag.Args()
		if len(args) == 2 {
			name := args[0]
			address = args[1]

			host, port, err := tools.ParseAddressWithDefault(address, "127.0.0.1", tools.DefaultPort)
			if err != nil {
				panic(err)
			}

			client.Run(net.JoinHostPort(host, port), name)
		}
	}

	fmt.Println("Usage: quantum <player_name> <server_ip> or quantum -server")
}

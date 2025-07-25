package main

import (
	"flag"
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
		args := flag.Args()
		if len(args) > 0 {
			address = args[0]
		} else {
			address = "0.0.0.0:18467"
		}

		host, port, err := tools.ParseAddress(address)
		if err != nil {
			slog.Error("Impsible de parser l'addresse", "error", err)
			panic(err)
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

			panic("Client not implemented yet")
		}
	}

}

package main

import (
	"flag"
	"fmt"
	"github.com/devleesch001/Quantum-go/game"
	"log/slog"
)

var flagPort int
var flagDebug bool

func main() {

	flag.BoolVar(&flagDebug, "debug", false, "Enable debug logging")
	flag.IntVar(&flagPort, "server", 18467, "Start server with <port>")

	flag.Parse()

	if flagDebug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	if flagPort > 0 && flagPort < 65536 {
		Server(uint16(flagPort))
	}

	fmt.Println("Usage: ./quantum -server <port> (default 18467)")
}

func Server(port uint16) {
	slog.Info("Starting server...")
	slog.Debug("Debug logging enabled")

	var g = new(game.Game)
	defer g.Close()

	if err := g.Start(port); err != nil {
		panic(err)
	}

	slog.Debug(g.String())

	select {}
}

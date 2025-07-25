package server

import (
	"github.com/devleesch001/Quantum-go/game"
	"log/slog"
)

func Run(addr string) {
	slog.Info("Starting server...")
	slog.Debug("Debug logging enabled")

	var g = game.New()
	defer g.Close()
	if err := g.Start(addr); err != nil {
		panic(err)
	}

	slog.Debug(g.String())

	select {}
}

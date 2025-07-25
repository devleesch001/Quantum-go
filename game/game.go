package game

import (
	"fmt"
	"github.com/devleesch001/Quantum-go/maps"
)

type Game struct {
	maps maps.Maps
}

func New() (*Game, error) {
	g := &Game{}

	m, err := maps.New()
	if err != nil {
		return nil, err
	}

	g.maps = *m

	return g, nil
}

func (g Game) String() string {
	return fmt.Sprintf("maps: %+v", g.maps)
}

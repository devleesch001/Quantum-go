package client

import (
	"fmt"
	"github.com/devleesch001/Quantum-go/frames"
	"github.com/devleesch001/Quantum-go/game"
	"github.com/devleesch001/Quantum-go/maps"
	"net"
	"os"
)

type Party struct {
	Game *game.Game // Reference to the game instance
	Maps *maps.Maps
}

func Run(addr, name string) {
	fmt.Printf("Connecting to %s with name %s\n", addr, name)
	fmt.Println("Client not implemented yet")
	os.Exit(0)

	/*	g, err := game.New()
		if err != nil {
			panic(err)
		}

		m, err := maps.New()
		if err != nil {
			panic(err)
		}

		p := &Party{
			Game: g,
			Maps: m,
		}*/

}

type Client struct {
	ID     uint8  // Unique identifier for the client
	Name   string // Name of the client, can be used for display purposes
	Color  byte
	X      uint8
	Y      uint8
	FaceID uint8
	BodyID uint8
	LegsID uint8
	Conn   net.Conn
}

func (c *Client) Label() string {
	return fmt.Sprintf("[Client %d]", c.ID)
}

func (c *Client) String() string {
	return fmt.Sprintf(
		"color: %d, x: %d, y: %d, faceID: %d, bodyID: %d, legsID: %d, conn: %s",
		c.Color, c.X, c.Y, c.FaceID, c.BodyID, c.LegsID, c.Conn.RemoteAddr().String(),
	)
}

func (c *Client) DisconnectFrame() ([]byte, error) {
	return frames.New(c.ID, frames.NewClientInfo(c.X, c.Y, c.Color, c.FaceID, c.BodyID, c.LegsID, c.Name, false)).MarshalBinary()
}

func New(id uint8, conn net.Conn) *Client {
	var c = &Client{
		ID:     id,
		Color:  0,
		X:      1,
		Y:      1,
		FaceID: 0,
		BodyID: 0,
		LegsID: 0,
		Conn:   conn,
	}

	return c
}

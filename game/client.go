package game

import (
	"fmt"
	"github.com/devleesch001/Quantum-go/frames"
	"net"
)

type client struct {
	id     uint8  // Unique identifier for the client
	name   string // Name of the client, can be used for display purposes
	color  byte
	x      uint8
	y      uint8
	faceID uint8
	bodyID uint8
	legsID uint8
	conn   net.Conn
}

func (c *client) Label() string {
	return fmt.Sprintf("[Client %d]", c.id)
}

func (c *client) String() string {
	return fmt.Sprintf(
		"color: %d, x: %d, y: %d, faceID: %d, bodyID: %d, legsID: %d, conn: %s",
		c.color, c.x, c.y, c.faceID, c.bodyID, c.legsID, c.conn.RemoteAddr().String(),
	)
}

func (c *client) DisconnectFrame() ([]byte, error) {
	return frames.New(c.id, frames.NewClientInfo(c.x, c.y, c.color, c.faceID, c.bodyID, c.legsID, c.name, false)).MarshalBinary()
}

func newClient(id uint8, conn net.Conn) *client {
	var c = &client{
		id:     id,
		color:  0,
		x:      1,
		y:      1,
		faceID: 0,
		bodyID: 0,
		legsID: 0,
		conn:   conn,
	}

	return c
}

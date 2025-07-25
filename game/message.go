package game

import "fmt"

type dataMessage struct {
	clientID uint8
	data     []byte
}

func (m dataMessage) String() string {
	return fmt.Sprintf("clientID: %d, Data: %x", m.clientID, m.data)
}

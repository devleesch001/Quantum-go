package game

import "fmt"

type DataMessage struct {
	clientID uint8
	data     []byte
}

func (d DataMessage) ClientID() uint8 {
	return d.clientID
}

func NewDataMessage(clientID uint8, data []byte) DataMessage {
	return DataMessage{
		clientID: clientID,
		data:     data,
	}
}

func (d DataMessage) String() string {
	return fmt.Sprintf("clientID: %d, Data: % 02x", d.clientID, d.data)
}

func (d DataMessage) Data() []byte {
	return d.data
}

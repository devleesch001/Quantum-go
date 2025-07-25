package frames

import "fmt"

type ClientMessage struct {
	text string
}

func (c *ClientMessage) UnmarshalBinary(data []byte) error {
	if len(data) < 1 {
		return ErrInvalidMessageLength
	}
	c.text = string(data[0:])
	return nil
}

func (c *ClientMessage) MarshalBinary() ([]byte, error) {
	data := make([]byte, len(c.text)+1)
	copy(data[0:], c.text)
	return data, nil
}

func (c *ClientMessage) Code() Code {
	return B_MESSAGE
}

func (c ClientMessage) String() string {
	return fmt.Sprintf("%s", c.text)
}

func (cm *ClientMessage) Text() string {
	return cm.text
}

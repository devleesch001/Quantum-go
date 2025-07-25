package frames

import "fmt"

//

// Payload 05 E1 02 05 4F
// 05 - Length of the message
// E1 - Frame type (0x1F)
// 02 - Client ID
// 05 - X coordinate
// 4F - Y coordinate

type ClientPosition struct {
	x uint8
	y uint8
}

func (c *ClientPosition) UnmarshalBinary(data []byte) error {
	if len(data) < 2 {
		return ErrInvalidMessageLength
	}
	c.x = data[0]
	c.y = data[1]
	return nil
}

func (c *ClientPosition) MarshalBinary() ([]byte, error) {
	data := make([]byte, 2)
	data[0] = c.x
	data[1] = c.y
	return data, nil
}

func (c *ClientPosition) Code() Code {
	return B_POS
}

func (c ClientPosition) String() string {
	return fmt.Sprintf("Pos: %d, %d", c.x, c.y)
}

func (c *ClientPosition) X() uint8 {
	return c.x
}

func (c *ClientPosition) Y() uint8 {
	return c.y
}

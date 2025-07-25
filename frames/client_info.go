package frames

import "fmt"

// Payload 0C 1F 03 01 0E 0F 4F 00 00 00 32 32
// 0C - Length of the message
// 1F - Frame type (0x1F)
// 03 - Client ID
// 01 - unknown
// 0E - X coordinate
// 0F - Y coordinate
// 4F - Color
// 00 - Face ID
// 00 - Body ID
// 06 - Legs ID

type ClientInfo struct {
	x      uint8
	y      uint8
	status uint8 // 0x01 for online, 0x00 for offline
	color  byte
	faceID uint8
	bodyID uint8
	legsID uint8
	name   string
}

func NewClientInfo(x, y uint8, color byte, faceID, bodyID, legsID uint8, name string, isOnline bool) *ClientInfo {
	c := &ClientInfo{
		x:      x,
		y:      y,
		color:  color,
		faceID: faceID,
		bodyID: bodyID,
		legsID: legsID,
		name:   name,
	}

	if isOnline {
		c.status = 0x01 // Online
	}

	return c
}

func (c ClientInfo) Code() Code {
	return B_CLIENT_INFOS
}

func (c ClientInfo) String() string {
	return fmt.Sprintf("ClientInfo<%d,%d,%d,%d>", c.x, c.y, c.color, c.faceID)
}

func (c *ClientInfo) MarshalBinary() ([]byte, error) {
	data := make([]byte, 0, 7+len(c.name))
	data = append(data, c.status) // Magic byte
	data = append(data, c.x)
	data = append(data, c.y)
	data = append(data, c.color)
	data = append(data, c.faceID)
	data = append(data, c.bodyID)
	data = append(data, c.legsID)
	data = append(data, []byte(c.name)...)

	return data, nil
}

func (c *ClientInfo) UnmarshalBinary(data []byte) error {
	if len(data) < 8 {
		return ErrInvalidDataLength
	}

	c.x = data[1]
	c.y = data[2]
	c.color = data[3]
	c.faceID = data[4]
	c.bodyID = data[5]
	c.legsID = data[6]

	if len(data) > 7 {
		c.name = string(data[7:])
	} else {
		c.name = ""
	}

	return nil
}

package frames

import "fmt"

type ClientHello struct {
	color  byte
	x      uint8
	y      uint8
	faceID uint8
	bodyID uint8
	legsID uint8
	name   string
}

func (m *ClientHello) UnmarshalBinary(data []byte) error {
	if len(data) < 9 {
		return ErrInvalidMessageLength
	}
	m.color = data[1]

	m.x = data[2]
	m.y = data[3]

	m.faceID = data[4]
	m.bodyID = data[5]
	m.legsID = data[6]

	m.name = string(data[7:])
	return nil
}

func (m *ClientHello) MarshalBinary() ([]byte, error) {
	var buf = make([]byte, 0, 9+len(m.name))
	buf = append(buf, 0x00)
	buf = append(buf, m.color)
	buf = append(buf, m.x)
	buf = append(buf, m.y)

	buf = append(buf, m.faceID)
	buf = append(buf, m.bodyID)
	buf = append(buf, m.legsID)

	buf = append(buf, m.name...)

	return buf, nil
}

func (m ClientHello) Code() byte {
	return B_NEW_CLIENT
}

func (m ClientHello) String() string {
	return fmt.Sprintf("Name: %s, Color %d, x %d, y %d", m.name, m.color, m.x, m.y)
}

/////////
// Getter
/////////

func (m *ClientHello) Y() uint8 {
	return m.y
}

func (m *ClientHello) X() uint8 {
	return m.x
}

func (m *ClientHello) Color() byte {
	return m.color
}

func (m *ClientHello) FaceID() uint8 {
	return m.faceID
}

func (m *ClientHello) BodyID() uint8 {
	return m.bodyID
}

func (m *ClientHello) LegsID() uint8 {
	return m.legsID
}

func (m ClientHello) Name() string {
	return m.name
}

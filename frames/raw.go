package frames

import "fmt"

func UnmarshalRaw(data []byte) (Raw, error) {
	var r Raw
	return r, r.UnmarshalBinary(data)
}

func MarshalRaw(code Code, clientID uint8, payload []byte) ([]byte, error) {
	r := Raw{
		code:     code,
		clientID: clientID,
		payload:  payload,
	}
	return r.MarshalBinary()
}

type Raw struct {
	code     Code
	clientID uint8
	payload  []byte
}

func (r *Raw) SetClientID(id uint8) {
	r.clientID = id
}

func (r *Raw) UnmarshalBinary(data []byte) error {
	lenData := len(data)
	if lenData < 3 {
		return fmt.Errorf("data too short to unmarshal message")
	}

	lenPayload := data[0]

	r.code = data[1]
	r.clientID = data[2]
	r.payload = data[3:lenPayload]
	return nil
}

func (r *Raw) MarshalBinary() ([]byte, error) {
	data := make([]byte, 3+len(r.payload))
	data[0] = byte(len(r.payload) + 3) // total length including code and length byte
	data[1] = r.code
	data[2] = r.clientID
	copy(data[3:], r.payload)

	return data, nil
}

func (r Raw) String() string {
	return fmt.Sprintf("Code: %02x, payload : %x", r.Code(), r.Payload())
}

func (r Raw) Bytes() []byte {
	return append([]byte{r.code}, r.payload...)
}

func (r Raw) Code() Code {
	return r.code
}

func (r Raw) Payload() []byte {
	return r.payload
}

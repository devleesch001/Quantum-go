package frames

type IPayload interface {
	UnmarshalBinary(data []byte) error
	MarshalBinary() ([]byte, error)
	String() string
	Code() Code
}

type Frame struct {
	clientID uint8
	IPayload
}

func New(clientID uint8, p IPayload) *Frame {
	return &Frame{
		clientID,
		p,
	}
}

func UnmarshalFrame(data []byte) (f Frame, err error) {
	return f, f.UnmarshalBinary(data)
}

func (f *Frame) MarshalBinary() ([]byte, error) {

	payload, err := f.IPayload.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return MarshalRaw(f.IPayload.Code(), f.clientID, payload)
}

func (f *Frame) UnmarshalBinary(data []byte) error {
	rawFrame, err := UnmarshalRaw(data)
	if err != nil {
		return err
	}

	switch rawFrame.Code() {
	case B_NEW_CLIENT:
		var hello ClientHello

		if err := hello.UnmarshalBinary(rawFrame.Payload()); err != nil {
			return err
		}

		f.IPayload = &hello

	case B_POS:
		var pos ClientPosition
		if err := pos.UnmarshalBinary(rawFrame.Payload()); err != nil {
			return err
		}

		f.IPayload = &pos
	case B_MESSAGE:
		var msg ClientMessage
		if err := msg.UnmarshalBinary(rawFrame.Payload()); err != nil {
			return err
		}

		f.IPayload = &msg
	default:
		return ErrInvalidFrameType
	}

	return nil
}

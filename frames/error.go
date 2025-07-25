package frames

import "errors"

var ErrInvalidMessageLength = errors.New("invalid message length")
var ErrInvalidDataLength = errors.New("invalid data length")
var ErrInvalidFrameType = errors.New("invalid frame type")

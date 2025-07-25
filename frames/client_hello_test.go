package frames

import (
	"reflect"
	"testing"
)

func TestClientHello(t *testing.T) {
	var clientHello = ClientHello{
		color:  1,
		x:      1,
		y:      1,
		faceID: 1,
		bodyID: 1,
		legsID: 1,
		name:   "test",
	}

	var rawFrame []byte
	t.Run("MarshalBinary", func(t *testing.T) {
		binary, err := clientHello.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}

		rawFrame = binary
		t.Logf("binary: % 02x", binary)
	})

	var newClientHello ClientHello
	t.Run("UnmarshalBinary", func(t *testing.T) {

		err := newClientHello.UnmarshalBinary(rawFrame)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("CheckValues", func(t *testing.T) {
		if !reflect.DeepEqual(clientHello, newClientHello) {
			t.Errorf("Expected: %+v, got: %+v", clientHello, newClientHello)
		}
	})
}

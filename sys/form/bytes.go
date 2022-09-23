package form

import (
	"context"
	"encoding/base64"
	"fmt"
)

// TODO: handle []byte transparently in the de/skeletization?
type Bytes []byte

func (x Bytes) Skeletize(context.Context) any {
	return EncodeBytesToString(x)
}

func (x *Bytes) DeSkeletize(_ context.Context, from any) error {
	s, ok := from.(string)
	if !ok {
		return fmt.Errorf("bytes are represented as a string")
	}
	b, err := DecodeBytesFromString(s)
	if err != nil {
		return err
	}
	*x = b
	return nil
}

func EncodeBytesToString(buf []byte) string {
	return base64.StdEncoding.EncodeToString(buf)
}

func DecodeBytesFromString(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

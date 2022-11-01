package form

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
)

type Bytes []byte

func (x Bytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(EncodeBytesToString(x))
}

func (x *Bytes) UnmarshalJSON(d []byte) error {
	var s string
	if err := json.Unmarshal(d, &s); err != nil {
		return err
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

func StringHashForFilename(s string) string {
	h := sha256.New()
	if _, err := h.Write([]byte(s)); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

package form

import (
	"encoding/base64"
	"encoding/json"
	"os"
)

func EncodeBytesToString(buf []byte) string {
	return base64.StdEncoding.EncodeToString(buf)
}

func DecodeBytesFromString(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func EncodeForm(form any) ([]byte, error) {
	return json.MarshalIndent(form, "", "   ")
}

func DecodeForm(data []byte, form any) error {
	return json.Unmarshal(data, form)
}

func EncodeFormToFile(form any, filepath string) error {
	data, err := EncodeForm(form)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath, data, 0644)
}

func DecodeFormFromFile(filepath string, form any) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	return DecodeForm(data, form)
}

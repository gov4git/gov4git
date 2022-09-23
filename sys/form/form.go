package form

import (
	"context"
	"encoding/json"
	"os"
)

type Form interface {
	Skeletize(context.Context) any
	DeSkeletize(context.Context, any) error
}

func EncodeForm(ctx context.Context, form Form) ([]byte, error) {
	return json.MarshalIndent(form.Skeletize(ctx), "", "   ")
}

func DecodeForm(ctx context.Context, data []byte, form Form) error {
	var skel any
	if err := json.Unmarshal(data, &skel); err != nil {
		return err
	}
	return form.DeSkeletize(ctx, skel)
}

func EncodeFormToFile(ctx context.Context, form Form, filepath string) error {
	data, err := EncodeForm(ctx, form)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath, data, 0644)
}

func DecodeFormFromFile(ctx context.Context, filepath string, form Form) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	return DecodeForm(ctx, data, form)
}

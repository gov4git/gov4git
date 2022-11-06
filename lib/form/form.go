package form

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/go-git/go-billy/v5"
	"github.com/gov4git/gov4git/lib/must"
)

type Form interface{}

func Pretty(form Form) string {
	data, err := json.MarshalIndent(form, "", "   ")
	if err != nil {
		panic(err)
	}
	return string(data)
}

func Encode[F Form](ctx context.Context, w io.Writer, f F) error {
	return json.NewEncoder(w).Encode(f)
}

func Decode[F Form](ctx context.Context, r io.Reader) (form F, err error) {
	err = json.NewDecoder(r).Decode(&form)
	return form, err
}

func EncodeBytes[F Form](ctx context.Context, form F) ([]byte, error) {
	return json.MarshalIndent(form, "", "   ")
}

func DecodeBytes[F Form](ctx context.Context, data []byte) (form F, err error) {
	err = json.Unmarshal(data, &form)
	return form, err
}

func EncodeToFile[F Form](ctx context.Context, fs billy.Filesystem, path string, form F) error {
	file, err := fs.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	return Encode(ctx, file, form)
}

func DecodeFromFile[F Form](ctx context.Context, fs billy.Filesystem, path string) (form F, err error) {
	file, err := fs.Open(path)
	if err != nil {
		return form, err
	}
	defer file.Close()
	return Decode[F](ctx, file)
}

func ToFile[F Form](ctx context.Context, fs billy.Filesystem, path string, form F) {
	if err := EncodeToFile(ctx, fs, path, form); err != nil {
		must.Panic(ctx, err)
	}
}

func FromFile[F Form](ctx context.Context, fs billy.Filesystem, path string) F {
	f, err := DecodeFromFile[F](ctx, fs, path)
	if err != nil {
		must.Panic(ctx, err)
	}
	return f
}

package git

import (
	"bytes"
	"context"
	"strings"

	"github.com/gov4git/gov4git/lib/form"
)

func PrepareCommitMsg(ctx context.Context, human string, embed form.Form) (string, error) {
	var w bytes.Buffer
	human = strings.TrimRight(human, "\n\r\t ")
	if strings.Index(human, "\n\n\n") >= 0 {
		panic("commit messages cannot contain triple new line")
	}
	w.WriteString(human)
	w.WriteString("\n\n\n")
	data, err := form.EncodeForm(ctx, embed)
	if err != nil {
		return "", err
	}
	w.Write(data)
	return w.String(), nil
}

package soul

import (
	"context"

	"github.com/petar/gitsoc/forms"
)

type Soul struct {
	Address forms.SoulAddress
}

func (x Soul) Init(ctx context.Context) (forms.PrivateInfo, error) {
	XXX
}

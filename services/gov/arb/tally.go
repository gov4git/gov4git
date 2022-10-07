package arb

import (
	"context"

	"github.com/gov4git/gov4git/lib/form"
)

type TallyIn struct {
	ReferendumBranch string `json:"referendum_branch"`
}

type TallyOut struct {
	ReferendumRepo   string `json:"referendum_repo"`
	ReferendumBranch string `json:"referendum_branch"`
}

func (x TallyOut) Human(ctx context.Context) string {
	data, err := form.EncodeForm(ctx, x)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func (x GovArbService) Tally(ctx context.Context, in *TallyIn) (*TallyOut, error) {
	panic("tally not implemented")
}

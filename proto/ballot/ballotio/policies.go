package ballotio

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotpolicies/sv"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/mod"
	"github.com/gov4git/lib4git/git"
)

var policyRegistry = mod.NewModuleRegistry[ballotproto.PolicyName, ballotproto.Policy]()

const (
	QVPolicyName ballotproto.PolicyName = "qv"
)

func init() {
	ctx := context.Background()
	Install(
		ctx,
		QVPolicyName,
		sv.SV{
			Kernel: sv.QVScoreKernel{},
		},
	)
}

func Install(ctx context.Context, name ballotproto.PolicyName, policy ballotproto.Policy) {
	policyRegistry.Set(ctx, name, policy)
}

func LoadPolicy(
	ctx context.Context,
	t *git.Tree,
	id ballotproto.BallotID,

) (ballotproto.Ad, ballotproto.Policy) {

	ad := git.FromFile[ballotproto.Ad](ctx, t, id.AdNS())
	return ad, policyRegistry.Get(ctx, ad.Policy)
}

func LookupPolicy(
	ctx context.Context,
	id ballotproto.PolicyName,

) ballotproto.Policy {

	return policyRegistry.Get(ctx, id)
}

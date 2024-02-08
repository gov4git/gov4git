package ballotio

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotpolicies/sv"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/mod"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
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
			Kernel: sv.MakeQVScoreKernel(ctx, 1.0),
		},
	)
}

func Install(ctx context.Context, name ballotproto.PolicyName, policy ballotproto.Policy) {
	policyRegistry.Set(ctx, name, policy)
}

func LoadAd_Local(
	ctx context.Context,
	t *git.Tree,
	id ballotproto.BallotID,

) ballotproto.Ad {

	return git.FromFile[ballotproto.Ad](ctx, t, id.AdNS())
}

func LoadAdPolicy_Local(
	ctx context.Context,
	t *git.Tree,
	id ballotproto.BallotID,

) (ballotproto.Ad, ballotproto.Policy) {

	ad := LoadAd_Local(ctx, t, id)

	p, err := must.Try1[ballotproto.Policy](
		func() ballotproto.Policy {
			return policyRegistry.Get(ctx, ad.Policy)
		},
	)
	must.Assertf(ctx, err == nil, "ballot policy not supported") // ERR

	return ad, p
}

func LookupPolicy(
	ctx context.Context,
	id ballotproto.PolicyName,

) ballotproto.Policy {

	return policyRegistry.Get(ctx, id)
}

func TryLookupPolicy(
	ctx context.Context,
	id ballotproto.PolicyName,

) ballotproto.Policy {

	p, _ := must.Try1[ballotproto.Policy](
		func() ballotproto.Policy {
			return policyRegistry.Get(ctx, id)
		},
	)
	return p
}

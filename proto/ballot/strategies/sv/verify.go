package sv

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/ballot/common"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/id"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/lib4git/must"
)

func (qv SV) VerifyElections(
	ctx context.Context,
	voterAddr id.OwnerAddress,
	govAddr gov.Address,
	voterCloned id.OwnerCloned,
	govCloned gov.Cloned,
	ad *common.Advertisement,
	prior *common.Tally,
	elections common.Elections,
) {

	voterCred := id.GetPublicCredentials(ctx, voterCloned.Public.Tree())
	user := member.LookupUserByID_Local(ctx, govCloned, voterCred.ID)
	if len(user) == 0 {
		must.Errorf(ctx, "cannot find user with id %v in the community", voterCred.ID)
	}

	// tally writes to the gov repo, but the repo is throw-away and won't be committed
	qv.tally(ctx, govCloned, ad, prior, map[member.User]common.Elections{user[0]: elections}, true)
}

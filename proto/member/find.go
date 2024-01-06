package member

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/id"
	"github.com/gov4git/lib4git/must"
)

func FindClonedUser_Local(
	ctx context.Context,
	cloned gov.Cloned,
	userCloned id.OwnerCloned,

) User {

	voterCred := id.GetPublicCredentials(ctx, userCloned.Public.Tree())
	users := LookupUserByID_Local(ctx, cloned, voterCred.ID)
	must.Assertf(ctx, len(users) > 0, "user not found in community")
	return users[0]
}

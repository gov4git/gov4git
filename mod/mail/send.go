package mail

import (
	"context"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/mod/id"
)

type SeqNo int64

func Send[M form.Form](
	ctx context.Context,
	sender *git.Tree,
	receiverAddr git.Address,
	topic string,
	msg M,
) git.Change[SeqNo] {

	receiverCred := id.FetchPublicCredentials(ctx, id.Public(receiverAddr))
	topicNS := SendTopicNS(receiverCred.ID, topic)
	git.TreeMkdirAll(ctx, sender, topicNS)

	panic("XXX")
}

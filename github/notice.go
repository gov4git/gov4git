package github

import (
	"context"

	"github.com/google/go-github/v55/github"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/notice"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/ns"
)

func FlushNotices_Local(
	ctx context.Context,
	cloned gov.Cloned,
	filepath ns.NS,
	ghc *github.Client,
	repo Repo,
	issueNum int,

) git.Change[*notice.NoticeQueue, form.None] {

	queue := notice.LoadNoticeQueue_Local(ctx, cloned, filepath)

	for _, nstate := range queue.NoticeStates {
		if nstate.IsDisplayed() {
			continue
		}
		payload := "### 🏛️ Gov4Git notice\n\n" + nstate.Notice.Body
		replyToIssue(ctx, repo, ghc, issueNum, payload)
		nstate.SetDisplayed()
	}

	return notice.SaveNoticeQueue_StageOnly(ctx, cloned, filepath, queue)
}

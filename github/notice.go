package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v55/github"
	"github.com/gov4git/gov4git/proto/docket/ops"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/notice"
	"github.com/gov4git/lib4git/base"
)

func DisplayNotices_StageOnly(
	ctx context.Context,
	repo Repo,
	ghc *github.Client,
	cloned gov.Cloned,
) {

	t := cloned.Tree()
	motions := ops.ListMotions_Local(ctx, t)
	for _, motion := range motions {
		issueNum, err := MotionIDToIssueNumber(motion.ID)
		if err != nil {
			base.Errorf("encountered motion %v whose id cannot be converted into a github issue number", motion.ID)
			continue
		}
		queue := ops.LoadMotionNotices_Local(ctx, cloned, motion.ID)
		flushNotices(ctx, repo, ghc, cloned, queue, issueNum)
		ops.SaveMotionNotices_StageOnly(ctx, cloned, motion.ID, queue)
	}
}

func flushNotices(
	ctx context.Context,
	repo Repo,
	ghc *github.Client,
	cloned gov.Cloned,
	queue *notice.NoticeQueue,
	issueNum int,
) {

	for _, nstate := range queue.NoticeStates {

		// check if notice already displayed, based on governance records
		if nstate.IsShown() {
			continue
		}

		// TODO: check if notice already displayed, according to github

		payload := fmt.Sprintf("### 🏛️ Gov4Git notice `%v`\n\n", nstate.ID) + nstate.Notice.Body
		replyToIssue(ctx, repo, ghc, issueNum, payload)
		nstate.MarkShown()
	}
}

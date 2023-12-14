package github

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/google/go-github/v55/github"
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/docket/ops"
	"github.com/gov4git/gov4git/proto/docket/policies/pmp/concern"
	"github.com/gov4git/gov4git/proto/docket/policies/pmp/proposal"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/gov4git/proto/notice"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func SyncManagedIssues(
	ctx context.Context,
	repo Repo,
	githubClient *github.Client,
	govAddr gov.OwnerAddress,
) git.Change[form.Map, *SyncManagedChanges] {

	govCloned := gov.CloneOwner(ctx, govAddr)
	syncChanges := SyncManagedIssues_StageOnly(ctx, repo, githubClient, govAddr, govCloned)
	chg := git.NewChange[form.Map, *SyncManagedChanges](
		fmt.Sprintf("Sync %d managed GitHub issues", len(syncChanges.IssuesCausingChange)),
		"github_sync",
		form.Map{},
		syncChanges,
		nil,
	)
	return proto.CommitIfChanged(ctx, govCloned.Public, chg)
}

type SyncManagedChanges struct {
	IssuesCausingChange ImportedIssues     `json:"issues_causing_change"`
	Updated             schema.MotionIDSet `json:"updated_motions"`
	Opened              schema.MotionIDSet `json:"opened_motions"`
	Closed              schema.MotionIDSet `json:"closed_motions"`
	Cancelled           schema.MotionIDSet `json:"cancelled_motions"`
	Froze               schema.MotionIDSet `json:"froze_motions"`
	Unfroze             schema.MotionIDSet `json:"unfroze_motions"`
	AddedRefs           schema.RefSet      `json:"added_refs"`
	RemovedRefs         schema.RefSet      `json:"removed_refs"`
}

func newSyncManagedChanges() *SyncManagedChanges {
	return &SyncManagedChanges{
		IssuesCausingChange: nil,
		Updated:             schema.MotionIDSet{},
		Opened:              schema.MotionIDSet{},
		Closed:              schema.MotionIDSet{},
		Cancelled:           schema.MotionIDSet{},
		Froze:               schema.MotionIDSet{},
		Unfroze:             schema.MotionIDSet{},
		AddedRefs:           schema.RefSet{},
		RemovedRefs:         schema.RefSet{},
	}
}

func SyncManagedIssues_StageOnly(
	ctx context.Context,
	repo Repo,
	githubClient *github.Client,
	govAddr gov.OwnerAddress,
	cloned gov.OwnerCloned,
) (syncChanges *SyncManagedChanges) {

	syncChanges = newSyncManagedChanges()

	t := cloned.Public.Tree()

	// load github issues and governance motions, and
	// index them under a common key space
	motions := indexMotions(ops.ListMotions_Local(ctx, t))
	loadPR := func(ctx context.Context,
		repo Repo,
		issue *github.Issue,
	) bool {

		id := IssueNumberToMotionID(int64(issue.GetNumber()))
		m, motionExists := motions[id]

		return IsIssueManaged(issue) && // merged state not relevant if issue is not managed
			issue.GetState() == "closed" && // merged state is not relevant for open prs
			(!motionExists || !m.Closed) // merged state is relevant, when no corresponding motion exists or motion is open
	}

	_, issues := LoadIssues(ctx, githubClient, repo, loadPR)

	// ensure every issue has a corresponding up-to-date motion
	for key, issue := range issues {
		id := schema.MotionID(key)
		if issue.Managed {
			if motion, ok := motions[id]; ok { // if motion for issue already exists, update it
				changed := syncMeta(ctx, cloned, syncChanges, issue, motion)
				switch {

				case issue.Closed && motion.Closed:

				case issue.Closed && !motion.Closed:
					if motion.IsConcern() {
						// manually closing an issue motion cancels it
						ops.CancelMotion_StageOnly(ctx, cloned, id)
						syncChanges.Cancelled.Add(id)
					} else if motion.IsProposal() {
						// manually closing a proposal motion closes it
						if issue.Merged {
							ops.CloseMotion_StageOnly(ctx, cloned, id, schema.Accept)
						} else {
							ops.CloseMotion_StageOnly(ctx, cloned, id, schema.Reject)
						}
						syncChanges.Closed.Add(id)
					} else {
						must.Errorf(ctx, "motion is neither a concern nor a proposal")
					}
					changed = true

				case !issue.Closed && motion.Closed:
					base.Infof("GitHub issue %v has been re-opened; corresonding motion remains closed", issue.Number)
					ops.AppendMotionNotices_StageOnly(
						ctx,
						cloned.PublicClone(),
						id,
						notice.Noticef(
							ctx,
							"Reopening %s %s [#%v](%v) does not reopen the corresponding motion. Consider creating a new %s instead.",
							motion.GithubArticle(),
							motion.GithubType(),
							id,
							motion.TrackerURL,
							motion.GithubType(),
						),
					)

				case !issue.Closed && !motion.Closed:

				}
				if changed {
					syncChanges.IssuesCausingChange = append(syncChanges.IssuesCausingChange, issue)
				}

			} else { // otherwise, no motion for this issue exists, so create one

				if !issue.Closed {
					syncCreateMotionForIssue(ctx, govAddr, cloned, syncChanges, issue, id)
					syncChanges.IssuesCausingChange = append(syncChanges.IssuesCausingChange, issue)
				}

			}

		} else { // issue is not governed, freeze motion if it exists and is open

			if motion, ok := motions[id]; ok { // motion for issue already exists, update it
				// if motion closed, do nothing
				// if motion frozen, do nothing
				// otherwise, freeze motion
				if !motion.Closed && !motion.Frozen {
					ops.AppendMotionNotices_StageOnly(
						ctx,
						cloned.PublicClone(),
						id,
						notice.Noticef(ctx, "The Gov4Git motion for this no longer managed issue/PR has been frozen."),
					)
					ops.FreezeMotion_StageOnly(notice.Mute(ctx), cloned, id)
					syncChanges.Froze.Add(id)
					syncChanges.IssuesCausingChange = append(syncChanges.IssuesCausingChange, issue)
				}
			}

		}
	}

	// don't touch motions that have no corresponding issue

	// update references
	matchingMotions := indexMotions(ops.ListMotions_Local(ctx, t))
	syncRefs(ctx, cloned, syncChanges, issues, matchingMotions)

	syncChanges.IssuesCausingChange.Sort()
	return
}

func syncMeta(
	ctx context.Context,
	cloned gov.OwnerCloned,
	chg *SyncManagedChanges,
	issue ImportedIssue,
	motion schema.Motion,

) bool {

	if motion.Closed {
		return false
	}

	author := findMemberForGithubLogin(ctx, cloned.PublicClone(), issue.Author)
	if motion.TrackerURL == issue.URL &&
		motion.Author == author &&
		motion.Title == issue.Title &&
		motion.Body == issue.Body &&
		slices.Equal(motion.Labels, issue.Labels) {
		return false
	}
	ops.EditMotionMeta_StageOnly(
		ctx,
		cloned,
		motion.ID,
		author,
		issue.Title,
		issue.Body,
		issue.URL,
		issue.Labels,
	)
	chg.Updated.Add(motion.ID)
	return true
}

const (
	MotionPolicyForIssue = concern.ConcernPolicyName
	MotionPolicyForPR    = proposal.ProposalPolicyName
)

func motionPolicyForIssue(issue ImportedIssue) schema.PolicyName {
	if issue.PullRequest {
		return MotionPolicyForPR
	}
	return MotionPolicyForIssue
}

func syncCreateMotionForIssue(
	ctx context.Context,
	addr gov.OwnerAddress,
	cloned gov.OwnerCloned,
	chg *SyncManagedChanges,
	issue ImportedIssue,
	id schema.MotionID,
) {

	must.Assertf(ctx, !issue.Closed, "issue is closed")

	ops.OpenMotion_StageOnly(
		ctx,
		cloned,
		id,
		issue.MotionType(),
		motionPolicyForIssue(issue),
		findMemberForGithubLogin(ctx, cloned.PublicClone(), issue.Author),
		issue.Title,
		issue.Body,
		issue.URL,
		issue.Labels,
	)
	chg.Opened.Add(id)
}

// findMemberForGithubLogin returns the community user corresponding to a GitHub login.
// If there is no corresponding community member, an empty string user is returned.
func findMemberForGithubLogin(ctx context.Context, cloned gov.Cloned, login string) member.User {

	var user member.User
	query := member.User(strings.ToLower(login))
	if member.IsUser_Local(ctx, cloned, query) {
		user = query
	}
	return user
}

func indexMotions(ms schema.Motions) map[schema.MotionID]schema.Motion {
	x := map[schema.MotionID]schema.Motion{}
	for _, m := range ms {
		x[m.ID] = m
	}
	return x
}

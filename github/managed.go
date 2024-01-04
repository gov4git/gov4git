package github

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/google/go-github/v55/github"
	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/gov4git/v2/proto/motion"
	"github.com/gov4git/gov4git/v2/proto/motion/motionapi"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp/concern"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp/proposal"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/gov4git/v2/proto/notice"
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
	IssuesCausingChange ImportedIssues          `json:"issues_causing_change"`
	Updated             motionproto.MotionIDSet `json:"updated_motions"`
	Opened              motionproto.MotionIDSet `json:"opened_motions"`
	Closed              motionproto.MotionIDSet `json:"closed_motions"`
	Cancelled           motionproto.MotionIDSet `json:"cancelled_motions"`
	Froze               motionproto.MotionIDSet `json:"froze_motions"`
	Unfroze             motionproto.MotionIDSet `json:"unfroze_motions"`
	AddedRefs           motionproto.RefSet      `json:"added_refs"`
	RemovedRefs         motionproto.RefSet      `json:"removed_refs"`
}

func newSyncManagedChanges() *SyncManagedChanges {
	return &SyncManagedChanges{
		IssuesCausingChange: nil,
		Updated:             motionproto.MotionIDSet{},
		Opened:              motionproto.MotionIDSet{},
		Closed:              motionproto.MotionIDSet{},
		Cancelled:           motionproto.MotionIDSet{},
		Froze:               motionproto.MotionIDSet{},
		Unfroze:             motionproto.MotionIDSet{},
		AddedRefs:           motionproto.RefSet{},
		RemovedRefs:         motionproto.RefSet{},
	}
}

func SyncManagedIssues_StageOnly(
	ctx context.Context,
	repo Repo,
	ghc *github.Client,
	addr gov.OwnerAddress,
	cloned gov.OwnerCloned,

) (syncChanges *SyncManagedChanges) {

	syncChanges = newSyncManagedChanges()

	t := cloned.Public.Tree()

	// load github issues and governance motions, and
	// index them under a common key space

	ms0 := indexMotions(motionapi.ListMotions_Local(ctx, t))
	loadPR := func(ctx context.Context,
		repo Repo,
		issue *github.Issue,
	) bool {

		id := IssueNumberToMotionID(int64(issue.GetNumber()))
		m, motionExists := ms0[id]

		return IsIssueManaged(issue) && // merged state not relevant if issue is not managed
			issue.GetState() == "closed" && // merged state is not relevant for open prs
			(!motionExists || !m.Closed) // merged state is relevant, when no corresponding motion exists or motion is open
	}

	_, issues := LoadIssues(ctx, ghc, repo, loadPR)

	// sync motions with issues
	syncMotionsWithIssues(
		ctx,
		repo,
		ghc,
		addr,
		cloned,
		syncChanges,
		ms0,
		issues,
	)

	// update references
	ms0Refs := indexMotions(motionapi.ListMotions_Local(ctx, t))
	syncRefs(ctx, cloned, syncChanges, issues, ms0Refs)

	// motions have changed
	ms1 := indexMotions(motionapi.ListMotions_Local(ctx, t))

	// resync with GitHub to flush changes to GitHub issues
	syncMotionsWithIssues(
		ctx,
		repo,
		ghc,
		addr,
		cloned,
		syncChanges,
		ms1,
		issues,
	)

	syncChanges.IssuesCausingChange.Sort()
	return
}

func syncMotionsWithIssues(
	ctx context.Context,
	repo Repo,
	ghc *github.Client,
	addr gov.OwnerAddress,
	cloned gov.OwnerCloned,
	syncChanges *SyncManagedChanges,
	motions map[motionproto.MotionID]motionproto.Motion,
	issues map[string]ImportedIssue,
) {

	for key, issue := range issues {
		id := motionproto.MotionID(key)
		syncMotion(
			ctx,
			repo,
			ghc,
			addr,
			cloned,
			syncChanges,
			motions,
			id,
			issue,
		)
	}
	// don't touch motions that have no corresponding issue
}

func syncMotion(
	ctx context.Context,
	repo Repo,
	ghc *github.Client,
	addr gov.OwnerAddress,
	cloned gov.OwnerCloned,
	syncChanges *SyncManagedChanges,
	motions map[motionproto.MotionID]motionproto.Motion,
	id motionproto.MotionID,
	issue ImportedIssue,
) {

	if issue.Managed {
		if motion, ok := motions[id]; ok { // if motion for issue already exists, update it
			changed := syncMeta(ctx, cloned, syncChanges, issue, motion)
			switch {

			case issue.Closed && motion.Closed:

			case issue.Closed && !motion.Closed:
				if motion.IsConcern() {
					// manually closing an issue motion cancels it
					motionapi.CancelMotion_StageOnly(ctx, cloned, id)
					syncChanges.Cancelled.Add(id)
				} else if motion.IsProposal() {
					// manually closing a proposal motion closes it
					if issue.Merged {
						motionapi.CloseMotion_StageOnly(ctx, cloned, id, motionproto.Accept)
					} else {
						motionapi.CloseMotion_StageOnly(ctx, cloned, id, motionproto.Reject)
					}
					syncChanges.Closed.Add(id)
				} else {
					must.Errorf(ctx, "motion is neither a concern nor a proposal")
				}
				changed = true

			case !issue.Closed && motion.Closed:

				err := must.Try(
					func() {
						closeIssue(ctx, repo, ghc, int(issue.Number))
					},
				)
				if err != nil {
					base.Infof("GitHub %s %v is open, while corresonding motion is closed. Failed to close GitHub issue (%v)",
						motion.GithubType(), issue.Number, err)
					motionapi.AppendMotionNotices_StageOnly(
						ctx,
						cloned.PublicClone(),
						id,
						notice.Noticef(
							ctx,
							"This %s must now be closed, as the corresponding Gov4Git motion has closed. Consider creating a new %s, if you want to revive it.",
							motion.GithubType(),
							motion.GithubType(),
						),
					)
				} else {
					motionapi.AppendMotionNotices_StageOnly(
						ctx,
						cloned.PublicClone(),
						id,
						notice.Noticef(
							ctx,
							"Gov4Git closed this issue, as the corresponding governance motion `%v` has now been closed.",
							id,
						),
					)
				}

			case !issue.Closed && !motion.Closed:

			}
			if changed {
				syncChanges.IssuesCausingChange = append(syncChanges.IssuesCausingChange, issue)
			}

		} else { // otherwise, no motion for this issue exists, so create one

			if !issue.Closed {
				syncCreateMotionForIssue(ctx, addr, cloned, syncChanges, issue, id)
				syncChanges.IssuesCausingChange = append(syncChanges.IssuesCausingChange, issue)
			}

		}

	} else { // issue is not governed, freeze motion if it exists and is open

		if motion, ok := motions[id]; ok { // motion for issue already exists, update it
			// if motion closed, do nothing
			// if motion frozen, do nothing
			// otherwise, freeze motion
			if !motion.Closed && !motion.Frozen {
				motionapi.AppendMotionNotices_StageOnly(
					ctx,
					cloned.PublicClone(),
					id,
					notice.Noticef(ctx, "The Gov4Git motion for this no longer managed issue/PR has been frozen."),
				)
				motionapi.FreezeMotion_StageOnly(notice.Mute(ctx), cloned, id)
				syncChanges.Froze.Add(id)
				syncChanges.IssuesCausingChange = append(syncChanges.IssuesCausingChange, issue)
			}
		}

	}
}

func syncMeta(
	ctx context.Context,
	cloned gov.OwnerCloned,
	chg *SyncManagedChanges,
	issue ImportedIssue,
	motion motionproto.Motion,

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
	motionapi.EditMotionMeta_StageOnly(
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

func motionPolicyForIssue(issue ImportedIssue) motion.PolicyName {
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
	id motionproto.MotionID,
) {

	must.Assertf(ctx, !issue.Closed, "issue is closed")

	motionapi.OpenMotion_StageOnly(
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

func indexMotions(ms motionproto.Motions) map[motionproto.MotionID]motionproto.Motion {
	x := map[motionproto.MotionID]motionproto.Motion{}
	for _, m := range ms {
		x[m.ID] = m
	}
	return x
}

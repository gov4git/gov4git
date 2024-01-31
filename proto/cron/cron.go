package cron

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/v55/github"
	"github.com/gov4git/gov4git/v2"
	govgh "github.com/gov4git/gov4git/v2/github"
	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotapi"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/motion/motionapi"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

var CronNS = ns.NS{"cron", "cron.json"}

func Cron(
	ctx context.Context,
	repo govgh.Repo,
	ghc *github.Client,
	govAddr gov.OwnerAddress,
	//
	githubFreq time.Duration, // frequency of importing from github
	communityFreq time.Duration, // frequency of fetching community votes and service requests
	//
	maxPar int, // parallelism for fetching community votes
) form.Map {

	cloned := gov.CloneOwner(ctx, govAddr)
	govTree := cloned.Public.Tree()

	// use a separate branch for cron logs
	cronAddr := git.Address(govAddr.Public)
	cronAddr.Branch = cronAddr.Branch + ".cron"
	cronCloned := git.CloneOne(ctx, cronAddr)
	cronTree := cronCloned.Tree()

	// read cron state
	state, err := git.TryFromFile[CronState](ctx, cronTree, CronNS)
	must.Assertf(ctx, err == nil || err == os.ErrNotExist, "opening cron state (%v)", err)

	now := time.Now()
	shouldSyncGithub := now.Sub(state.LastGithubImport) > githubFreq
	shouldSyncCommunity := now.Sub(state.LastCommunityTally) > communityFreq

	report := form.Map{}

	// import from github
	if shouldSyncGithub {

		// fetch repo maintainers
		maintainers := govgh.FetchRepoMaintainers(ctx, repo, ghc)
		base.Infof("maintainers for %v are %v", repo, form.SprintJSON(maintainers))

		// process managed issues and pull requests
		base.Infof("CRON: syncing managed issues and pull requests")
		report["processed_managed_issues"] = govgh.SyncManagedIssues_StageOnly(ctx, repo, ghc, govAddr, cloned)

		// process joins
		base.Infof("CRON: processing join requests")
		report["processed_joins"] = govgh.ProcessJoinRequestIssues_StageOnly(ctx, repo, ghc, govAddr, cloned, maintainers, false)

		// process directives
		base.Infof("CRON: processing directives")
		report["processed_directives"] = govgh.ProcessDirectiveIssues_StageOnly(ctx, repo, ghc, govAddr, cloned, maintainers)

		state.LastGithubImport = time.Now()
	}

	// sync community
	if shouldSyncCommunity {

		// tally votes for all ballots from all community members
		base.Infof("CRON: tallying community votes")
		report["tally"] = ballotapi.TallyAll_StageOnly(ctx, cloned, maxPar).Result

		state.LastCommunityTally = time.Now()
	}

	motionapi.MainPipeline(ctx, cloned)

	// display notices on github
	govgh.DisplayNotices_StageOnly(ctx, repo, ghc, cloned.PublicClone())

	if shouldSyncCommunity {

		// update community dashboard on github
		base.Infof("CRON: publishing community dashboard")
		govgh.PublishDashboard(ctx, repo, ghc, cloned.PublicClone())
	}

	// prepare commit message
	report["cron"] = state
	ver := gov4git.GetVersionInfo()
	latestChange := LatestChange{
		Stamp:          now,
		Gov4GitVersion: ver,
	}

	git.ToFileStage[LatestChange](ctx, cloned.PublicClone().Tree(), LatestChangeMetaNS, latestChange)

	cronChg := git.NewChange[form.Map, LatestChange](
		fmt.Sprintf("Gov4Git %s cron job.", ver.Version),
		"cron",
		nil,
		// We used to include the report in the commit message. However this causes a problem on GitHub.
		// The report includes the bodies of the issues that were processed.
		// It turns out GitHub scans the commit message for "resolves issue" text and automatically closes issues based on those.
		// This triggers spurious closures.
		latestChange,
		nil,
	)

	// push gov state
	govStatus, err := govTree.Status()
	must.NoError(ctx, err)
	if !govStatus.IsClean() {
		proto.Commit(ctx, cloned.Public.Tree(), cronChg)
		cloned.Public.Push(ctx)
	}

	// push cron state
	git.ToFileStage(ctx, cronTree, CronNS, state)
	proto.Commit(ctx, cronTree, cronChg)
	cronCloned.Push(ctx)

	return report
}

var LatestChangeMetaNS = ns.ParseFromGitPath("latest_change.json")

type LatestChange struct {
	Stamp          time.Time           `json:"change_stamp"`
	Gov4GitVersion gov4git.VersionInfo `json:"gov4git_version"`
}

type CronState struct {
	LastGithubImport   time.Time `json:"last_github_import"`
	LastCommunityTally time.Time `json:"last_community_tally"`
}

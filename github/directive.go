package github

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/go-github/v55/github"
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/account"
	"github.com/gov4git/gov4git/proto/docket/ops"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/util"
)

type ProcessDirectiveIssueReports []ProcessDirectiveIssueReport

type ProcessDirectiveIssueReport struct {
	Directive DirectiveIssue `json:"directive"`
	Success   *bool          `json:"success,omitempty"`
	Error     *string        `json:"error,omitempty"`
}

type DirectiveIssue struct {
	IssueVotingCredits    *IssueVotingCreditsDirective    `json:"issue_voting_credits,omitempty"`
	TransferVotingCredits *TransferVotingCreditsDirective `json:"transfer_voting_credits,omitempty"`
	Freeze                *FreezeDirective                `json:"freeze,omitempty"`
	Unfreeze              *UnfreezeDirective              `json:"unfreeze,omitempty"`
}

type IssueVotingCreditsDirective struct {
	Amount float64 `json:"amount"`
	To     string  `json:"to_user"`
}

type TransferVotingCreditsDirective struct {
	Amount float64 `json:"amount"`
	From   string  `json:"from_user"`
	To     string  `json:"to_user"`
}

type FreezeDirective struct {
	IssueURL    string `json:"issue_url"`
	IssueNumber int64  `json:"issue_number"`
}

type UnfreezeDirective struct {
	IssueURL    string `json:"issue_url"`
	IssueNumber int64  `json:"issue_number"`
}

func ProcessDirectiveIssuesByMaintainer(
	ctx context.Context,
	repo Repo,
	ghc *github.Client,
	govAddr gov.OwnerAddress,

) git.Change[form.Map, ProcessDirectiveIssueReports] {

	maintainers := FetchRepoMaintainers(ctx, repo, ghc)
	base.Infof("maintainers for %v are %v", repo, form.SprintJSON(maintainers))
	return ProcessDirectiveIssues(ctx, repo, ghc, govAddr, maintainers)
}

func ProcessDirectiveIssues(
	ctx context.Context,
	repo Repo,
	ghc *github.Client,
	govAddr gov.OwnerAddress,
	approverGitHubUsers []string,

) git.Change[form.Map, ProcessDirectiveIssueReports] {

	cloned := gov.CloneOwner(ctx, govAddr)
	report := ProcessDirectiveIssues_StageOnly(ctx, repo, ghc, govAddr, cloned, approverGitHubUsers)
	chg := git.NewChange[form.Map, ProcessDirectiveIssueReports](
		fmt.Sprintf("Process %d organizer directives", len(report)),
		"github_directive_issues",
		form.Map{},
		report,
		nil,
	)
	status, err := cloned.Public.Tree().Status()
	must.NoError(ctx, err)
	if !status.IsClean() {
		proto.Commit(ctx, cloned.Public.Tree(), chg)
		cloned.Public.Push(ctx)
	}
	return chg
}

func ProcessDirectiveIssues_StageOnly(
	ctx context.Context,
	repo Repo,
	ghc *github.Client, // if nil, a new client for repo will be created
	govAddr gov.OwnerAddress,
	govCloned gov.OwnerCloned,
	maintainers []string,

) ProcessDirectiveIssueReports { // return list of processed directives

	report := ProcessDirectiveIssueReports{}

	// fetch open issues labelled gov4git:directive
	issues := fetchOpenIssues(ctx, repo, ghc, DirectiveLabel)
	for _, issue := range issues {
		directive, err := processDirectiveIssue_StageOnly(ctx, repo, ghc, govAddr, govCloned, maintainers, issue)
		if err != nil {
			report = append(report, ProcessDirectiveIssueReport{
				Directive: directive,
				Error:     github.String(err.Error()),
			})
		} else {
			report = append(report, ProcessDirectiveIssueReport{
				Directive: directive,
				Success:   github.Bool(true),
			})
		}
	}
	return report
}

func processDirectiveIssue_StageOnly(
	ctx context.Context,
	repo Repo,
	ghc *github.Client,
	govAddr gov.OwnerAddress,
	cloned gov.OwnerCloned,
	maintainers []string,
	issue *github.Issue,

) (DirectiveIssue, error) {

	must.Assertf(ctx, len(maintainers) > 0, "no maintainers found")

	u := issue.GetUser()
	if u == nil {
		base.Infof("github identity of issue author is not available: %v", form.SprintJSON(issue))
		replyAndCloseIssue(ctx, repo, ghc, issue, FollowUpSubject, "The GitHub identity of the issue's author is not available.")
		return DirectiveIssue{}, fmt.Errorf("identity of issue author is not available")
	}
	login := strings.ToLower(u.GetLogin())
	if login == "" {
		base.Infof("github user of issue author is not available: %v", form.SprintJSON(issue))
		replyAndCloseIssue(ctx, repo, ghc, issue, FollowUpSubject, "The GitHub login of the issue's author is not available.")
		return DirectiveIssue{}, fmt.Errorf("user of issue author is not available")
	}

	d, err := parseDirective(issue.GetBody())
	if err != nil {
		base.Infof("directive cannot be parsed (%v): %q", err, issue.GetBody())
		replyAndCloseIssue(ctx, repo, ghc, issue, FollowUpSubject, "Your directive cannot be parsed.")
		return DirectiveIssue{}, err
	}

	switch {

	case d.IssueVotingCredits != nil:
		err = must.Try(
			func() {
				account.Issue_StageOnly(
					ctx,
					cloned.PublicClone(),
					member.UserAccountID(member.User(d.IssueVotingCredits.To)),
					account.H(account.PluralAsset, d.IssueVotingCredits.Amount),
					fmt.Sprintf("directive from GitHub issue #%v", issue.GetNumber()),
				)
			},
		)
		if err != nil {
			base.Infof("could not issue %v credits to member %v (%v)",
				d.IssueVotingCredits.Amount, d.IssueVotingCredits.To, err)
			replyAndCloseIssue(
				ctx, repo, ghc, issue,
				FollowUpSubject,
				fmt.Sprintf("Could not issue `%v` credits to member @%v. Reopen the issue to retry.\n\nBecause: `%v`",
					d.IssueVotingCredits.Amount, d.IssueVotingCredits.To, err))
			return DirectiveIssue{}, err
		}
		replyAndCloseIssue(ctx, repo, ghc, issue,
			FollowUpSubject,
			fmt.Sprintf("Issued `%v` credits to member @%v.",
				d.IssueVotingCredits.Amount, d.IssueVotingCredits.To))
		return d, nil

	case d.TransferVotingCredits != nil:
		err = must.Try(
			func() {
				account.Transfer_StageOnly(
					ctx,
					cloned.PublicClone(),
					member.UserAccountID(member.User(d.TransferVotingCredits.From)),
					member.UserAccountID(member.User(d.TransferVotingCredits.To)),
					account.H(account.PluralAsset, d.TransferVotingCredits.Amount),
					fmt.Sprintf("directive from GitHub issue #%v", issue.GetNumber()),
				)
			},
		)
		if err != nil {
			base.Infof("could not transfer %v credits from member %v to member %v (%v)",
				d.TransferVotingCredits.Amount, d.TransferVotingCredits.From, d.TransferVotingCredits.To, err)
			replyAndCloseIssue(
				ctx, repo, ghc, issue,
				FollowUpSubject,
				fmt.Sprintf("Could not transfer `%v` credits from member @%v to member @%v. Reopen the issue to retry.\n\nBecause: `%v`",
					d.TransferVotingCredits.Amount, d.TransferVotingCredits.From, d.TransferVotingCredits.To, err))
			return DirectiveIssue{}, err
		}
		replyAndCloseIssue(ctx, repo, ghc, issue,
			FollowUpSubject,
			fmt.Sprintf("Transferred `%v` credits from member @%v to member @%v.",
				d.TransferVotingCredits.Amount, d.TransferVotingCredits.From, d.TransferVotingCredits.To))
		return d, nil

	case d.Freeze != nil:
		id := IssueNumberToMotionID(d.Freeze.IssueNumber)
		err = must.Try(
			func() {
				ops.FreezeMotion_StageOnly(ctx, cloned, id)
			},
		)
		if err != nil {
			base.Infof("could not freeze issue/PR %v (%v)", d.Freeze.IssueURL, err)
			replyAndCloseIssue(
				ctx, repo, ghc, issue, FollowUpSubject,
				fmt.Sprintf("Could not freeze issue/PR %v. Reopen this issue to retry.\n\nBecause: `%v`", d.Freeze.IssueURL, err))
			return DirectiveIssue{}, err
		}
		replyAndCloseIssue(
			ctx, repo, ghc, issue, FollowUpSubject,
			fmt.Sprintf("Froze issue/PR %v.", d.Freeze.IssueURL),
		)
		return d, nil

	case d.Unfreeze != nil:
		id := IssueNumberToMotionID(d.Unfreeze.IssueNumber)
		err = must.Try(
			func() {
				ops.UnfreezeMotion_StageOnly(ctx, cloned, id)
			},
		)
		if err != nil {
			base.Infof("could not unfreeze issue/PR %v (%v)", d.Unfreeze.IssueURL, err)
			replyAndCloseIssue(
				ctx, repo, ghc, issue, FollowUpSubject,
				fmt.Sprintf("Could not unfreeze issue/PR %v. Reopen this issue to retry.\n\nBecause: `%v`", d.Unfreeze.IssueURL, err))
			return DirectiveIssue{}, err
		}
		replyAndCloseIssue(
			ctx, repo, ghc, issue, FollowUpSubject,
			fmt.Sprintf("Unfroze issue/PR %v.", d.Unfreeze.IssueURL),
		)
		return d, nil

	}
	panic("unknown directive")
}

// example directives:
//
//	"issue 30 credits to @user"
//	"transfer 20 credits from @user1 to @user2"
func parseDirective(body string) (DirectiveIssue, error) {
	body = strings.ToLower(body)
	body = strings.ReplaceAll(body, "\n", " ")
	body = strings.ReplaceAll(body, "\r", " ")
	body = strings.ReplaceAll(body, "\t", " ")
	body = strings.Trim(body, "\t .")
	words := []string{}
	for _, w := range strings.Split(body, " ") {
		if w != "" {
			words = append(words, w)
		}
	}

	if d, err := parseIssueCreditsDirective(words); err == nil {
		return d, nil
	}

	if d, err := parseTransferCreditsDirective(words); err == nil {
		return d, nil
	}

	if d, err := parseFreezeIssueDirective(words); err == nil {
		return d, nil
	}

	if d, err := parseUnfreezeIssueDirective(words); err == nil {
		return d, nil
	}

	return DirectiveIssue{}, fmt.Errorf("unrecognized directive")
}

// "issue 30 credits to @user"
func parseIssueCreditsDirective(words []string) (DirectiveIssue, error) {
	if len(words) == 5 &&
		words[0] == "issue" &&
		util.IsIn(words[2], "credit", "credits", "token", "tokens") &&
		words[3] == "to" {
		amount, err := strconv.ParseFloat(words[1], 64)
		if err != nil {
			return DirectiveIssue{}, fmt.Errorf("cannot parse amount of credits")
		}
		user, err := parseUser(words[4])
		if err != nil {
			return DirectiveIssue{}, err
		}
		return DirectiveIssue{
			IssueVotingCredits: &IssueVotingCreditsDirective{Amount: amount, To: user},
		}, nil
	}
	return DirectiveIssue{}, fmt.Errorf("cannot parse issue credits directive")
}

// "transfer 20 credits from @user1 to @user2"
func parseTransferCreditsDirective(words []string) (DirectiveIssue, error) {
	if len(words) == 7 &&
		words[0] == "transfer" &&
		util.IsIn(words[2], "credit", "credits", "token", "tokens") &&
		words[3] == "from" &&
		words[5] == "to" {
		amount, err := strconv.ParseFloat(words[1], 64)
		if err != nil {
			return DirectiveIssue{}, fmt.Errorf("cannot parse amount of credits")
		}
		from, err := parseUser(words[4])
		if err != nil {
			return DirectiveIssue{}, err
		}
		to, err := parseUser(words[6])
		if err != nil {
			return DirectiveIssue{}, err
		}
		return DirectiveIssue{
			TransferVotingCredits: &TransferVotingCreditsDirective{Amount: amount, From: from, To: to},
		}, nil
	}
	return DirectiveIssue{}, fmt.Errorf("cannot parse transfer credits directive")
}

// "freeze issueURL"
func parseFreezeIssueDirective(words []string) (DirectiveIssue, error) {
	if len(words) < 2 || words[0] != "freeze" {
		return DirectiveIssue{}, fmt.Errorf("cannot parse freeze issue directive")
	}

	matches := refRegexp.FindStringSubmatch(words[1])
	if matches == nil {
		return DirectiveIssue{}, fmt.Errorf("cannot parse freeze issue directive (expecting an issue URL)")
	}
	n, err := strconv.Atoi(matches[5])
	if err != nil {
		return DirectiveIssue{}, fmt.Errorf("cannot parse freeze issue directive (issue number not parsable)")
	}

	return DirectiveIssue{
		Freeze: &FreezeDirective{
			IssueURL:    matches[0],
			IssueNumber: int64(n),
		},
	}, nil
}

// "unfreeze issueURL"
func parseUnfreezeIssueDirective(words []string) (DirectiveIssue, error) {
	if len(words) < 2 || words[0] != "unfreeze" {
		return DirectiveIssue{}, fmt.Errorf("cannot parse unfreeze issue directive")
	}

	matches := refRegexp.FindStringSubmatch(words[1])
	if matches == nil {
		return DirectiveIssue{}, fmt.Errorf("cannot parse unfreeze issue directive (expecting an issue URL)")
	}
	n, err := strconv.Atoi(matches[5])
	if err != nil {
		return DirectiveIssue{}, fmt.Errorf("cannot parse unfreeze issue directive (issue number not parsable)")
	}

	return DirectiveIssue{
		Unfreeze: &UnfreezeDirective{
			IssueURL:    matches[0],
			IssueNumber: int64(n),
		},
	}, nil
}

func parseUser(s string) (string, error) {
	if len(s) < 2 {
		return "", fmt.Errorf("cannot parse user")
	}
	if s[0] != '@' {
		return "", fmt.Errorf("user must start with @")
	}
	return strings.ToLower(s[1:]), nil
}

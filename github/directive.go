package github

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/go-github/v55/github"
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/balance"
	"github.com/gov4git/gov4git/proto/ballot/qv"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

type ProcessDirectiveIssueReports []ProcessDirectiveIssueReport

type ProcessDirectiveIssueReport struct {
	Directive DirectiveIssue `json:"directive"`
	Success   *bool          `json:"success,omitempty"`
	Error     *string        `json:"error,omitempty"`
}

type DirectiveIssue struct {
	IssueVotingCredits    *IssueVotingCreditsDirective    `json:"issue_voting_credits"`
	TransferVotingCredits *TransferVotingCreditsDirective `json:"transfer_voting_credits"`
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

func ProcessDirectiveIssuesByMaintainer(
	ctx context.Context,
	repo GithubRepo,
	ghc *github.Client,
	govAddr gov.OrganizerAddress,
) git.Change[form.Map, ProcessDirectiveIssueReports] {

	maintainers := fetchRepoMaintainers(ctx, repo, ghc)
	base.Infof("maintainers for %v are %v", repo, form.SprintJSON(maintainers))
	return ProcessDirectiveIssues(ctx, repo, ghc, govAddr, maintainers)
}

func ProcessDirectiveIssues(
	ctx context.Context,
	repo GithubRepo,
	ghc *github.Client,
	govAddr gov.OrganizerAddress,
	approverGitHubUsers []string,
) git.Change[form.Map, ProcessDirectiveIssueReports] {

	govCloned := id.CloneOwner(ctx, id.OwnerAddress(govAddr))
	report := ProcessDirectiveIssues_Local(ctx, repo, ghc, govAddr, govCloned, approverGitHubUsers)
	chg := git.NewChange[form.Map, ProcessDirectiveIssueReports](
		fmt.Sprintf("Process %d organizer directives", len(report)),
		"github_directive_issues",
		form.Map{},
		report,
		nil,
	)
	status, err := govCloned.Public.Tree().Status()
	must.NoError(ctx, err)
	if !status.IsClean() {
		proto.Commit(ctx, govCloned.Public.Tree(), chg)
		govCloned.Public.Push(ctx)
	}
	return chg
}

func ProcessDirectiveIssues_Local(
	ctx context.Context,
	repo GithubRepo,
	ghc *github.Client, // if nil, a new client for repo will be created
	govAddr gov.OrganizerAddress,
	govCloned id.OwnerCloned,
	maintainers []string,
) ProcessDirectiveIssueReports { // return list of processed directives

	report := ProcessDirectiveIssueReports{}

	// fetch open issues labelled gov4git:directive
	issues := fetchOpenIssues(ctx, repo, ghc, DirectiveLabel)
	for _, issue := range issues {
		directive, err := processDirectiveIssue_Local(ctx, repo, ghc, govAddr, govCloned, maintainers, issue)
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

func processDirectiveIssue_Local(
	ctx context.Context,
	repo GithubRepo,
	ghc *github.Client,
	govAddr gov.OrganizerAddress,
	govCloned id.OwnerCloned,
	maintainers []string,
	issue *github.Issue,
) (DirectiveIssue, error) {

	must.Assertf(ctx, len(maintainers) > 0, "no maintainers found")

	u := issue.GetUser()
	if u == nil {
		base.Infof("github identity of issue author is not available: %v", form.SprintJSON(issue))
		replyAndCloseIssue(ctx, repo, ghc, issue, "GitHub identity of issue author is not available.")
		return DirectiveIssue{}, fmt.Errorf("identity of issue author is not available")
	}
	login := u.GetLogin()
	if login == "" {
		base.Infof("github user of issue author is not available: %v", form.SprintJSON(issue))
		replyAndCloseIssue(ctx, repo, ghc, issue, "GitHub user of issue author is not available.")
		return DirectiveIssue{}, fmt.Errorf("user of issue author is not available")
	}

	d, err := parseDirective(issue.GetBody())
	if err != nil {
		base.Infof("directive cannot be parsed (%v): %q", err, issue.GetBody())
		replyAndCloseIssue(ctx, repo, ghc, issue, "Directive cannot be parsed.")
		return DirectiveIssue{}, err
	}

	switch {
	case d.IssueVotingCredits != nil:
		err = must.Try(
			func() {
				balance.AddStageOnly(
					ctx,
					govCloned.Public.Tree(),
					member.User(d.IssueVotingCredits.To),
					qv.VotingCredits,
					d.IssueVotingCredits.Amount,
				)
			},
		)
		if err != nil {
			base.Infof("could not issue %v voting credits to member %v (%v)",
				d.IssueVotingCredits.Amount, d.IssueVotingCredits.To, err)
			replyAndCloseIssue(
				ctx, repo, ghc, issue,
				fmt.Sprintf("Could not issue %v voting credits to member %v (%v). Reopen the issue to retry.",
					d.IssueVotingCredits.Amount, d.IssueVotingCredits.To, err))
			return DirectiveIssue{}, err
		}
		replyAndCloseIssue(ctx, repo, ghc, issue,
			fmt.Sprintf("Issued %v voting credits to member %v.",
				d.IssueVotingCredits.Amount, d.IssueVotingCredits.To))
		return d, nil

	case d.TransferVotingCredits != nil:
		err = must.Try(
			func() {
				balance.TransferStageOnly(
					ctx,
					govCloned.Public.Tree(),
					member.User(d.TransferVotingCredits.From),
					qv.VotingCredits,
					member.User(d.TransferVotingCredits.To),
					qv.VotingCredits,
					d.TransferVotingCredits.Amount,
				)
			},
		)
		if err != nil {
			base.Infof("could not transfer %v voting credits from member %v to member %v (%v)",
				d.TransferVotingCredits.Amount, d.TransferVotingCredits.From, d.TransferVotingCredits.To, err)
			replyAndCloseIssue(
				ctx, repo, ghc, issue,
				fmt.Sprintf("Could not transfer %v voting credits from member %v to member %v (%v). Reopen the issue to retry.",
					d.TransferVotingCredits.Amount, d.TransferVotingCredits.From, d.TransferVotingCredits.To, err))
			return DirectiveIssue{}, err
		}
		replyAndCloseIssue(ctx, repo, ghc, issue,
			fmt.Sprintf("Transferred %v voting credits from member %v to member %v.",
				d.TransferVotingCredits.Amount, d.TransferVotingCredits.From, d.TransferVotingCredits.To))
		return d, nil

	}
	panic("unknown directive")
}

// example directives:
//
//	"issue 30 voting credits to @user"
//	"transfer 20 voting credits from @user1 to @user2"
func parseDirective(body string) (DirectiveIssue, error) {
	body = strings.ToLower(body)
	body = strings.ReplaceAll(body, "\n", " ")
	body = strings.ReplaceAll(body, "\r", " ")
	body = strings.ReplaceAll(body, "\t", " ")
	body = strings.Trim(body, "\t ")
	words := []string{}
	for _, w := range strings.Split(body, " ") {
		if w != "" {
			words = append(words, w)
		}
	}

	// "issue 30 voting credits to @user"
	if len(words) == 6 &&
		words[0] == "issue" &&
		words[2] == "voting" &&
		words[3] == "credits" &&
		words[4] == "to" {
		amount, err := strconv.ParseFloat(words[1], 64)
		if err != nil {
			return DirectiveIssue{}, fmt.Errorf("cannot parse amount of voting credits")
		}
		user, err := parseUser(words[5])
		if err != nil {
			return DirectiveIssue{}, err
		}
		return DirectiveIssue{
			IssueVotingCredits: &IssueVotingCreditsDirective{Amount: amount, To: user},
		}, nil
	}

	// "transfer 20 voting credits from @user1 to @user2"
	if len(words) == 8 &&
		words[0] == "transfer" &&
		words[2] == "voting" &&
		words[3] == "credits" &&
		words[4] == "from" &&
		words[6] == "to" {
		amount, err := strconv.ParseFloat(words[1], 64)
		if err != nil {
			return DirectiveIssue{}, fmt.Errorf("cannot parse amount of voting credits")
		}
		from, err := parseUser(words[5])
		if err != nil {
			return DirectiveIssue{}, err
		}
		to, err := parseUser(words[7])
		if err != nil {
			return DirectiveIssue{}, err
		}
		return DirectiveIssue{
			TransferVotingCredits: &TransferVotingCreditsDirective{Amount: amount, From: from, To: to},
		}, nil
	}

	return DirectiveIssue{}, fmt.Errorf("unrecognized directive")
}

func parseUser(s string) (string, error) {
	if len(s) < 2 {
		return "", fmt.Errorf("cannot parse user")
	}
	if s[0] != '@' {
		return "", fmt.Errorf("user must start with @")
	}
	return s[1:], nil
}

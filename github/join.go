package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v55/github"
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func ProcessJoinRequestIssues(
	ctx context.Context,
	repo GithubRepo,
	githubClient *github.Client, // if nil, a new client for repo will be created
	govAddr gov.OrganizerAddress,
	approverGitHubUsers []string,
) git.Change[form.Map, []string] {

	govCloned := id.CloneOwner(ctx, id.OwnerAddress(govAddr))
	newMembers := ProcessJoinRequestIssues_Local(ctx, repo, githubClient, govAddr, govCloned, approverGitHubUsers)
	chg := git.NewChange[form.Map, []string](
		fmt.Sprintf("Add %d new community members", len(newMembers)),
		"github_process_join_request_issues",
		form.Map{},
		newMembers,
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

func ProcessJoinRequestIssues_Local(
	ctx context.Context,
	repo GithubRepo,
	githubClient *github.Client, // if nil, a new client for repo will be created
	govAddr gov.OrganizerAddress,
	govCloned id.OwnerCloned,
	approvers []string,
) []string { // return list of new member usernames

	// fetch open issues labelled gov4git:join
	issues := fetchOpenJoinRequestIssues(ctx, repo, githubClient)
	newMembers := []string{}
	for _, issue := range issues {
		newMember := processJoinRequestIssue_Local(ctx, repo, githubClient, govAddr, govCloned, approvers, issue)
		if newMember != "" {
			newMembers = append(newMembers, newMember)
		}
	}
	return newMembers
}

func fetchOpenJoinRequestIssues(ctx context.Context, repo GithubRepo, ghc *github.Client) []*github.Issue {
	opt := &github.IssueListByRepoOptions{State: "open", Labels: []string{JoinRequestLabel}}
	var issues []*github.Issue
	for {
		issues, resp, err := ghc.Issues.ListByRepo(ctx, repo.Owner, repo.Name, opt)
		must.NoError(ctx, err)
		issues = append(issues, issues...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return issues
}

//	parse form info
//	fetch comments
//	find approval comment by an approver GitHub user
//	add requesting user to community
//		if success, reply to GitHub issue, close issue
//		if user already exists, reply to GitHub issue, close issue
//		if other error, log it

func processJoinRequestIssue_Local(
	ctx context.Context,
	repo GithubRepo,
	githubClient *github.Client, // if nil, a new client for repo will be created
	govAddr gov.OrganizerAddress,
	govCloned id.OwnerCloned,
	approverGitHubUsers []string,
	issue *github.Issue,
) string { // return new member username, if joined

	u := issue.GetUser()
	if u == nil {
		// XXX // reply to issue
		// XXX // close issue
		return ""
	}
	login := u.GetLogin()
	if login == "" {
		// XXX // reply to issue
		// XXX // close issue
		return ""
	}

	info, err := parseJoinRequest(login, issue.GetBody())
	if err != nil {
		// XXX // reply to issue
		// XXX // close issue
		return ""
	}
	if info.Email == "" {
		info.Email = u.GetEmail()
	}

	panic("XXX")
}

type JoinRequest struct {
	User         string     `json:"github_user"`
	PublicURL    git.URL    `json:"public_url"`
	PublicBranch git.Branch `json:"public_branch"`
	Email        string     `json:"email"`
}

// example request body:
// "### Your public repo\n\nhttps://github.com/petar/gov4git.public.git\n\n### Your public branch\n\nmain\n\n### Your email (optional)\n\npetar@protocol.ai"
/*
### Your public repo

https://github.com/petar/gov4git.public.git

### Your public branch

main

### Your email (optional)

petar@protocol.ai
*/

var ErrJoinSyntax = fmt.Errorf("join request format is unrecognizable")

func parseJoinRequest(authorLogin string, body string) (*JoinRequest, error) {
	lines := strings.Split(body, "\n")
	if len(lines) < 7 {
		return nil, ErrJoinSyntax
	}
	if strings.Index(lines[0], "public repo") < 0 {
		return nil, ErrJoinSyntax
	}
	if strings.Index(lines[4], "public branch") < 0 {
		return nil, ErrJoinSyntax
	}
	if lines[1] != "" || lines[3] != "" || lines[5] != "" {
		return nil, ErrJoinSyntax
	}
	if lines[2] == "" || lines[6] == "" {
		return nil, ErrJoinSyntax
	}
	return &JoinRequest{
		User:         authorLogin,
		PublicURL:    git.URL(lines[2]),
		PublicBranch: git.Branch(lines[6]),
	}, nil
}

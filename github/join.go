package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v55/github"
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/util"
)

func ProcessJoinRequestIssuesApprovedByMaintainer(
	ctx context.Context,
	repo GithubRepo,
	ghc *github.Client, // if nil, a new client for repo will be created
	govAddr gov.OrganizerAddress,
) git.Change[form.Map, ProcessJoinRequestIssuesReport] {

	maintainers := fetchRepoMaintainers(ctx, repo, ghc)
	base.Infof("maintainers for %v are %v", repo, form.SprintJSON(maintainers))
	return ProcessJoinRequestIssues(ctx, repo, ghc, govAddr, maintainers)
}

func fetchRepoMaintainers(
	ctx context.Context,
	repo GithubRepo,
	ghc *github.Client, // if nil, a new client for repo will be created
) []string {

	opts := &github.ListCollaboratorsOptions{}
	users, _, err := ghc.Repositories.ListCollaborators(ctx, repo.Owner, repo.Name, opts)
	must.NoError(ctx, err)

	m := []string{}
	for _, u := range users {
		if u.GetPermissions()["maintainer"] || u.GetPermissions()["admin"] {
			m = append(m, u.GetLogin())
		}
	}
	return m
}

func ProcessJoinRequestIssues(
	ctx context.Context,
	repo GithubRepo,
	ghc *github.Client, // if nil, a new client for repo will be created
	govAddr gov.OrganizerAddress,
	approverGitHubUsers []string,
) git.Change[form.Map, ProcessJoinRequestIssuesReport] {

	govCloned := id.CloneOwner(ctx, id.OwnerAddress(govAddr))
	report := ProcessJoinRequestIssues_Local(ctx, repo, ghc, govAddr, govCloned, approverGitHubUsers)
	chg := git.NewChange[form.Map, ProcessJoinRequestIssuesReport](
		fmt.Sprintf("Add %d new community members; skipped %d", len(report.Joined), len(report.NotJoined)),
		"github_process_join_request_issues",
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

type ProcessJoinRequestIssuesReport struct {
	Joined    []string `json:"joined"`
	NotJoined []string `json:"not_joined"`
}

func ProcessJoinRequestIssues_Local(
	ctx context.Context,
	repo GithubRepo,
	ghc *github.Client, // if nil, a new client for repo will be created
	govAddr gov.OrganizerAddress,
	govCloned id.OwnerCloned,
	approvers []string,
) ProcessJoinRequestIssuesReport { // return list of new member usernames

	report := ProcessJoinRequestIssuesReport{}

	// fetch open issues labelled gov4git:join
	issues := fetchOpenJoinRequestIssues(ctx, repo, ghc)
	for _, issue := range issues {
		newMember := processJoinRequestIssue_Local(ctx, repo, ghc, govAddr, govCloned, approvers, issue)
		if newMember != "" {
			report.Joined = append(report.Joined, newMember)
		} else {
			if issue.User != nil {
				report.NotJoined = append(report.NotJoined, issue.User.GetLogin())
			}
		}
	}
	return report
}

func fetchOpenJoinRequestIssues(ctx context.Context, repo GithubRepo, ghc *github.Client) []*github.Issue {
	opt := &github.IssueListByRepoOptions{State: "open", Labels: []string{JoinRequestLabel}}
	var allIssues []*github.Issue
	for {
		issues, resp, err := ghc.Issues.ListByRepo(ctx, repo.Owner, repo.Name, opt)
		must.NoError(ctx, err)
		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allIssues
}

func processJoinRequestIssue_Local(
	ctx context.Context,
	repo GithubRepo,
	ghc *github.Client, // if nil, a new client for repo will be created
	govAddr gov.OrganizerAddress,
	govCloned id.OwnerCloned,
	approverGitHubUsers []string,
	issue *github.Issue,
) string { // return new member username, if joined

	must.Assertf(ctx, len(approverGitHubUsers) > 0, "no membership approvers")

	u := issue.GetUser()
	if u == nil {
		base.Infof("github identity of issue author is not available: %v", form.SprintJSON(issue))
		replyAndCloseIssue(ctx, repo, ghc, issue, "GitHub identity of issue author is not available.")
		return ""
	}
	login := u.GetLogin()
	if login == "" {
		base.Infof("github user of issue author is not available: %v", form.SprintJSON(issue))
		replyAndCloseIssue(ctx, repo, ghc, issue, "GitHub user of issue author is not available.")
		return ""
	}

	info, err := parseJoinRequest(login, issue.GetBody())
	if err != nil {
		base.Infof("request form cannot be parsed: %q", issue.GetBody())
		replyAndCloseIssue(ctx, repo, ghc, issue, "Request form cannot be parsed.")
		return ""
	}
	if info.Email == "" {
		info.Email = u.GetEmail()
	}

	// fetch comments and find a join approval
	comments := fetchIssueComments(ctx, repo, ghc, issue)
	if !isJoinApprovalPresent(ctx, approverGitHubUsers, comments) {
		return ""
	}

	// add user to community members
	err = must.Try(
		func() {
			member.AddUserByPublicAddressStageOnly(ctx, govCloned.Public.Tree(), member.User(login), info.PublicAddress())
		},
	)
	if err != nil {
		base.Infof("could not add member %v (%v)", login, err)
		replyAndCloseIssue(ctx, repo, ghc, issue, fmt.Sprintf("Could not add member due to (%v). Reopen the issue to retry.", err))
		return ""
	}

	replyAndCloseIssue(ctx, repo, ghc, issue, fmt.Sprintf("%v added to community.", login))
	return login
}

type JoinRequest struct {
	User         string     `json:"github_user"`
	PublicURL    git.URL    `json:"public_url"`
	PublicBranch git.Branch `json:"public_branch"`
	Email        string     `json:"email"`
}

func (x JoinRequest) PublicAddress() id.PublicAddress {
	return id.PublicAddress{
		Repo:   git.URL(x.PublicURL),
		Branch: git.Branch(x.PublicBranch),
	}
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

func replyAndCloseIssue(
	ctx context.Context,
	repo GithubRepo,
	ghc *github.Client,
	issue *github.Issue,
	payload string,
) {
	replyToIssue(ctx, repo, ghc, issue, payload)
	closeIssue(ctx, repo, ghc, issue)
}

func replyToIssue(
	ctx context.Context,
	repo GithubRepo,
	ghc *github.Client,
	issue *github.Issue,
	payload string,
) {

	comment := &github.IssueComment{
		Body: github.String(payload),
	}
	_, _, err := ghc.Issues.CreateComment(ctx, repo.Owner, repo.Name, issue.GetNumber(), comment)
	must.NoError(ctx, err)
}

func closeIssue(
	ctx context.Context,
	repo GithubRepo,
	ghc *github.Client,
	issue *github.Issue,
) {
	req := &github.IssueRequest{
		State: github.String("closed"),
	}
	_, _, err := ghc.Issues.Edit(ctx, repo.Owner, repo.Name, issue.GetNumber(), req)
	must.NoError(ctx, err)
}

func fetchIssueComments(
	ctx context.Context,
	repo GithubRepo,
	ghc *github.Client,
	issue *github.Issue,
) []*github.IssueComment {

	if issue.GetComments() == 0 {
		return nil
	}
	opts := &github.IssueListCommentsOptions{}
	comments, _, err := ghc.Issues.ListComments(ctx, repo.Owner, repo.Name, issue.GetNumber(), opts)
	must.NoError(ctx, err)
	return comments
}

func isJoinApprovalPresent(ctx context.Context, approvers []string, comments []*github.IssueComment) bool {
	for _, comment := range comments {
		u := comment.GetUser()
		if u == nil {
			continue
		}
		if !util.IsIn(u.GetLogin(), approvers...) {
			continue
		}
		// trim empty lines and spaces
		trimmed := strings.ToLower(strings.Trim(comment.GetBody(), ". \t\n\r"))
		if strings.Index(trimmed, JoinRequestApprovalWord) < 0 {
			continue
		}
		return true
	}
	return false
}

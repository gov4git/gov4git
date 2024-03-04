package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v58/github"
	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/id"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func ProcessJoinRequestIssuesApprovedByMaintainer(
	ctx context.Context,
	repo Repo,
	ghc *github.Client, // if nil, a new client for repo will be created
	govAddr gov.OwnerAddress,
	allowNonGithubJoins bool,
) git.Change[form.Map, ProcessJoinRequestIssuesReport] {

	maintainers := FetchRepoMaintainers(ctx, repo, ghc)
	base.Infof("maintainers for %v are %v", repo, form.SprintJSON(maintainers))
	return ProcessJoinRequestIssues(ctx, repo, ghc, govAddr, maintainers, allowNonGithubJoins)
}

func ProcessJoinRequestIssues(
	ctx context.Context,
	repo Repo,
	ghc *github.Client, // if nil, a new client for repo will be created
	govAddr gov.OwnerAddress,
	approverGitHubUsers []string,
	allowNonGithubJoins bool,
) git.Change[form.Map, ProcessJoinRequestIssuesReport] {

	govCloned := gov.CloneOwner(ctx, govAddr)
	report := ProcessJoinRequestIssues_StageOnly(
		ctx,
		repo,
		ghc,
		govAddr,
		govCloned,
		approverGitHubUsers,
		allowNonGithubJoins,
	)
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

func ProcessJoinRequestIssues_StageOnly(
	ctx context.Context,
	repo Repo,
	ghc *github.Client, // if nil, a new client for repo will be created
	govAddr gov.OwnerAddress,
	govCloned gov.OwnerCloned,
	approvers []string,
	allowNonGithubJoins bool,
) ProcessJoinRequestIssuesReport { // return list of new member usernames

	report := ProcessJoinRequestIssuesReport{}

	// fetch open issues
	issues := fetchOpenIssues(ctx, repo, ghc)
	for _, issue := range issues {
		if !isJoinRequestIssue(issue) {
			continue
		}
		newMember := processJoinRequestIssue_StageOnly(ctx, repo, ghc, govAddr, govCloned, approvers, allowNonGithubJoins, issue)
		if newMember != "" {
			report.Joined = append(report.Joined, newMember)
		} else {
			if issue.User != nil {
				report.NotJoined = append(report.NotJoined, strings.ToLower(issue.User.GetLogin()))
			}
		}
	}
	return report
}

func isJoinRequestIssue(issue *github.Issue) bool {
	_, _, err := parseJoinBody(issue.GetBody())
	return err == nil
}

func processJoinRequestIssue_StageOnly(
	ctx context.Context,
	repo Repo,
	ghc *github.Client, // if nil, a new client for repo will be created
	_ gov.OwnerAddress,
	govCloned gov.OwnerCloned,
	approverGitHubUsers []string,
	allowNonGithubJoins bool,
	issue *github.Issue,
) string { // return new member username, if joined

	must.Assertf(ctx, len(approverGitHubUsers) > 0, "no membership approvers")

	if !isJoinRequestIssue(issue) {
		return ""
	}

	// find the github login of the requesting user
	u := issue.GetUser()
	if u == nil {
		base.Infof("github identity of issue author is not available: %v", form.SprintJSON(issue))
		replyAndCloseIssue(ctx, repo, ghc, issue, FollowUpSubject, "The GitHub identity of the issue's author is not available.")
		return ""
	}
	login := strings.ToLower(u.GetLogin())
	if login == "" {
		base.Infof("github user of issue author is not available: %v", form.SprintJSON(issue))
		replyAndCloseIssue(ctx, repo, ghc, issue, FollowUpSubject, "The GitHub user of the issue's author is not available.")
		return ""
	}

	// extract the join request from the github issue body
	info, err := parseJoinRequest(login, issue.GetBody())
	if err != nil {
		base.Infof("request form cannot be parsed: %q", issue.GetBody())
		replyAndCloseIssue(ctx, repo, ghc, issue, FollowUpSubject, "The join request form cannot be parsed.")
		return ""
	}
	if info.Email == "" {
		info.Email = u.GetEmail()
	}

	// verify that the gov4git repo url matches the login of the requesting user
	if !allowNonGithubJoins && info.PublicRepo.Owner != info.User {
		base.Infof("reguster's GitHub login %s does not match the public repo owner %s", info.User, info.PublicRepo.Owner)
		replyAndCloseIssue(
			ctx, repo, ghc, issue, FollowUpSubject,
			fmt.Sprintf(
				"The regusting user, @%s, does not match the owner, @%s, of the provided Gov4Git public identity repo.",
				info.User, info.PublicRepo.Owner,
			),
		)
		return ""
	}

	// fetch comments and find a join approval
	comments := fetchIssueComments(ctx, repo, ghc, issue)
	if !isJoinApprovalPresent(ctx, approverGitHubUsers, comments) {
		return ""
	}

	// add user to community members
	err = must.Try(
		func() {
			member.AddUserByPublicAddress_StageOnly(ctx, govCloned.PublicClone(), member.User(login), info.PublicAddress())
		},
	)
	if err != nil {
		base.Infof("could not add member %v (%v)", login, err)
		replyAndCloseIssue(ctx, repo, ghc, issue, FollowUpSubject, fmt.Sprintf("Could not add member due to `%v`. Reopen the issue to retry.", err))
		return ""
	}

	replyAndCloseIssue(ctx, repo, ghc, issue, FollowUpSubject, fmt.Sprintf("@%v was added to the community.", login))
	return login
}

type JoinRequest struct {
	User         string     `json:"github_user"`
	PublicRepo   Repo       `json:"public_repo"`
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
	publicURL, publicBranch, err := parseJoinBody(body)
	if err != nil {
		return nil, err
	}
	publicURL = git.URL(strings.TrimSpace(string(publicURL)))
	repo, _ := parseGithubRepoHTTPSURL(string(publicURL))
	// if repo is not a GitHub URL, that's ok
	return &JoinRequest{
		User:         authorLogin,
		PublicRepo:   repo,
		PublicURL:    publicURL,
		PublicBranch: publicBranch,
	}, nil
}

func parseJoinBody(body string) (publicURL git.URL, publicBranch git.Branch, err error) {
	lines := strings.Split(body, "\n")
	if len(lines) < 7 {
		return "", "", ErrJoinSyntax
	}
	if strings.Index(lines[0], "public repo") < 0 {
		return "", "", ErrJoinSyntax
	}
	if strings.Index(lines[4], "public branch") < 0 {
		return "", "", ErrJoinSyntax
	}
	if lines[1] != "" || lines[3] != "" || lines[5] != "" {
		return "", "", ErrJoinSyntax
	}
	if lines[2] == "" || lines[6] == "" {
		return "", "", ErrJoinSyntax
	}
	return git.URL(lines[2]), git.Branch(lines[6]), nil
}

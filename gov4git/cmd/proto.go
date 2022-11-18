package cmd

import (
	"context"

	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/git"
)

const (
	LocalAgentPath     = ".gov4git"
	LocalAgentTempPath = "gov4git"
)

type Setup struct {
	Community gov.CommunityAddress
	Organizer gov.OrganizerAddress
	Member    id.OwnerAddress
}

type Config struct {
	// auth
	SSHPrivateKeysFile *string       `json:"ssh_private_keys_file"`
	AccessToken        *string       `json:"access_token"`
	UserPassword       *UserPassword `json:"user_password"`

	//
	CommunityPublicURL    git.URL    `json:"community_public_url"`
	CommunityPublicBranch git.Branch `json:"community_public_branch"`
	//
	CommunityPrivateURL    git.URL    `json:"community_private_url"`
	CommunityPrivateBranch git.Branch `json:"community_private_branch"`

	//
	MemberPublicURL    git.URL    `json:"member_public_url"`
	MemberPublicBranch git.Branch `json:"member_public_branch"`
	//
	MemberPrivateURL    git.URL    `json:"member_private_url"`
	MemberPrivateBranch git.Branch `json:"member_private_branch"`
}

type UserPassword struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func (cfg Config) Setup(ctx context.Context) Setup {
	switch {
	case cfg.SSHPrivateKeysFile != nil:
		git.SetSSHFileAuth(ctx, "git", *cfg.SSHPrivateKeysFile)
	case cfg.AccessToken != nil:
		git.SetTokenAuth(ctx, *cfg.AccessToken)
	case cfg.UserPassword != nil:
		git.SetPasswordAuth(ctx, cfg.UserPassword.User, cfg.UserPassword.Password)
	}
	return Setup{
		Community: gov.CommunityAddress{Repo: cfg.CommunityPublicURL, Branch: cfg.CommunityPublicBranch},
		Organizer: gov.OrganizerAddress{
			Public:  id.PublicAddress{Repo: cfg.CommunityPublicURL, Branch: cfg.CommunityPublicBranch},
			Private: id.PrivateAddress{Repo: cfg.CommunityPrivateURL, Branch: cfg.CommunityPrivateBranch},
		},
		Member: id.OwnerAddress{
			Public:  id.PublicAddress{Repo: cfg.MemberPublicURL, Branch: cfg.MemberPublicBranch},
			Private: id.PrivateAddress{Repo: cfg.MemberPrivateURL, Branch: cfg.MemberPrivateBranch},
		},
	}
}

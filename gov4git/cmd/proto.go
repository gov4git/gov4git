package cmd

import (
	"context"
	"time"

	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/git"
)

const (
	LocalAgentPath     = ".gov4git"
	LocalAgentTempPath = "gov4git"
)

type Setup struct {
	CacheDir  string
	CacheTTL  time.Duration
	Gov       gov.GovAddress
	Organizer gov.OrganizerAddress
	Member    id.OwnerAddress
}

type Config struct {
	Auth map[git.URL]AuthConfig `json:"auth"`
	//
	GovPublicURL     git.URL    `json:"gov_public_url"`
	GovPublicBranch  git.Branch `json:"gov_public_branch"`
	GovPrivateURL    git.URL    `json:"gov_private_url"`
	GovPrivateBranch git.Branch `json:"gov_private_branch"`
	//
	MemberPublicURL     git.URL    `json:"member_public_url"`
	MemberPublicBranch  git.Branch `json:"member_public_branch"`
	MemberPrivateURL    git.URL    `json:"member_private_url"`
	MemberPrivateBranch git.Branch `json:"member_private_branch"`
	//
	CacheDir        string `json:"cache_dir"`
	CacheTTLSeconds int    `json:"cache_ttl_seconds"` // ttl of repo cache replicas in seconds
}

type AuthConfig struct {
	SSHPrivateKeysFile *string       `json:"ssh_private_keys_file"`
	AccessToken        *string       `json:"access_token"`
	UserPassword       *UserPassword `json:"user_password"`
}

type UserPassword struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func (cfg Config) Setup(ctx context.Context) Setup {

	git.SetAuthor("gov4git governance", "no-reply@gov4git")

	// attach auth information to context
	for url, auth := range cfg.Auth {
		switch {
		case auth.SSHPrivateKeysFile != nil:
			git.SetAuth(ctx, url, git.MakeSSHFileAuth(ctx, "git", *auth.SSHPrivateKeysFile))
		case auth.AccessToken != nil:
			git.SetAuth(ctx, url, git.MakeTokenAuth(ctx, *auth.AccessToken))
		case auth.UserPassword != nil:
			git.SetAuth(ctx, url, git.MakePasswordAuth(ctx, auth.UserPassword.User, auth.UserPassword.Password))
		}
	}

	s := Setup{
		CacheDir: cfg.CacheDir,
		CacheTTL: time.Second * time.Duration(cfg.CacheTTLSeconds),
		Gov:      gov.GovAddress{Repo: cfg.GovPublicURL, Branch: cfg.GovPublicBranch},
		Organizer: gov.OrganizerAddress{
			Public:  id.PublicAddress{Repo: cfg.GovPublicURL, Branch: cfg.GovPublicBranch},
			Private: id.PrivateAddress{Repo: cfg.GovPrivateURL, Branch: cfg.GovPrivateBranch},
		},
		Member: id.OwnerAddress{
			Public:  id.PublicAddress{Repo: cfg.MemberPublicURL, Branch: cfg.MemberPublicBranch},
			Private: id.PrivateAddress{Repo: cfg.MemberPrivateURL, Branch: cfg.MemberPrivateBranch},
		},
	}

	// attach cache ttl information to context
	git.SetTTL(ctx, setup.Gov.Repo, s.CacheTTL)
	git.SetTTL(ctx, setup.Organizer.Public.Repo, s.CacheTTL)
	git.SetTTL(ctx, setup.Organizer.Private.Repo, s.CacheTTL)
	git.SetTTL(ctx, setup.Member.Public.Repo, s.CacheTTL)
	git.SetTTL(ctx, setup.Member.Private.Repo, s.CacheTTL)

	return s
}

package testutil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/services/identity"
)

func CreateTestCommunity(dir string, numUsers int) (*TestCommunity, error) {
	c := &TestCommunity{Dir: dir, NumUsers: numUsers}
	if err := c.Init(context.Background()); err != nil {
		return nil, err
	}
	return c, nil
}

type TestCommunity struct {
	// Dir is an absolute local path to the git repository
	Dir      string
	NumUsers int
}

func (x *TestCommunity) Init(ctx context.Context) error {
	ctx = x.WithWorkDir(ctx, "TestCommunity")

	// prep test directory
	fmt.Printf("initializing a test community in %v\n", x.Dir)
	if err := os.RemoveAll(x.Dir); err != nil {
		return err
	}

	// init community repo
	if err := x.initCommunityRepo(ctx); err != nil {
		return err
	}

	// init user repos
	for i := 0; i < x.NumUsers; i++ {
		if err := x.initUserRepos(ctx, i); err != nil {
			return err
		}
	}

	return nil
}

func (x *TestCommunity) initCommunityRepo(ctx context.Context) error {
	ctx = x.WithWorkDir(ctx, "TestCommunity/initCommunityRepo")
	fmt.Printf("initializing community repo\n")

	// init community repo
	communityRepo := x.CommunityRepoLocal()
	if err := communityRepo.InitBare(ctx); err != nil {
		return err
	}

	// clone and make first commit
	clonedCommunityRepo := git.Local{Path: filepath.Join(x.Dir, "community_clone")}

	if err := clonedCommunityRepo.CloneOrInitBranch(ctx, x.CommunityRepoURL(), proto.MainBranch); err != nil {
		return err
	}
	if err := clonedCommunityRepo.Dir().WriteByteFile("empty", nil); err != nil {
		return err
	}
	if err := clonedCommunityRepo.Add(ctx, []string{"empty"}); err != nil {
		return err
	}
	if err := clonedCommunityRepo.Commit(ctx, "first"); err != nil {
		return err
	}
	if err := clonedCommunityRepo.PushUpstream(ctx); err != nil {
		return err
	}

	return nil
}

func (x *TestCommunity) initUserRepos(ctx context.Context, i int) error {
	ctx = x.WithWorkDir(ctx, fmt.Sprintf("TestCommunity/initUserRepos_%d", i))
	fmt.Printf("initializing user %d repos\n", i)

	// create user identity repos
	userPublicRepo := x.UserPublicRepoLocal(i)
	if err := userPublicRepo.InitBare(ctx); err != nil {
		return err
	}
	userPrivateRepo := x.UserPrivateRepoLocal(i)
	if err := userPrivateRepo.InitBare(ctx); err != nil {
		return err
	}

	// initialize user
	idService := identity.IdentityService{
		IdentityConfig: proto.IdentityConfig{
			PublicURL:  x.UserPublicRepoURL(i),
			PrivateURL: x.UserPrivateRepoURL(i),
		},
	}
	if _, err := idService.Init(ctx, &identity.InitIn{}); err != nil {
		return err
	}

	return nil
}

func (x *TestCommunity) WithWorkDir(ctx context.Context, suffix ...string) context.Context {
	return files.WithWorkDir(ctx, files.Dir{Path: filepath.Join(x.Dir, "working", filepath.Join(suffix...))})
}

// community repo

func (x *TestCommunity) CommunityRepoDir() string {
	return filepath.Join(x.Dir, "community_repo")
}

func (x *TestCommunity) CommunityRepoLocal() git.Local {
	return git.Local{Path: x.CommunityRepoDir()}
}

func (x *TestCommunity) CommunityRepoURL() string {
	return x.CommunityRepoDir()
}

func (x *TestCommunity) CommunityGovConfig() proto.GovConfig {
	return proto.GovConfig{
		CommunityURL: x.CommunityRepoURL(),
	}
}

// user public repos

func (x *TestCommunity) UserPublicRepoDir(u int) string {
	return filepath.Join(x.Dir, fmt.Sprintf("%d_public_repo", u))
}

func (x *TestCommunity) UserPublicRepoLocal(u int) git.Local {
	return git.Local{Path: x.UserPublicRepoDir(u)}
}

func (x *TestCommunity) UserPublicRepoURL(u int) string {
	return x.UserPublicRepoDir(u)
}

func (x *TestCommunity) UserIdentityConfig(i int) proto.IdentityConfig {
	return proto.IdentityConfig{
		PublicURL:  x.UserPublicRepoURL(i),
		PrivateURL: x.UserPrivateRepoURL(i),
	}
}

// user private repos

func (x *TestCommunity) UserPrivateRepoDir(u int) string {
	return filepath.Join(x.Dir, fmt.Sprintf("%d_private_repo", u))
}

func (x *TestCommunity) UserPrivateRepoLocal(u int) git.Local {
	return git.Local{Path: x.UserPrivateRepoDir(u)}
}

func (x *TestCommunity) UserPrivateRepoURL(u int) string {
	return x.UserPrivateRepoDir(u)
}

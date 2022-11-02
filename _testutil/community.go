package testutil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gov4git/gov4git/lib/base"
	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/govproto"
	"github.com/gov4git/gov4git/proto/idproto"
	"github.com/gov4git/gov4git/services/gov/group"
	"github.com/gov4git/gov4git/services/gov/member"
	"github.com/gov4git/gov4git/services/gov/user"
	"github.com/gov4git/gov4git/services/id"
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
	base.Infof("initializing a test community in %v\n", x.Dir)
	if err := os.RemoveAll(x.Dir); err != nil {
		return err
	}

	// init user repos
	for i := 0; i < x.NumUsers; i++ {
		if err := x.initUserRepos(ctx, i); err != nil {
			return err
		}
	}

	// init community repo
	if err := x.initCommunityRepo(ctx); err != nil {
		return err
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

	// add users to group all
	if err := x.addUsersToGroupAll(ctx, clonedCommunityRepo); err != nil {
		return err
	}

	return nil
}

func (x *TestCommunity) initUserRepos(ctx context.Context, i int) error {
	ctx = x.WithWorkDir(ctx, fmt.Sprintf("TestCommunity/initUserRepos_%d", i))
	fmt.Printf("initializing user %d with public repo %v and private repo %v\n", i, x.UserPublicRepoURL(i), x.UserPrivateRepoURL(i))

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
	idService := id.IdentityService{
		IdentityConfig: idproto.IdentityConfig{
			PublicURL:  x.UserPublicRepoURL(i),
			PrivateURL: x.UserPrivateRepoURL(i),
		},
	}
	if _, err := idService.Init(ctx, &id.InitIn{}); err != nil {
		return err
	}

	return nil
}

func (x *TestCommunity) addUsersToGroupAll(ctx context.Context, clonedCommunityRepo git.Local) error {

	// make group "all"
	if err := x.CommunityGroupService().AddLocalStageOnly(ctx, clonedCommunityRepo, "all"); err != nil {
		return err
	}

	// create users
	for i := 0; i < x.NumUsers; i++ {
		user := x.User(i)
		if err := x.CommunityUserService().AddLocalStageOnly(ctx, clonedCommunityRepo, user, x.UserPublicRepoURL(i)); err != nil {
			return err
		}
		if err := x.CommunityMemberService().AddLocalStageOnly(ctx, clonedCommunityRepo, user, "all"); err != nil {
			return err
		}
	}

	// commit changes
	if err := clonedCommunityRepo.Commit(ctx, "add users and groups"); err != nil {
		return err
	}
	if err := clonedCommunityRepo.PushUpstream(ctx); err != nil {
		return err
	}

	return nil
}

func (x *TestCommunity) Background() context.Context {
	return x.WithWorkDir(context.Background())
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

func (x *TestCommunity) CommunityGovConfig() govproto.GovConfig {
	return govproto.GovConfig{
		CommunityURL: x.CommunityRepoURL(),
	}
}

func (x *TestCommunity) CommunityUserService() user.GovUserService {
	return user.GovUserService{
		GovConfig: x.CommunityGovConfig(),
	}
}

func (x *TestCommunity) CommunityGroupService() group.GovGroupService {
	return group.GovGroupService{
		GovConfig: x.CommunityGovConfig(),
	}
}

func (x *TestCommunity) CommunityMemberService() member.GovMemberService {
	return member.GovMemberService{
		GovConfig: x.CommunityGovConfig(),
	}
}

// user public repos

func (x *TestCommunity) User(i int) string {
	return strconv.Itoa(i)
}

func (x *TestCommunity) UserPublicRepoDir(u int) string {
	return filepath.Join(x.Dir, fmt.Sprintf("%d_public_repo", u))
}

func (x *TestCommunity) UserPublicRepoLocal(u int) git.Local {
	return git.Local{Path: x.UserPublicRepoDir(u)}
}

func (x *TestCommunity) UserPublicRepoURL(u int) string {
	return x.UserPublicRepoDir(u)
}

func (x *TestCommunity) UserIdentityConfig(i int) idproto.IdentityConfig {
	return idproto.IdentityConfig{
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

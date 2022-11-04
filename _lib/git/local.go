package git

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	. "github.com/gov4git/gov4git/lib/base"
	"github.com/gov4git/gov4git/lib/files"
)

type Commit string

type Branch string

var MainBranch = Branch("main")

type URL string

type Origin struct {
	Repo   URL    `json:"repo"`
	Branch Branch `json:"branch"`
}

func (o Origin) String() string {
	return string(o.Repo) + "/" + string(o.Branch)
}

const (
	CommitMsgHeader = "gov4git: "
)

type Local struct {
	// Path is an absolute local path to the git repository
	Path string `json:"path"`
}

func MakeLocal(ctx context.Context) (Local, error) {
	return MakeLocalInCtx(ctx, "")
}

func MakeLocalInCtx(ctx context.Context, label string) (Local, error) {
	eph, err := files.WorkDir(ctx).MkEphemeralDir(label, "")
	if err != nil {
		return Local{}, err
	}
	return LocalInDir(eph), nil
}

func CloneOrigin(ctx context.Context, origin Origin) (Local, error) {
	clone, err := MakeLocal(ctx)
	if err != nil {
		return Local{}, err
	}
	if err := clone.CloneOrigin(ctx, origin); err != nil {
		return Local{}, err
	}
	return clone, nil
}

func LocalInDir(d files.Dir) Local {
	return Local{Path: d.Path}
}

func (x Local) Dir() files.Dir {
	return files.Dir{Path: x.Path}
}

func (x Local) Invoke(ctx context.Context, args ...string) (stdout, stderr string, err error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = x.Path
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout, cmd.Stderr = &outbuf, &errbuf
	err = cmd.Run()
	stdout, stderr = outbuf.String(), errbuf.String()
	Infof("\n$ git %s # repo: %s\nstdout> %s\nstderr> %s", strings.Join(args, " "), x.Path, stdout, stderr)
	return stdout, stderr, err
}

func (x Local) InvokeStdin(ctx context.Context, stdin string, args ...string) (stdout, stderr string, err error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = x.Path
	cmd.Stdin = bytes.NewBufferString(stdin)
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout, cmd.Stderr = &outbuf, &errbuf
	err = cmd.Run()
	stdout, stderr = outbuf.String(), errbuf.String()
	Infof("\n$ git %s # repo: %s\nstdin> %s\nstdout> %s\nstderr> %s", strings.Join(args, " "), x.Path, stdin, stdout, stderr)
	return stdout, stderr, err
}

func (x Local) Version(ctx context.Context) (string, error) {
	stdout, _, err := x.Invoke(ctx, "version")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(stdout), nil
}

func (x Local) Init(ctx context.Context) error {
	err := x.Dir().Mk()
	if err != nil {
		return err
	}
	_, _, err = x.Invoke(ctx, "init")
	return err
}

func (x Local) InitBare(ctx context.Context) error {
	err := x.Dir().Mk()
	if err != nil {
		return err
	}
	_, _, err = x.Invoke(ctx, "init", "--bare")
	return err
}

func (x Local) RenameBranch(ctx context.Context, newBranchName string) error {
	_, _, err := x.Invoke(ctx, "branch", "-M", newBranchName)
	return err
}

func (x Local) Commit(ctx context.Context, msg string) error {
	_, stderr, err1 := x.InvokeStdin(ctx, CommitMsgHeader+msg, "commit", "-F", "-")
	if err2 := ParseCommitError(stderr); err2 != nil {
		return err2
	}
	return err1
}

func (x Local) Commitf(ctx context.Context, f string, args ...any) error {
	_, _, err := x.InvokeStdin(ctx, CommitMsgHeader+fmt.Sprintf(f, args...), "commit", "-F", "-")
	return err
}

func (x Local) AddRemote(ctx context.Context, remoteName string, remoteURL string) error {
	_, _, err := x.Invoke(ctx, "remote", "add", remoteName, remoteURL)
	return err
}

func (x Local) AddRemoteOrigin(ctx context.Context, remoteURL string) error {
	return x.AddRemote(ctx, "origin", remoteURL)
}

func (x Local) PushBranchUpstream(ctx context.Context, srcBranch string) error {
	_, _, err := x.Invoke(ctx, "push", "-u", "origin", srcBranch)
	return err
}

func (x Local) PullBranchUpstream(ctx context.Context, srcBranch string) error {
	_, _, err := x.Invoke(ctx, "pull", "origin", srcBranch)
	return err
}

func (x Local) PushUpstream(ctx context.Context) error {
	return x.PushBranchUpstream(ctx, "HEAD")
}

func (x Local) PullUpstream(ctx context.Context) error {
	return x.PullBranchUpstream(ctx, "HEAD")
}

func (x Local) Push(ctx context.Context) error {
	_, _, err := x.Invoke(ctx, "push", "origin")
	return err
}

func (x Local) Add(ctx context.Context, paths []string) error {
	_, _, err := x.InvokeStdin(ctx, strings.Join(MakeNonAbsPaths(paths), "\n"), "add", "--pathspec-from-file=-")
	return err
}

func (x Local) Remove(ctx context.Context, paths []string) error {
	_, _, err := x.InvokeStdin(ctx, strings.Join(MakeNonAbsPaths(paths), "\n"), "rm", "-r", "--pathspec-from-file=-")
	return err
}

func (x Local) CheckoutBranch(ctx context.Context, branch string) error {
	_, stderr, err1 := x.Invoke(ctx, "checkout", branch)
	if err2 := ParseCheckoutError(stderr); err2 != nil && err2 != ErrAlreadyOnBranch {
		return err2
	}
	return err1
}

func (x Local) CheckoutBranchForce(ctx context.Context, branch string) error {
	_, stderr, err1 := x.Invoke(ctx, "checkout", "-f", branch)
	if err2 := ParseCheckoutError(stderr); err2 != nil && err2 != ErrAlreadyOnBranch {
		return err2
	}
	return err1
}

func (x Local) CheckoutNewBranch(ctx context.Context, branch string) error {
	_, _, err := x.Invoke(ctx, "checkout", "-b", branch)
	return err
}

func (x Local) CheckoutNewOrphan(ctx context.Context, branch string) error {
	_, _, err := x.Invoke(ctx, "checkout", "--orphan", branch)
	return err
}

func (x Local) ResetHard(ctx context.Context) error {
	_, _, err := x.Invoke(ctx, "reset", "--hard")
	return err
}

func (x Local) LogOneline(ctx context.Context) (string, error) {
	stdout, _, err := x.Invoke(ctx, "log", "--pretty=oneline")
	return stdout, err
}

func (x Local) HeadCommitHash(ctx context.Context) (string, error) {
	stdout, _, err := x.Invoke(ctx, "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	h := strings.Trim(stdout, " \t\n\r")
	if h == "" {
		return "", fmt.Errorf("head commit missing")
	}
	return h, nil
}

func (x *Local) CloneOrigin(ctx context.Context, origin Origin) error {
	if err := x.Dir().Mk(); err != nil {
		return nil
	}
	_, stderr, err1 := x.Invoke(ctx, "clone", "--branch", string(origin.Branch), "--single-branch", string(origin.Repo), x.Path)
	if err2 := ParseCloneError(stderr, string(origin.Branch), "origin"); err2 != nil {
		return err2
	}
	return err1
}

func (x Local) CloneOrInitOrigin(ctx context.Context, origin Origin) error {
	if err := x.CloneOrigin(ctx, origin); err != nil {
		if err != ErrRemoteBranchNotFound {
			return err
		}
		if err := x.InitWithRemoteOrigin(ctx, origin); err != nil {
			return err
		}
	}
	return nil
}

func (x Local) InitWithRemoteOrigin(ctx context.Context, origin Origin) error {
	if err := x.Init(ctx); err != nil {
		return err
	}
	if err := x.RenameBranch(ctx, string(origin.Branch)); err != nil {
		return err
	}
	if err := x.AddRemoteOrigin(ctx, string(origin.Repo)); err != nil {
		return err
	}
	return nil
}

func (x Local) AddCommitPush(ctx context.Context, addPaths []string, commitMsg string) error {
	if err := x.Add(ctx, addPaths); err != nil {
		return err
	}
	if err := x.Commit(ctx, commitMsg); err != nil {
		return err
	}
	return x.PushUpstream(ctx)
}

func Init() {
	p, err := exec.LookPath("git")
	if err != nil {
		Fatalf("no path to git")
	}
	r := Local{Path: "."}
	v, err := r.Version(context.Background())
	if err != nil {
		Fatalf("cannot determine git version")
	}
	Infof("using %v, %v", p, v)
}

package git

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/petar/gitty/proto/layout"
	. "github.com/petar/gitty/sys/base"
	"github.com/petar/gitty/sys/files"
)

type Local struct {
	Path string
}

func LocalFromDir(d files.Dir) Local {
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
	Infof("$ git %s\nstdout> %s\nstderr> %s\n", strings.Join(args, " "), stdout, stderr)
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
	Infof("$ git %s\nstdin> %s\nstdout> %s\nstderr> %s\n", strings.Join(args, " "), stdin, stdout, stderr)
	return stdout, stderr, err
}

func (x Local) Init(ctx context.Context) error {
	err := os.MkdirAll(x.Path, 0755)
	if err != nil {
		return err
	}
	_, _, err = x.Invoke(ctx, "init")
	return err
}

func (x Local) RenameBranch(ctx context.Context, newBranchName string) error {
	_, _, err := x.Invoke(ctx, "branch", "-M", newBranchName)
	return err
}

func (x Local) Commit(ctx context.Context, msg string) error {
	_, _, err := x.InvokeStdin(ctx, msg, "commit", "-F", "-")
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

func (x Local) PushUpstream(ctx context.Context) error {
	return x.PushBranchUpstream(ctx, "HEAD")
}

func (x Local) Push(ctx context.Context) error {
	_, _, err := x.Invoke(ctx, "push", "origin")
	return err
}

func (x Local) Add(ctx context.Context, paths []string) error {
	_, _, err := x.InvokeStdin(ctx, strings.Join(paths, "\n"), "add", "--pathspec-from-file=-")
	return err
}

func (x Local) CloneBranch(ctx context.Context, remoteURL, branch string) error {
	if err := x.Dir().Mk(); err != nil {
		return nil
	}
	_, stderr, err1 := x.Invoke(ctx, "clone", "--branch", branch, "--single-branch", remoteURL, x.Path)
	if err2 := ParseCloneError(stderr, branch, "origin"); err2 != nil {
		return err2
	}
	if err1 != nil {
		return err1
	}
	return nil
}

func (x Local) CloneOrInitBranch(ctx context.Context, remoteURL, branch string) error {
	if err := x.CloneBranch(ctx, remoteURL, branch); err != nil {
		if err != ErrRemoteBranchNotFound {
			return err
		}
		if err := x.InitWithRemoteBranch(ctx, remoteURL, branch); err != nil {
			return err
		}
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

func (x Local) InitWithRemoteBranch(ctx context.Context, remoteURL, branch string) error {
	if err := x.Init(ctx); err != nil {
		return err
	}
	if err := x.RenameBranch(ctx, layout.MainBranch); err != nil {
		return err
	}
	if err := x.AddRemoteOrigin(ctx, remoteURL); err != nil {
		return err
	}
	return nil
}

func init() {
	p, err := exec.LookPath("git")
	if err != nil {
		Fatalf("did not find git in path")
	} else {
		Infof("using %s", p)
	}
}

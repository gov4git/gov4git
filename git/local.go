package git

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/petar/gitsoc/files"
)

type Local struct {
	Path string
}

func (x Local) Dir() files.Dir {
	return files.Dir{Path: x.Path}
}

func (x Local) Invoke(ctx context.Context, args ...string) (stdout string, err error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = x.Path
	buf, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func (x Local) InvokeStdin(ctx context.Context, stdin string, args ...string) (stdout string, err error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = x.Path
	cmd.Stdin = bytes.NewBufferString(stdin)
	buf, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func (x Local) Init(ctx context.Context) error {
	err := os.MkdirAll(x.Path, 0755)
	if err != nil {
		return err
	}
	_, err = x.Invoke(ctx, "init")
	return err
}

func (x Local) RenameBranch(ctx context.Context, newBranchName string) error {
	_, err := x.Invoke(ctx, "branch", "-M", newBranchName)
	return err
}

func (x Local) Commit(ctx context.Context, msg string) error {
	_, err := x.InvokeStdin(ctx, msg, "commit", "-F", "-")
	return err
}

func (x Local) AddRemote(ctx context.Context, remoteName string, remoteURL string) error {
	_, err := x.Invoke(ctx, "remote", "add", remoteName, remoteURL)
	return err
}

func (x Local) AddRemoteOrigin(ctx context.Context, remoteURL string) error {
	return x.AddRemote(ctx, "origin", remoteURL)
}

func (x Local) PushToOrigin(ctx context.Context, srcBranch string) error {
	_, err := x.Invoke(ctx, "push", "-u", "origin", srcBranch)
	return err
}

func (x Local) Add(ctx context.Context, paths ...string) error {
	_, err := x.InvokeStdin(ctx, strings.Join(paths, "\n"), "add", "--pathspec-from-file=-")
	return err
}

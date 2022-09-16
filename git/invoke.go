package git

import (
	"bytes"
	"context"
	"os"
	"os/exec"
)

func Invoke(ctx context.Context, dir string, args ...string) (stdout string, err error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	buf, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func InvokeStdin(ctx context.Context, dir string, stdin string, args ...string) (stdout string, err error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	cmd.Stdin = bytes.NewBufferString(stdin)
	buf, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func Init(ctx context.Context, repo string) error {
	err := os.MkdirAll(repo, 0755)
	if err != nil {
		return err
	}
	_, err = Invoke(ctx, repo, "init")
	return err
}

func RenameBranch(ctx context.Context, repo string, newBranchName string) error {
	_, err := Invoke(ctx, repo, "branch", "-M", newBranchName)
	return err
}

func Commit(ctx context.Context, repo string, msg string) error {
	_, err := InvokeStdin(ctx, repo, msg, "commit", "-F", "-")
	return err
}

func AddRemote(ctx context.Context, repo string, remoteName string, remoteURL string) error {
	_, err := Invoke(ctx, repo, "remote", "add", remoteName, remoteURL)
	return err
}

func AddRemoteOrigin(ctx context.Context, repo string, remoteURL string) error {
	return AddRemote(ctx, repo, "origin", remoteURL)
}

func PushToOrigin(ctx context.Context, repo string, srcBranch string) error {
	_, err := Invoke(ctx, repo, "push", "-u", "origin", srcBranch)
	return err
}

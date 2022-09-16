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

type LocalRepo struct {
	Dir string
}

func Init(ctx context.Context, repo LocalRepo) error {
	err := os.MkdirAll(repo.Dir, 0755)
	if err != nil {
		return err
	}
	_, err = Invoke(ctx, repo.Dir, "init")
	return err
}

func RenameBranch(ctx context.Context, repo LocalRepo, newBranchName string) error {
	_, err := Invoke(ctx, repo.Dir, "branch", "-M", newBranchName)
	return err
}

func Commit(ctx context.Context, repo LocalRepo, msg string) error {
	_, err := InvokeStdin(ctx, repo.Dir, msg, "commit", "-F", "-")
	return err
}

func AddRemote(ctx context.Context, repo LocalRepo, name string, url string) error {
	_, err := Invoke(ctx, repo.Dir, "remote", "add", name, url)
	return err
}

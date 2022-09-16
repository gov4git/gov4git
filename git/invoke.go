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

func Init(ctx context.Context, dir string) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	_, err = Invoke(ctx, dir, "init")
	return err
}

func RenameBranch(ctx context.Context, dir, newBranchName string) error {
	_, err := Invoke(ctx, dir, "branch", "-M", newBranchName)
	return err
}

func Commit(ctx context.Context, dir string, msg string) error {
	_, err := InvokeStdin(ctx, dir, msg, "commit", "-F", "-")
	return err
}

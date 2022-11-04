package git

import (
	"context"
	"fmt"
	"testing"
)

func TestRenameMain(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	fmt.Println(dir)
	repo := InitPlain(ctx, dir, false)
	RenameMain(ctx, repo, MainBranch)
	<-(chan int)(nil)
}

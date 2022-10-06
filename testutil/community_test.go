package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCommunityTestInit(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "gov4git_test") // t.TempDir()
	if _, err := CreateTestCommunity(dir, 2); err != nil {
		t.Fatal(err)
	}
}

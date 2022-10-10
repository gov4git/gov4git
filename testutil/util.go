package testutil

import (
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
)

func MakeStickyTestDir(path ...string) string {
	return filepath.Join(os.TempDir(), filepath.Join(path...), strconv.FormatUint(uint64(rand.Int63()), 36))
}

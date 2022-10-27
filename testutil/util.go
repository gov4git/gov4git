package testutil

import (
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func MakeStickyTestDir(path ...string) string {
	rand.Seed(time.Now().Unix())
	return filepath.Join(os.TempDir(), filepath.Join(path...), strconv.FormatUint(uint64(rand.Int63()), 36))
}

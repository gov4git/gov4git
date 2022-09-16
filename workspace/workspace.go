package workspace

import (
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Workspace struct {
	Dir string
}

func (x Workspace) AbsPath(rel string) string {
	return filepath.Join(x.Dir, rel)
}

func (x Workspace) MakeEphemeralDir(prefix string) (abs string, err error) {
	t := time.Now()
	abs = filepath.Join(
		"ephemeral",
		strconv.Itoa(t.Year()),
		strconv.Itoa(int(t.Month())),
		strconv.Itoa(t.Day()),
		strconv.Itoa(t.Hour()),
		strings.Join([]string{prefix, strconv.FormatUint(uint64(rand.Int63()), 64)}, "."),
	)
	if err = os.MkdirAll(abs, 0755); err != nil {
		return "", err
	}
	return abs, nil
}

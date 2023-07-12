package balance

import (
	"path/filepath"
	"strings"
)

type Balance string

func userPropKey(balanceKey Balance) string {
	key := strings.TrimLeft(filepath.ToSlash(string(balanceKey)), "/")
	return strings.Join([]string{"balance", key}, "/")
}

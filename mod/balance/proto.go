package balance

import "path/filepath"

func userPropKey(balanceKey string) string {
	return filepath.Join("balance", balanceKey)
}

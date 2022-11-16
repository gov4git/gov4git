package balance

import "path/filepath"

type Balance string

func userPropKey(balanceKey Balance) string {
	return filepath.Join("balance", string(balanceKey))
}

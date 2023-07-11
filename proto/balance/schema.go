package balance

import (
	"strings"
)

type Balance string

func userPropKey(balanceKey Balance) string {
	key := strings.TrimLeft(string(balanceKey), "/")
	return strings.Join([]string{"balance", key}, "/")
}

package govproto

import "fmt"

func BalanceKey(balance string) string {
	return fmt.Sprintf("balance:%v", balance)
}

package member

import "github.com/gov4git/gov4git/mod/id"

type Account struct {
	Home id.PublicAddress `json:"home"`
}

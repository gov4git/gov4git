package git

import "github.com/gov4git/gov4git/lib/form"

type Change[R any] struct {
	Result R
	Msg    string //commit msg
}

type ChangeNoResult = Change[form.None]

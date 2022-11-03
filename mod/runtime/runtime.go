package runtime

type Change[R any] struct {
	Result R
	Msg    string //commit msg
}

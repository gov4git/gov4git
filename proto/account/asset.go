package account

type Asset string

func (a Asset) String() string {
	return string(a)
}

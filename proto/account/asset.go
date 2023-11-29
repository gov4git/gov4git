package account

var (
	PluralAsset = Asset("plural")
)

type Asset string

func (a Asset) String() string {
	return string(a)
}

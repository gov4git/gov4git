package account

import "github.com/gov4git/gov4git/v2/proto/history"

var (
	PluralAsset = Asset("plural")
)

type Asset string

func (a Asset) String() string {
	return string(a)
}

func (a Asset) HistoryAsset() history.Asset {
	return history.Asset(a)
}

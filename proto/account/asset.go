package account

import "github.com/gov4git/gov4git/v2/proto/history/metric"

var (
	PluralAsset = Asset("plural")
)

type Asset string

func (a Asset) String() string {
	return string(a)
}

func (a Asset) MetricAsset() metric.Asset {
	return metric.Asset(a)
}

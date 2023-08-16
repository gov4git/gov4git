package qv

type QV struct{}

const QVStrategyName = "qv"

func (x QV) Name() string {
	return QVStrategyName
}

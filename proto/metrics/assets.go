package metrics

type AssetURLCalculator func(assetRepoPath string) (url string)

type ReportAssets struct {
	Series   *ReportSeries
	ReportMD string
	Assets   map[string][]byte // path in git repo assets branch -> content
}

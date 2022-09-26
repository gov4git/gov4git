package git

import "testing"

func TestURL2Path(t *testing.T) {
	for _, testURL := range testURL2Path {
		u2p, err := RepoURLToPath(testURL[0])
		if err != nil {
			t.Error(err)
		} else if u2p != testURL[1] {
			t.Errorf("expecting %v, got %v", u2p, testURL[1])
		}
	}
}

var testURL2Path = [][2]string{
	{"git@github.com:ipfs/kubo.git", "github.com/ipfs/kubo.git"},
	{"https://github.com/petar/maymounkov.org.git", "github.com/petar/maymounkov.org.git"},
}

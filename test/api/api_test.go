package api

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/gov4git/gov4git/v2/gov4git/cmd"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	resRun := testscript.RunMain(m, map[string]func() int{
		"gov4git": func() int {
			return cmd.ExecuteWithConfig("config.json")
		},
	})
	os.Exit(resRun)
}

func TestScript(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/script",
		Cmds: map[string]func(ts *testscript.TestScript, neg bool, args []string){
			"provision": scriptProvisionCommunity,
		},
	})
}

func scriptProvisionCommunity(ts *testscript.TestScript, neg bool, args []string) {

	ts.Setenv("HOME", ts.Getenv("WORK"))
	ts.Exec("git", "config", "--global", "init.defaultBranch", "main")
	cachePath := path.Join(ts.Getenv("WORK"), "cache")
	os.MkdirAll(cachePath, 0755)

	govPublicPath := path.Join(ts.Getenv("WORK"), "gov_public")
	os.MkdirAll(govPublicPath, 0755)
	ts.Exec("git", "init", "--bare", govPublicPath)

	govPrivatePath := path.Join(ts.Getenv("WORK"), "gov_private")
	os.MkdirAll(govPrivatePath, 0755)
	ts.Exec("git", "init", "--bare", govPrivatePath)

	memberPublicPath := path.Join(ts.Getenv("WORK"), "member_public")
	os.MkdirAll(memberPublicPath, 0755)
	ts.Exec("git", "init", "--bare", memberPublicPath)

	memberPrivatePath := path.Join(ts.Getenv("WORK"), "member_private")
	os.MkdirAll(memberPrivatePath, 0755)
	ts.Exec("git", "init", "--bare", memberPrivatePath)

	// create config.json
	configPath := path.Join(ts.Getenv("WORK"), "config.json")
	src := fmt.Sprintf(srcConfigJSONFmt,
		"cache",
		"gov_public",
		"gov_private",
		"member_public",
		"member_private",
	)
	fmt.Println("CONFIG:", src)
	os.WriteFile(configPath, []byte(src), 0644)
}

const (
	srcConfigJSONFmt = `

	{
		"cache_dir": %q,
		"gov_public_url": %q,
		"gov_public_branch": "main",
		"gov_private_url": %q,
		"gov_private_branch": "main",
		"member_public_url": %q,
		"member_public_branch": "main",
		"member_private_url": %q,
		"member_private_branch": "main"
	}

`
)

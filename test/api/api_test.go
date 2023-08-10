package api

import (
	"os"
	"path"
	"testing"

	"github.com/gov4git/gov4git/gov4git/cmd"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	resRun := testscript.RunMain(m, map[string]func() int{
		"gov4git": cmd.Execute,
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

	// create config.json
	configPath := path.Join(ts.Getenv("WORK"), "config.json")
	os.WriteFile(configPath, []byte(srcConfigJSON), 0644)

	ts.Setenv("HOME", ts.Getenv("WORK"))
	ts.Exec("git", "config", "--global", "init.defaultBranch", "main")
	os.MkdirAll(path.Join(ts.Getenv("WORK"), "cache"), 0755)

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
}

const (
	srcConfigJSON = `
{
	"cache_dir": "$WORK/cache",
	"gov_public_url": "~/gov_public",
	"gov_public_branch": "main",
	"gov_private_url": "~/gov_private",
	"gov_private_branch": "main",
	"member_public_url": "~/member_public",
	"member_public_branch": "main",
	"member_private_url": "~/member_private",
	"member_private_branch": "main"
}
`
)

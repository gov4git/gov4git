package runtime

import (
	"context"
	"os/exec"

	"github.com/gov4git/lib4git/base"
)

func RunPython(ctx context.Context, script string) ([]byte, error) {
	py3path, err := exec.LookPath("python3")
	if err != nil {
		return nil, err
	}
	return RunPythonWithPath(ctx, py3path, script)
}

func RunPythonWithPath(ctx context.Context, exepath, script string) ([]byte, error) {
	base.Infof("running %s -c ...", exepath)
	cmd := exec.CommandContext(ctx, exepath, `-c`, script)
	return cmd.CombinedOutput()
}

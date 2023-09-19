//go:build windows

package runtime

import (
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
	"github.com/go-git/go-git/v5/plumbing/transport/server"
)

// init ensures that go-git does not use an external git binary for git file urls.
func init() {
	client.InstallProtocol("file", server.NewClient(server.NewFilesystemLoader(osfs.New(""))))
}

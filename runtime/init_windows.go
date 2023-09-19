//go:build windows

package runtime

// init ensures that go-git does not use an external git binary for git file urls.
func init() {
	// client.InstallProtocol("file", server.NewClient(server.NewFilesystemLoader(osfs.New(""))))
}

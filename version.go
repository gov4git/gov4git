// Package gov4git_root embeds VERSION.txt into the binary.
package gov4git_root

import _ "embed"

// VersionDotTxt is the contents of VERSION.txt.
//
//go:embed VERSION.txt
var VersionDotTxt string

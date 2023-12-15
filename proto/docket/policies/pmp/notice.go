package pmp

import (
	"fmt"

	"github.com/gov4git/gov4git/v2/materials"
)

var Welcome = fmt.Sprintf(
	`

This project is managed by [Gov4Git](%s), a decentralized governance system for collaborative git projects.
To participate in governance, install the [Gov4Git desktop app](%s).
	`,
	materials.Gov4GitWebsiteURL,
	materials.Gov4GitDesktopAppInstall,
)

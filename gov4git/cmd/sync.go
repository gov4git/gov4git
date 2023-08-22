package cmd

import (
	"fmt"
	"os"

	"github.com/gov4git/gov4git/proto/sync"
	"github.com/gov4git/lib4git/form"
	"github.com/spf13/cobra"
)

var (
	syncCmd = &cobra.Command{
		Use:   "sync",
		Short: "Sync governance with the community",
		Long: `
Sync is the heartbeat that advances the state of the governance forward.
Sync fetches all outstanding votes from community users and incorporates them in ballot tallies.
Sync also fetches and processes all outstanding service requests from community users.`,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			chg := sync.Sync(ctx, setup.Organizer, syncMaxPar)
			fmt.Fprint(os.Stdout, form.SprintJSON(chg.Result))
		},
	}
)

var (
	syncMaxPar int
)

func init() {
	syncCmd.Flags().IntVar(&syncMaxPar, "maxpar", 5, "parallelism while clonging member repos for vote collection")
}

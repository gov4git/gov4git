package cmd

import (
	"fmt"
	"os"

	"github.com/petar/gov4git/lib/base"
	"github.com/petar/gov4git/lib/files"
	man "github.com/petar/gov4git/man/arb"
	"github.com/petar/gov4git/proto"
	"github.com/petar/gov4git/services/gov/arb"
	"github.com/spf13/cobra"
)

var (
	pollCmd = &cobra.Command{
		Use:   "poll",
		Short: "Create a new poll",
		Long:  man.GovArbPoll,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := arb.GovArbService{
				GovConfig: proto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(proto.LocalAgentTempPath, "gov-arb-poll")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.Poll(ctx, &arb.GovArbPollIn{
				Path:            pollPath,
				Choices:         pollChoices,
				Group:           pollGroup,
				Strategy:        pollStrategy,
				GoverningBranch: pollGoverningBranch,
			})
			if err == nil {
				fmt.Fprint(os.Stdout, r.Human(cmd.Context()))
			} else {
				fmt.Fprint(os.Stderr, err.Error())
			}
			return err
		},
	}
)

var (
	pollPath            string
	pollChoices         []string
	pollGroup           string
	pollStrategy        string
	pollGoverningBranch string
)

func init() {
	pollCmd.Flags().StringVar(&pollPath, "path", "", "community repo path for poll results and proofs")
	pollCmd.Flags().StringArrayVar(&pollChoices, "choices", nil, "poll choices")
	pollCmd.Flags().StringVar(&pollGroup, "group", "", "group of users participating in poll")
	pollCmd.Flags().StringVar(&pollStrategy, "strategy", "", "polling strategy (XXX)")
	pollCmd.Flags().StringVar(&pollGoverningBranch, "branch", "", "branch governing the poll")
}

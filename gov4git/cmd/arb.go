package cmd

import (
	"fmt"
	"os"

	"github.com/gov4git/gov4git/lib/base"
	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/form"
	man "github.com/gov4git/gov4git/man/arb"
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/services/gov/arb"
	"github.com/spf13/cobra"
)

var (
	pollCmd = &cobra.Command{
		Use:   "poll",
		Short: "Create a new poll",
		Long:  man.Poll,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := arb.GovArbService{
				GovConfig: proto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(proto.LocalAgentTempPath, "gov-arb-poll")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.Poll(ctx, &arb.PollIn{
				Path:            pollPath,
				Choices:         pollChoices,
				Group:           pollGroup,
				Strategy:        pollStrategy,
				GoverningBranch: pollGoverningBranch,
				PollBranch:      pollBranch,
			})
			if err == nil {
				fmt.Fprint(os.Stdout, form.Pretty(r))
			} else {
				fmt.Fprint(os.Stderr, err.Error())
			}
			return err
		},
	}

	voteCmd = &cobra.Command{
		Use:   "vote",
		Short: "Vote on a referendum",
		Long:  man.Vote,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := arb.GovArbService{
				GovConfig: proto.GovConfig{
					CommunityURL: communityURL,
				},
				IdentityConfig: proto.IdentityConfig{
					PublicURL:  publicURL,
					PrivateURL: privateURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(proto.LocalAgentTempPath, "gov-arb-vote")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.Vote(ctx, &arb.VoteIn{
				ReferendumBranch: voteReferendumBranch,
				ReferendumPath:   voteReferendumPath,
				VoteChoice:       voteChoice,
				VoteStrength:     voteStrength,
			})
			if err == nil {
				fmt.Fprint(os.Stdout, form.Pretty(r))
			} else {
				fmt.Fprint(os.Stderr, err.Error())
			}
			return err
		},
	}

	tallyCmd = &cobra.Command{
		Use:   "tally",
		Short: "Tally votes on a referendum",
		Long:  man.Tally,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := arb.GovArbService{
				GovConfig: proto.GovConfig{
					CommunityURL: communityURL,
				},
				IdentityConfig: proto.IdentityConfig{
					PublicURL:  publicURL,
					PrivateURL: privateURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(proto.LocalAgentTempPath, "gov-arb-tally")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.Tally(ctx, &arb.TallyIn{
				ReferendumBranch: tallyReferendumBranch,
				ReferendumPath:   tallyReferendumPath,
			})
			if err == nil {
				fmt.Fprint(os.Stdout, form.Pretty(r))
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
	pollBranch          string

	voteReferendumBranch string
	voteReferendumPath   string
	voteChoice           string
	voteStrength         float64

	tallyReferendumBranch string
	tallyReferendumPath   string
)

func init() {
	pollCmd.Flags().StringVar(&pollPath, "path", "", "community repo path for poll results and proofs")
	pollCmd.Flags().StringArrayVar(&pollChoices, "choices", nil, "poll choices")
	pollCmd.Flags().StringVar(&pollGroup, "group", "", "group of users participating in poll")
	pollCmd.Flags().StringVar(&pollStrategy, "strategy", "", "polling strategy (available strategy: prioritize)")
	pollCmd.Flags().StringVar(&pollGoverningBranch, "govern-branch", "", "branch governing the poll")
	pollCmd.Flags().StringVar(&pollBranch, "poll-branch", "", "branch where poll is created (if empty, use governing branch)")

	voteCmd.Flags().StringVar(&voteReferendumBranch, "--refm-branch", "", "referendum branch (e.g. poll branch)")
	voteCmd.Flags().StringVar(&voteReferendumPath, "--refm-path", "", "referendum path")
	voteCmd.Flags().StringVar(&voteChoice, "--choice", "", "vote choice")
	voteCmd.Flags().Float64Var(&voteStrength, "--strength", 0, "vote strength")

	tallyCmd.Flags().StringVar(&tallyReferendumBranch, "--refm-branch", "", "referendum branch (e.g. poll branch)")
	tallyCmd.Flags().StringVar(&tallyReferendumPath, "--refm-path", "", "referendum path")
}

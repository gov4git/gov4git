package cmd

import (
	"fmt"
	"os"

	"github.com/gov4git/gov4git/lib/base"
	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/form"
	man "github.com/gov4git/gov4git/man/arb"
	"github.com/gov4git/gov4git/proto/cmdproto"
	"github.com/gov4git/gov4git/proto/govproto"
	"github.com/gov4git/gov4git/proto/idproto"
	"github.com/gov4git/gov4git/services/gov/arb"
	"github.com/spf13/cobra"
)

var (
	ballotCmd = &cobra.Command{
		Use:   "ballot",
		Short: "Create a new ballot (poll, merge proposal, etc.)",
		Long:  man.Ballot,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := arb.GovArbService{
				GovConfig: govproto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(cmdproto.LocalAgentTempPath, "gov-arb-ballot")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.CreateBallot(ctx, &arb.CreateBallotIn{
				Path:            ballotPath,
				Choices:         ballotChoices,
				Group:           ballotGroup,
				Strategy:        ballotStrategy,
				GoverningBranch: ballotGoverningBranch,
				BallotBranch:    ballotBranch,
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
		Short: "Vote on a ballot",
		Long:  man.Vote,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := arb.GovArbService{
				GovConfig: govproto.GovConfig{
					CommunityURL: communityURL,
				},
				IdentityConfig: idproto.IdentityConfig{
					PublicURL:  publicURL,
					PrivateURL: privateURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(cmdproto.LocalAgentTempPath, "gov-arb-vote")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)

			if len(voteChoice) != len(voteStrength) {
				err := fmt.Errorf("number of choices much match number of strengths")
				fmt.Fprint(os.Stderr, err.Error())
				return err
			}
			votes := make([]govproto.Election, len(voteChoice))
			for i := 0; i < len(voteChoice); i++ {
				votes[i] = govproto.Election{Choice: voteChoice[i], Strength: voteStrength[i]}
			}

			r, err := s.Vote(ctx, &arb.VoteIn{
				BallotBranch: voteBallotBranch,
				BallotPath:   voteBallotPath,
				Votes:        votes,
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
		Short: "Tally votes on a ballot",
		Long:  man.Tally,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := arb.GovArbService{
				GovConfig: govproto.GovConfig{
					CommunityURL: communityURL,
				},
				IdentityConfig: idproto.IdentityConfig{
					PublicURL:  publicURL,
					PrivateURL: privateURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(cmdproto.LocalAgentTempPath, "gov-arb-tally")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.Tally(ctx, &arb.TallyIn{
				BallotBranch: tallyBallotBranch,
				BallotPath:   tallyBallotPath,
			})
			if err == nil {
				fmt.Fprint(os.Stdout, form.Pretty(r))
			} else {
				fmt.Fprint(os.Stderr, err.Error())
			}
			return err
		},
	}

	sealCmd = &cobra.Command{
		Use:   "seal",
		Short: "Seal a ballot at the current tally.",
		Long:  man.Seal,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := arb.GovArbService{
				GovConfig: govproto.GovConfig{
					CommunityURL: communityURL,
				},
				IdentityConfig: idproto.IdentityConfig{
					PublicURL:  publicURL,
					PrivateURL: privateURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(cmdproto.LocalAgentTempPath, "gov-arb-seal")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.Seal(ctx, &arb.SealIn{
				BallotBranch: tallyBallotBranch,
				BallotPath:   tallyBallotPath,
			})
			if err == nil {
				fmt.Fprint(os.Stdout, form.Pretty(r))
			} else {
				fmt.Fprint(os.Stderr, err.Error())
			}
			return err
		},
	}

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List currently open and closed ballots.",
		Long:  man.List,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := arb.GovArbService{
				GovConfig: govproto.GovConfig{
					CommunityURL: communityURL,
				},
				IdentityConfig: idproto.IdentityConfig{
					PublicURL:  publicURL,
					PrivateURL: privateURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(cmdproto.LocalAgentTempPath, "gov-arb-list")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.List(ctx, &arb.ListIn{
				BallotBranch: tallyBallotBranch,
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
	ballotPath            string
	ballotChoices         []string
	ballotGroup           string
	ballotStrategy        string
	ballotGoverningBranch string
	ballotBranch          string

	voteBallotBranch string
	voteBallotPath   string
	voteChoice       []string
	voteStrength     []float64

	tallyBallotBranch string
	tallyBallotPath   string

	sealBallotBranch string
	sealBallotPath   string

	listBallotBranch string
)

func init() {
	ballotCmd.Flags().StringVar(&ballotPath, "path", "", "community repo path for ballot results and proofs")
	ballotCmd.Flags().StringArrayVar(&ballotChoices, "choices", nil, "ballot choices")
	ballotCmd.Flags().StringVar(&ballotGroup, "group", "", "group of users participating in ballot")
	ballotCmd.Flags().StringVar(&ballotStrategy, "strategy", "", "balloting strategy (available strategy: priority-poll)")
	ballotCmd.Flags().StringVar(&ballotGoverningBranch, "govern-branch", "", "branch governing the ballot")
	ballotCmd.Flags().StringVar(&ballotBranch, "ballot-branch", "", "branch where ballot is created (if empty, use governing branch)")

	voteCmd.Flags().StringVar(&voteBallotBranch, "ballot-branch", "", "ballot branch")
	voteCmd.Flags().StringVar(&voteBallotPath, "ballot-path", "", "ballot path")
	voteCmd.Flags().StringSliceVar(&voteChoice, "choice", nil, "vote choice")
	voteCmd.Flags().Float64SliceVar(&voteStrength, "strength", nil, "vote strength")

	tallyCmd.Flags().StringVar(&tallyBallotBranch, "ballot-branch", "", "ballot branch")
	tallyCmd.Flags().StringVar(&tallyBallotPath, "ballot-path", "", "ballot path")

	sealCmd.Flags().StringVar(&sealBallotBranch, "ballot-branch", "", "ballot branch")
	sealCmd.Flags().StringVar(&sealBallotPath, "ballot-path", "", "ballot path")

	listCmd.Flags().StringVar(&listBallotBranch, "ballot-branch", "", "ballot branch")
}

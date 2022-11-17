package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/gov4git/gov4git/mod/ballot/core"
	"github.com/gov4git/gov4git/mod/ballot/proto"
	"github.com/gov4git/gov4git/mod/ballot/qv"
	"github.com/gov4git/gov4git/mod/member"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
	"github.com/spf13/cobra"
)

var (
	ballotCmd = &cobra.Command{
		Use:   "ballot",
		Short: "Manage ballots",
		Long:  ``,
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	ballotOpenCmd = &cobra.Command{
		Use:   "open",
		Short: "Open a new ballot",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			strat := qv.PriorityPoll{UseVotingCredits: ballotUseVotingCredits}
			chg := core.Open(
				ctx,
				strat,
				setup.Community,
				ns.NS(ballotName),
				ballotTitle,
				ballotDescription,
				ballotChoices,
				member.Group(ballotGroup),
			)
			fmt.Fprint(os.Stdout, form.Pretty(chg.Result))
		},
	}

	ballotCloseCmd = &cobra.Command{
		Use:   "close",
		Short: "Close an open ballot",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			chg := core.Close(
				ctx,
				setup.Community,
				ns.NS(ballotName),
			)
			fmt.Fprint(os.Stdout, form.Pretty(chg.Result))
		},
	}

	ballotListCmd = &cobra.Command{
		Use:   "list",
		Short: "List open ballots",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			ls := core.ListOpen(
				ctx,
				setup.Community,
			)
			fmt.Fprint(os.Stdout, form.Pretty(ls))
		},
	}

	ballotTallyCmd = &cobra.Command{
		Use:   "tally",
		Short: "Fetch current votes and record latest tally",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			chg := core.Tally(
				ctx,
				setup.Organizer,
				ns.NS(ballotName),
			)
			fmt.Fprint(os.Stdout, form.Pretty(chg.Result))
		},
	}

	ballotVoteCmd = &cobra.Command{
		Use:   "vote",
		Short: "Cast a vote on an open ballot",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			chg := core.Vote(
				ctx,
				setup.Member,
				setup.Community,
				ns.NS(ballotName),
				parseElections(ctx, ballotElectionChoice, ballotElectionStrength),
			)
			fmt.Fprint(os.Stdout, form.Pretty(chg.Result))
		},
	}
)

var (
	ballotName             string
	ballotTitle            string
	ballotDescription      string
	ballotChoices          []string
	ballotGroup            string
	ballotElectionChoice   []string
	ballotElectionStrength []float64
	ballotUseVotingCredits bool
)

func init() {
	// open
	ballotCmd.AddCommand(ballotOpenCmd)
	ballotOpenCmd.Flags().StringVar(&ballotName, "name", "", "ballot name")
	ballotOpenCmd.MarkFlagRequired("name")
	ballotOpenCmd.Flags().StringVar(&ballotTitle, "title", "", "ballot title")
	ballotOpenCmd.MarkFlagRequired("title")
	ballotOpenCmd.Flags().StringVar(&ballotDescription, "desc", "", "ballot description")
	ballotOpenCmd.MarkFlagRequired("desc")
	ballotOpenCmd.Flags().StringSliceVar(&ballotChoices, "choices", nil, "ballot choices")
	ballotOpenCmd.MarkFlagRequired("choices")
	ballotOpenCmd.Flags().StringVar(&ballotGroup, "group", "", "group of ballot participants")
	ballotOpenCmd.MarkFlagRequired("group")
	ballotOpenCmd.Flags().BoolVar(&ballotUseVotingCredits, "use_credits", false, "use voting credits")

	// close
	ballotCmd.AddCommand(ballotCloseCmd)
	ballotCloseCmd.Flags().StringVar(&ballotName, "name", "", "ballot name")
	ballotCloseCmd.MarkFlagRequired("name")

	// list
	ballotCmd.AddCommand(ballotListCmd)

	// tally
	ballotCmd.AddCommand(ballotTallyCmd)
	ballotTallyCmd.Flags().StringVar(&ballotName, "name", "", "ballot name")
	ballotTallyCmd.MarkFlagRequired("name")

	// vote
	ballotCmd.AddCommand(ballotVoteCmd)
	ballotVoteCmd.Flags().StringVar(&ballotName, "name", "", "ballot name")
	ballotVoteCmd.MarkFlagRequired("name")
	ballotVoteCmd.Flags().StringSliceVar(&ballotElectionChoice, "choices", nil, "list of elected choices")
	ballotVoteCmd.MarkFlagRequired("choices")
	ballotVoteCmd.Flags().Float64SliceVar(&ballotElectionStrength, "strengths", nil, "list of elected vote strengths")
	ballotVoteCmd.MarkFlagRequired("strengths")
}

func parseElections(ctx context.Context, choices []string, strengths []float64) proto.Elections {
	if len(choices) != len(strengths) {
		must.Errorf(ctx, "elected choices must match elected strengths in count")
	}
	el := make(proto.Elections, len(choices))
	for i := range choices {
		el[i] = proto.Election{VoteChoice: choices[i], VoteStrengthChange: strengths[i]}
	}
	return el
}

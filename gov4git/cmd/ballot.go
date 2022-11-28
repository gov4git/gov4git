package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/gov4git/gov4git/proto/ballot/ballot"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/ballot/qv"
	"github.com/gov4git/gov4git/proto/member"
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
			chg := ballot.Open(
				ctx,
				strat,
				setup.Gov,
				ns.ParseFromPath(ballotName),
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
			chg := ballot.Close(
				ctx,
				setup.Organizer,
				ns.ParseFromPath(ballotName),
				common.Summary(ballotSummary),
			)
			fmt.Fprint(os.Stdout, form.Pretty(chg.Result))
		},
	}

	ballotShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show open ballot",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			r := ballot.Show(
				ctx,
				setup.Gov,
				ns.ParseFromPath(ballotName),
			)
			fmt.Fprint(os.Stdout, r)
		},
	}

	ballotListCmd = &cobra.Command{
		Use:   "list",
		Short: "List open ballots",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			ls := ballot.ListOpen(
				ctx,
				setup.Gov,
			)
			fmt.Fprint(os.Stdout, form.Pretty(ls))
		},
	}

	ballotTallyCmd = &cobra.Command{
		Use:   "tally",
		Short: "Fetch current votes and record latest tally",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			chg := ballot.Tally(
				ctx,
				setup.Organizer,
				ns.ParseFromPath(ballotName),
			)
			fmt.Fprint(os.Stdout, form.Pretty(chg.Result))
		},
	}

	ballotVoteCmd = &cobra.Command{
		Use:   "vote",
		Short: "Cast a vote on an open ballot",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			chg := ballot.Vote(
				ctx,
				setup.Member,
				setup.Gov,
				ns.ParseFromPath(ballotName),
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
	ballotSummary          string
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
	ballotCmd.AddCommand(ballotShowCmd)
	ballotShowCmd.Flags().StringVar(&ballotName, "name", "", "ballot name")
	ballotShowCmd.MarkFlagRequired("name")
	ballotShowCmd.Flags().StringVar(&ballotSummary, "summary", "", "summary")
	ballotShowCmd.MarkFlagRequired("summary")

	// show
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

func parseElections(ctx context.Context, choices []string, strengths []float64) common.Elections {
	if len(choices) != len(strengths) {
		must.Errorf(ctx, "elected choices must match elected strengths in count")
	}
	el := make(common.Elections, len(choices))
	for i := range choices {
		el[i] = common.Election{VoteChoice: choices[i], VoteStrengthChange: strengths[i]}
	}
	return el
}

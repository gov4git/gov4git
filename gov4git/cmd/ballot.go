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
			LoadConfig()
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
			fmt.Fprint(os.Stdout, form.SprintJSON(chg.Result))
		},
	}

	ballotFreezeCmd = &cobra.Command{
		Use:   "freeze",
		Short: "Freeze an open ballot",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			chg := ballot.Freeze(
				ctx,
				setup.Organizer,
				ns.ParseFromPath(ballotName),
			)
			fmt.Fprint(os.Stdout, form.SprintJSON(chg.Result))
		},
	}

	ballotUnfreezeCmd = &cobra.Command{
		Use:   "unfreeze",
		Short: "Unfreeze an open ballot",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			chg := ballot.Unfreeze(
				ctx,
				setup.Organizer,
				ns.ParseFromPath(ballotName),
			)
			fmt.Fprint(os.Stdout, form.SprintJSON(chg.Result))
		},
	}

	ballotCloseCmd = &cobra.Command{
		Use:   "close",
		Short: "Close an open ballot",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			chg := ballot.Close(
				ctx,
				setup.Organizer,
				ns.ParseFromPath(ballotName),
				common.Summary(ballotSummary),
			)
			fmt.Fprint(os.Stdout, form.SprintJSON(chg.Result))
		},
	}

	ballotShowOpenCmd = &cobra.Command{
		Use:   "show-open",
		Short: "Show open ballot",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			r := ballot.Show(
				ctx,
				setup.Gov,
				ns.ParseFromPath(ballotName),
				false,
			)
			fmt.Fprint(os.Stdout, form.SprintJSON(r))
		},
	}

	ballotShowClosedCmd = &cobra.Command{
		Use:   "show-closed",
		Short: "Show closed ballot",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			r := ballot.Show(
				ctx,
				setup.Gov,
				ns.ParseFromPath(ballotName),
				true,
			)
			fmt.Fprint(os.Stdout, form.SprintJSON(r))
		},
	}

	ballotListOpenCmd = &cobra.Command{
		Use:   "list-open",
		Short: "List open ballots",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			ls := ballot.List(
				ctx,
				setup.Gov,
				false,
			)
			if ballotOnlyNames {
				for _, n := range common.AdsToBallotNames(ls) {
					fmt.Println(n)
				}
			} else {
				fmt.Fprint(os.Stdout, form.SprintJSON(ls))
			}
		},
	}

	ballotListClosedCmd = &cobra.Command{
		Use:   "list-closed",
		Short: "List closed ballots",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			ls := ballot.List(
				ctx,
				setup.Gov,
				true,
			)
			if ballotOnlyNames {
				for _, n := range common.AdsToBallotNames(ls) {
					fmt.Println(n)
				}
			} else {
				fmt.Fprint(os.Stdout, form.SprintJSON(ls))
			}
		},
	}

	ballotTallyCmd = &cobra.Command{
		Use:   "tally",
		Short: "Fetch current votes and record latest tally",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			chg := ballot.Tally(
				ctx,
				setup.Organizer,
				ns.ParseFromPath(ballotName),
			)
			fmt.Fprint(os.Stdout, form.SprintJSON(chg.Result))
		},
	}

	ballotVoteCmd = &cobra.Command{
		Use:   "vote",
		Short: "Cast a vote on an open ballot",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			chg := ballot.Vote(
				ctx,
				setup.Member,
				setup.Gov,
				ns.ParseFromPath(ballotName),
				parseElections(ctx, ballotElectionChoice, ballotElectionStrength),
			)
			fmt.Fprint(os.Stdout, form.SprintJSON(chg.Result))
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
	ballotOnlyNames        bool
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
	ballotCloseCmd.Flags().StringVar(&ballotSummary, "summary", "", "summary")
	ballotCloseCmd.MarkFlagRequired("summary")

	// freeze
	ballotCmd.AddCommand(ballotFreezeCmd)
	ballotFreezeCmd.Flags().StringVar(&ballotName, "name", "", "ballot name")
	ballotFreezeCmd.MarkFlagRequired("name")

	// unfreeze
	ballotCmd.AddCommand(ballotUnfreezeCmd)
	ballotUnfreezeCmd.Flags().StringVar(&ballotName, "name", "", "ballot name")
	ballotUnfreezeCmd.MarkFlagRequired("name")

	// show open
	ballotCmd.AddCommand(ballotShowOpenCmd)
	ballotShowOpenCmd.Flags().StringVar(&ballotName, "name", "", "ballot name")
	ballotShowOpenCmd.MarkFlagRequired("name")

	// show closed
	ballotCmd.AddCommand(ballotShowClosedCmd)
	ballotShowClosedCmd.Flags().StringVar(&ballotName, "name", "", "ballot name")
	ballotShowClosedCmd.MarkFlagRequired("name")

	// list open
	ballotCmd.AddCommand(ballotListOpenCmd)
	ballotListOpenCmd.Flags().BoolVar(&ballotOnlyNames, "only_names", false, "list only ballot names")

	// list closed
	ballotCmd.AddCommand(ballotListClosedCmd)
	ballotListClosedCmd.Flags().BoolVar(&ballotOnlyNames, "only_names", false, "list only ballot names")

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

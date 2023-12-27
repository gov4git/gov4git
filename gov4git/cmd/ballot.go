package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotapi"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotio"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/gov4git/v2/proto/purpose"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/must"
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
			chg := ballotapi.Open(
				ctx,
				ballotio.QVStrategyName,
				setup.Organizer,
				ballotproto.ParseBallotNameFromPath(ballotName),
				account.NobodyAccountID,
				purpose.Unspecified,
				"",
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
			chg := ballotapi.Freeze(
				ctx,
				setup.Organizer,
				ballotproto.ParseBallotNameFromPath(ballotName),
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
			chg := ballotapi.Unfreeze(
				ctx,
				setup.Organizer,
				ballotproto.ParseBallotNameFromPath(ballotName),
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
			chg := ballotapi.Close(
				ctx,
				setup.Organizer,
				ballotproto.ParseBallotNameFromPath(ballotName),
				account.AccountID(ballotEscrowTo),
			)
			fmt.Fprint(os.Stdout, form.SprintJSON(chg.Result))
		},
	}

	ballotCancelCmd = &cobra.Command{
		Use:   "cancel",
		Short: "Cancel an open ballot",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			chg := ballotapi.Cancel(
				ctx,
				setup.Organizer,
				ballotproto.ParseBallotNameFromPath(ballotName),
			)
			fmt.Fprint(os.Stdout, form.SprintJSON(chg.Result))
		},
	}

	ballotShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show ballot details",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			r := ballotapi.Show(
				ctx,
				setup.Gov,
				ballotproto.ParseBallotNameFromPath(ballotName),
			)
			fmt.Fprint(os.Stdout, form.SprintJSON(r))
		},
	}

	ballotListCmd = &cobra.Command{
		Use:   "list",
		Short: "List ballots",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			ads := ballotapi.ListFilter(
				ctx,
				setup.Gov,
				ballotOnlyOpen,
				ballotOnlyClosed,
				ballotOnlyFrozen,
				member.User(ballotWithParticipant),
			)
			if ballotOnlyNames {
				for _, n := range ballotproto.AdsToBallotNames(ads) {
					fmt.Println(n)
				}
			} else {
				fmt.Fprint(os.Stdout, form.SprintJSON(ads))
			}
		},
	}

	ballotTallyCmd = &cobra.Command{
		Use:   "tally",
		Short: "Fetch current votes and record latest tally",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			chg := ballotapi.Tally(
				ctx,
				setup.Organizer,
				ballotproto.ParseBallotNameFromPath(ballotName),
				ballotFetchPar,
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
			chg := ballotapi.Vote(
				ctx,
				setup.Member,
				setup.Gov,
				ballotproto.ParseBallotNameFromPath(ballotName),
				parseElections(ctx, ballotElectionChoice, ballotElectionStrength),
			)
			fmt.Fprint(os.Stdout, form.SprintJSON(chg.Result))
		},
	}

	ballotTrackCmd = &cobra.Command{
		Use:   "track",
		Short: "Track the status of votes",
		Long: `Track reports on the status of the user's votes on a given ballot.
It returns a report listing accepted, rejected and pending votes.
Rejected votes are associated with a reason, such as "ballot is frozen".
Pending votes have not yet been processed by the community's governance.`,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			status := ballotapi.Track(
				ctx,
				setup.Member,
				setup.Gov,
				ballotproto.ParseBallotNameFromPath(ballotName),
			)
			fmt.Fprint(os.Stdout, form.SprintJSON(status))
		},
	}

	ballotEraseCmd = &cobra.Command{
		Use:   "erase",
		Short: "Erase a ballot",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			LoadConfig()
			chg := ballotapi.Erase(
				ctx,
				setup.Organizer,
				ballotproto.ParseBallotNameFromPath(ballotName),
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
	ballotOnlyNames        bool
	ballotOnlyOpen         bool
	ballotOnlyClosed       bool
	ballotOnlyFrozen       bool
	ballotWithParticipant  string
	ballotFetchPar         int
	ballotEscrowTo         string
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
	ballotCloseCmd.Flags().StringVar(&ballotEscrowTo, "to", account.BurnAccountID.String(), "account id to receive ballot escrows")

	// cancel
	ballotCmd.AddCommand(ballotCancelCmd)
	ballotCancelCmd.Flags().StringVar(&ballotName, "name", "", "ballot name")
	ballotCancelCmd.MarkFlagRequired("name")

	// freeze
	ballotCmd.AddCommand(ballotFreezeCmd)
	ballotFreezeCmd.Flags().StringVar(&ballotName, "name", "", "ballot name")
	ballotFreezeCmd.MarkFlagRequired("name")

	// unfreeze
	ballotCmd.AddCommand(ballotUnfreezeCmd)
	ballotUnfreezeCmd.Flags().StringVar(&ballotName, "name", "", "ballot name")
	ballotUnfreezeCmd.MarkFlagRequired("name")

	// show
	ballotCmd.AddCommand(ballotShowCmd)
	ballotShowCmd.Flags().StringVar(&ballotName, "name", "", "ballot name")
	ballotShowCmd.MarkFlagRequired("name")

	// list
	ballotCmd.AddCommand(ballotListCmd)
	ballotListCmd.Flags().BoolVar(&ballotOnlyNames, "only_names", false, "list only ballot names")
	ballotListCmd.Flags().BoolVar(&ballotOnlyOpen, "open", false, "list only open ballots")
	ballotListCmd.Flags().BoolVar(&ballotOnlyClosed, "closed", false, "list only closed ballots")
	ballotListCmd.Flags().BoolVar(&ballotOnlyFrozen, "frozen", false, "list only frozen ballots")
	ballotListCmd.Flags().StringVar(&ballotWithParticipant, "participant", "", "list only ballots where the given user is a participant")

	// tally
	ballotCmd.AddCommand(ballotTallyCmd)
	ballotTallyCmd.Flags().StringVar(&ballotName, "name", "", "ballot name")
	ballotTallyCmd.Flags().IntVar(&ballotFetchPar, "fetch_par", 5, "parallelism while clonging member repos for vote collection")
	ballotTallyCmd.MarkFlagRequired("name")

	// vote
	ballotCmd.AddCommand(ballotVoteCmd)
	ballotVoteCmd.Flags().StringVar(&ballotName, "name", "", "ballot name")
	ballotVoteCmd.MarkFlagRequired("name")
	ballotVoteCmd.Flags().StringSliceVar(&ballotElectionChoice, "choices", nil, "list of elected choices")
	ballotVoteCmd.MarkFlagRequired("choices")
	ballotVoteCmd.Flags().Float64SliceVar(&ballotElectionStrength, "strengths", nil, "list of elected vote strengths")
	ballotVoteCmd.MarkFlagRequired("strengths")

	// track
	ballotCmd.AddCommand(ballotTrackCmd)
	ballotTrackCmd.Flags().StringVar(&ballotName, "name", "", "ballot name")
	ballotTrackCmd.MarkFlagRequired("name")

	// erase
	ballotCmd.AddCommand(ballotEraseCmd)
	ballotEraseCmd.Flags().StringVar(&ballotName, "name", "", "ballot name")
	ballotEraseCmd.MarkFlagRequired("name")
}

func parseElections(ctx context.Context, choices []string, strengths []float64) ballotproto.Elections {
	if len(choices) != len(strengths) {
		must.Errorf(ctx, "elected choices must match elected strengths in count")
	}
	el := make(ballotproto.Elections, len(choices))
	for i := range choices {
		el[i] = ballotproto.NewElection(choices[i], strengths[i])
	}
	return el
}

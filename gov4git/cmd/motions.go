package cmd

import (
	"strings"

	"github.com/gov4git/gov4git/v2/gov4git/api"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/gov4git/v2/proto/motion"
	"github.com/gov4git/gov4git/v2/proto/motion/motionapi"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/spf13/cobra"
)

var (
	motionCmd = &cobra.Command{
		Use:   "motion",
		Short: "Manage motions (concerns and proposals)",
		Long:  ``,
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	motionOpenCmd = &cobra.Command{
		Use:   "open",
		Short: "Open a new motion",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			api.Invoke(
				func() {
					LoadConfig()
					motionapi.OpenMotion(
						ctx,
						setup.Organizer,
						motionproto.MotionID(motionName),
						motionproto.ParseMotionType(ctx, motionType),
						motion.PolicyName(motionPolicy),
						member.User(motionAuthor),
						motionTitle,
						motionDesc,
						motionTrackerURL,
						nil,
					)
				},
			)
		},
	}

	motionCloseCmd = &cobra.Command{
		Use:   "close",
		Short: "Close a motion",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			api.Invoke(
				func() {
					LoadConfig()
					var d motionproto.Decision
					if motionAccept {
						d = motionproto.Accept
					} else {
						d = motionproto.Reject
					}
					motionapi.CloseMotion(
						ctx,
						setup.Organizer,
						motionproto.MotionID(motionName),
						d,
					)
				},
			)
		},
	}

	motionListCmd = &cobra.Command{
		Use:   "list",
		Short: "List motions",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			api.Invoke1(
				func() any {
					LoadConfig()
					if motionTrack {
						return motionapi.TrackMotionBatch(ctx, setup.Gov, setup.Member)
					} else {
						return motionapi.ListMotionViews(ctx, setup.Gov)
					}
				},
			)
		},
	}

	motionShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show motion state and associated objects",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			api.Invoke1(
				func() any {
					LoadConfig()
					if motionTrack {
						return motionapi.TrackMotion(ctx, setup.Gov, setup.Member, motionproto.MotionID(motionName))
					} else {
						return motionapi.ShowMotion(ctx, setup.Gov, motionproto.MotionID(motionName))
					}
				},
			)
		},
	}
)

var (
	motionName       string
	motionPolicy     string
	motionAuthor     string
	motionTitle      string
	motionDesc       string
	motionType       string
	motionTrackerURL string
	motionAccept     bool
	motionTrack      bool
)

func init() {
	motionCmd.AddCommand(motionOpenCmd)
	motionOpenCmd.Flags().StringVar(&motionName, "name", "", "unique name for motion")
	motionOpenCmd.MarkFlagRequired("name")
	motionOpenCmd.Flags().StringVar(&motionPolicy, "policy", "", "policy ("+strings.Join(motionproto.InstalledMotionPolicies(), ", ")+")")
	motionOpenCmd.MarkFlagRequired("policy")
	motionOpenCmd.Flags().StringVar(&motionAuthor, "author", "", "author user name")
	motionOpenCmd.MarkFlagRequired("author")

	motionOpenCmd.Flags().StringVar(&motionTitle, "title", "", "title for motion")
	motionOpenCmd.Flags().StringVar(&motionDesc, "desc", "", "description for motion")
	motionOpenCmd.Flags().StringVar(&motionType, "type", "concern", "type of motion (concern, proposal)")
	motionOpenCmd.Flags().StringVar(&motionTrackerURL, "tracking", "", "tracking URL for motion")

	motionCmd.AddCommand(motionCloseCmd)
	motionCloseCmd.Flags().StringVar(&motionName, "name", "", "name of motion")
	motionCloseCmd.MarkFlagRequired("name")
	motionCloseCmd.Flags().BoolVar(&motionAccept, "accept", false, "accept/reject")

	motionCmd.AddCommand(motionListCmd)
	motionListCmd.Flags().BoolVar(&motionTrack, "track", false, "include this voter's tracking info with every motion")

	motionCmd.AddCommand(motionShowCmd)
	motionShowCmd.Flags().StringVar(&motionName, "name", "", "name of motion")
	motionShowCmd.MarkFlagRequired("name")
	motionShowCmd.Flags().BoolVar(&motionTrack, "track", false, "include this voter's tracking info")
}

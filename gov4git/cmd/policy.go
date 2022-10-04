package cmd

import (
	"fmt"
	"os"

	"github.com/petar/gov4git/lib/base"
	"github.com/petar/gov4git/lib/files"
	man "github.com/petar/gov4git/man/policy"
	"github.com/petar/gov4git/proto"
	"github.com/petar/gov4git/services/gov/policy"
	"github.com/spf13/cobra"
)

var (
	// policy management
	policyCmd = &cobra.Command{
		Use:   "policy",
		Short: "Manage directory governance policy",
		Long:  man.GovPolicy,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	policySetCmd = &cobra.Command{
		Use:   "set",
		Short: "Set directory governance",
		Long:  man.GovPolicySet,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := policy.GovPolicyService{
				GovConfig: proto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(proto.LocalAgentTempPath, "gov-policy-set")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.PolicySet(ctx, &policy.GovPolicySetIn{
				Dir:             policyDir,
				Arb:             policyArb,
				Group:           policyGroup,
				Threshold:       policyThreshold,
				CommunityBranch: communityBranch,
			})
			if err == nil {
				fmt.Fprint(os.Stdout, r.Human(cmd.Context()))
			} else {
				fmt.Fprint(os.Stderr, err.Error())
			}
			return err
		},
	}

	policyGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get policy property",
		Long:  man.GovPolicyGet,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := policy.GovPolicyService{
				GovConfig: proto.GovConfig{
					CommunityURL: communityURL,
				},
			}
			workDir, err := files.TempDir().MkEphemeralDir(proto.LocalAgentTempPath, "gov-policy-get")
			base.AssertNoErr(err)
			ctx := files.WithWorkDir(cmd.Context(), workDir)
			r, err := s.PolicyGet(ctx, &policy.GovPolicyGetIn{
				Dir:             policyDir,
				CommunityBranch: communityBranch,
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
	policyDir       string
	policyArb       string
	policyGroup     string
	policyThreshold float64
)

func init() {
	policyCmd.AddCommand(policySetCmd)
	policyCmd.AddCommand(policyGetCmd)

	policyCmd.PersistentFlags().StringVar(&policyDir, "dir", "", "governed directory")

	policySetCmd.Flags().StringVar(&policyArb, "arb", "", "arbitration policy (supported: quorum)")
	policySetCmd.Flags().StringVar(&policyGroup, "group", "", "arbitration group")
	policySetCmd.Flags().Float64Var(&policyThreshold, "thresh", 0.0, "voting or quorum threshold")
}

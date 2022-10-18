package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/base"
	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/cmdproto"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "gov4git",
		Short: "gov4git is a command-line client for transparent community operations",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

var (
	configPath string
	privateURL string
	publicURL  string
	verbose    bool
)

func init() {
	cobra.OnInitialize(initAfterFlags)

	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "config file (default is $HOME/.gov4git/config.json)")
	rootCmd.PersistentFlags().StringVar(&privateURL, "private_url", "", "private url of soul")
	rootCmd.PersistentFlags().StringVar(&publicURL, "public_url", "", "public url of soul")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "run in developer mode with verbose logging")
	// rootCmd.MarkPersistentFlagRequired("private_url")
	// rootCmd.MarkPersistentFlagRequired("public_url")
	// viper.BindPFlag("private_url", rootCmd.PersistentFlags().Lookup("private_url"))
	// viper.BindPFlag("public_url", rootCmd.PersistentFlags().Lookup("public_url"))

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(govCmd)
}

func initAfterFlags() {
	if verbose {
		base.LogVerbosely()
	} else {
		base.LogQuietly()
	}

	if configPath == "" {
		// find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			base.Fatalf("looking for home dir (%v)", err)
		}
		base.AssertNoErr(err)

		// search for config in ~/.gov4git/ directory with name "config" (without extension).
		configPath = filepath.Join(home, cmdproto.LocalAgentPath, "config.json")
	}

	if err := form.DecodeFormFromFile(context.Background(), configPath, &config); err == nil {
		if publicURL == "" {
			publicURL = config.PublicURL
		}
		if privateURL == "" {
			privateURL = config.PrivateURL
		}
		if communityURL == "" {
			communityURL = config.CommunityURL
		}
		if communityBranch == "" {
			communityBranch = config.CommunityBranch
		}
	}

	base.Infof("private_url=%v public_url=%v community_url=%v", privateURL, publicURL, communityURL)

	git.Init()
}

var config cmdproto.Config // used directly by some commands

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

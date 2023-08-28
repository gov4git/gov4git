package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	gov4git_root "github.com/gov4git/gov4git"
	"github.com/gov4git/gov4git/github"
	"github.com/gov4git/gov4git/gov4git/api"
	_ "github.com/gov4git/gov4git/runtime"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "gov4git",
		Short: "gov4git is a command-line client for the gov4git community governance protocol",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

var ctx = github.WithTokenSource(git.WithTTL(git.WithAuth(context.Background(), nil), nil), nil)

var (
	configPath string
	verbose    bool
)

func init() {
	cobra.OnInitialize(initAfterFlags)

	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "config file (default is $HOME/.gov4git/config.json)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "run in developer mode with verbose logging")

	rootCmd.AddCommand(initIDCmd)
	rootCmd.AddCommand(initGovCmd)
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(groupCmd)
	rootCmd.AddCommand(memberCmd)
	rootCmd.AddCommand(ballotCmd)
	rootCmd.AddCommand(balanceCmd)
	rootCmd.AddCommand(bureauCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(cronCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(cacheCmd)
	rootCmd.AddCommand(githubCmd)
	rootCmd.AddCommand(collabCmd)
}

func initAfterFlags() {
	if verbose {
		base.LogVerbosely()
	} else {
		base.LogQuietly()
	}
	base.Infof("gov4git version: %v", gov4git_root.Short())
}

func LoadConfig() {
	if configPath == "" {
		// find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			base.Fatalf("looking for home dir (%v)", err)
		}
		base.AssertNoErr(err)

		// search for config in ~/.gov4git/config.json
		configPath = filepath.Join(home, api.LocalAgentPath, "config.json")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		base.Fatalf("reading config file (%v)", err)
	}

	config, err := form.DecodeBytes[api.Config](ctx, data)
	if err != nil {
		base.Fatalf("decoding config file (%v)", err)
	}

	if config.CacheDir != "" {
		ctx = git.WithCache(ctx, config.CacheDir)
	}

	setup = config.Setup(ctx)
}

var (
	setup api.Setup
)

func Execute() int {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return 0
}

func ExecuteWithConfig(cfgPath string) int {
	configPath = cfgPath
	return Execute()
}

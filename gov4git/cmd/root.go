package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/base"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/cmdproto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "config file (default is $HOME/.gov/config.json)")
	rootCmd.PersistentFlags().StringVar(&privateURL, "private_url", "", "private url of soul")
	rootCmd.PersistentFlags().StringVar(&publicURL, "public_url", "", "public url of soul")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "run in developer mode with verbose logging")
	// rootCmd.MarkPersistentFlagRequired("private_url")
	// rootCmd.MarkPersistentFlagRequired("public_url")
	viper.BindPFlag("private_url", rootCmd.PersistentFlags().Lookup("private_url"))
	viper.BindPFlag("public_url", rootCmd.PersistentFlags().Lookup("public_url"))

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(govCmd)
}

func initAfterFlags() {
	if verbose {
		base.LogVerbosely()
	} else {
		base.LogQuietly()
	}

	if configPath != "" {
		viper.SetConfigFile(configPath) // use config file from the flag
	} else {
		// find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			base.Fatalf("looking for home dir (%v)", err)
		}
		base.AssertNoErr(err)

		// search for config in ~/.gov/ directory with name "config" (without extension).
		viper.AddConfigPath(filepath.Join(home, cmdproto.LocalAgentPath))
		viper.SetConfigType("json")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		base.Infof("using config file %v", viper.ConfigFileUsed())
	}

	git.Init()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

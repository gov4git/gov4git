package cmd

import (
	"os"
	"path/filepath"

	"github.com/petar/gitty/lib/base"
	"github.com/petar/gitty/proto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "ana",
		Short: "ana is a command-line tool for ...",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	// govCmd = &cobra.Command{
	// 	Use:   "gov",
	// 	Short: "",
	// 	Long:  ``,
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 	},
	// }
)

var configPath string
var privateURL string
var publicURL string

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "config file (default is $HOME/.ana/config.json)")
	rootCmd.PersistentFlags().StringVar(&privateURL, "private_url", "", "private url of soul")
	rootCmd.PersistentFlags().StringVar(&publicURL, "public_url", "", "public url of soul")
	// rootCmd.MarkPersistentFlagRequired("private_url")
	// rootCmd.MarkPersistentFlagRequired("public_url")
	viper.BindPFlag("private_url", rootCmd.PersistentFlags().Lookup("private_url"))
	viper.BindPFlag("public_url", rootCmd.PersistentFlags().Lookup("public_url"))

	rootCmd.AddCommand(initCmd)

	// rootCmd.AddCommand(sendCmd)
	// rootCmd.AddCommand(receiveCmd)
	// rootCmd.AddCommand(govCmd)
}

func initConfig() {
	if configPath != "" {
		viper.SetConfigFile(configPath) // use config file from the flag
	} else {
		// find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			base.Fatalf("looking for home dir (%v)", err)
		}
		base.AssertNoErr(err)

		// search for config in ~/.ana/ directory with name "config" (without extension).
		viper.AddConfigPath(filepath.Join(home, proto.LocalAgentPath))
		viper.SetConfigType("json")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		base.Infof("using config file %v", viper.ConfigFileUsed())
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		base.Fatalf("command error (%v)", err)
		base.Sync()
	}
}

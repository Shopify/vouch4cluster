package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vouch4cluster",
	Short: "vouch4cluster lets you run voucher against all of the images running in a cluster or deployment",
	Long:  `vouch4cluster lets you run voucher against all of the images running in a cluster or deployment`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vouch4cluster.yaml)")
	rootCmd.PersistentFlags().StringVar(&defaultConfig.Hostname, "voucher", "v", "Voucher server to connect to.")
	viper.BindPFlag("voucher.hostname", rootCmd.PersistentFlags().Lookup("voucher"))
	rootCmd.PersistentFlags().StringVar(&defaultConfig.Username, "username", "", "Username to authenticate against Voucher with")
	viper.BindPFlag("voucher.username", rootCmd.PersistentFlags().Lookup("username"))
	rootCmd.PersistentFlags().StringVar(&defaultConfig.Password, "password", "", "Password to authenticate against Voucher with")
	viper.BindPFlag("voucher.password", rootCmd.PersistentFlags().Lookup("password"))
	rootCmd.PersistentFlags().IntVar(&defaultConfig.Workers, "workers", 100, "The number of workers to spawn to call Voucher with.")
	viper.BindPFlag("voucher.workers", rootCmd.PersistentFlags().Lookup("workers"))
	rootCmd.PersistentFlags().StringSliceVar(&defaultConfig.Checks, "checks", []string{"all"}, "The Voucher checks to require an image to pass.")
	viper.BindPFlag("voucher.checks", rootCmd.PersistentFlags().Lookup("checks"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".vouch4cluster" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".vouch4cluster")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
		if err = viper.UnmarshalKey("voucher", &defaultConfig); nil != err {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func errorf(format string, v interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", v)
}

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/bgreen/space-traders-go/st"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "st",
	Short: "A command line utility for the Space Traders API",
	Long:  `A command line utility for the Space Traders API.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	client = st.NewClient()
	client.Start()
	defer client.Stop()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var client *st.Client

func init() {

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.space-traders-go.yaml)")
	rootCmd.PersistentFlags().BoolP("long", "l", false, "Longer output")
	rootCmd.PersistentFlags().BoolP("raw", "r", false, "Print raw messages")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

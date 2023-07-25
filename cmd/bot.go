/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/bgreen/space-traders-go/bot"
	"github.com/spf13/cobra"
)

// botCmd represents the bot command
var botCmd = &cobra.Command{
	Use:   "bot",
	Short: "Automated tasks",
	Long:  `Automated tasks`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("bot called")
	},
}

var minerCmd = &cobra.Command{
	Use:   "miner [ship]",
	Short: "An auto-mining bot",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 1 {
			ship, err := client.GetMyShip(args[0])
			if err != nil {
				return
			}
			bot.MinerBot{
				Client: client,
				Ship:   ship,
			}.RunOnce()
		} else if len(args) == 0 {
			m := bot.NewManager(client)
			m.Run()
		}
	},
}

func init() {
	rootCmd.AddCommand(botCmd)

	botCmd.AddCommand(minerCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// botCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// botCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/bgreen/space-traders-go/stapi"
	"github.com/spf13/cobra"
)

// agentCmd represents the agent command
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Print agent info",
	Long:  `Print agent info.`,
	Run: func(cmd *cobra.Command, args []string) {
		a, err := client.GetMyAgent()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Print(agentInfoLong(a))
	},
}

func agentInfoLong(a stapi.Agent) string {
	var s string

	s += fmt.Sprintf("Name:    %v\n", a.Symbol)
	s += fmt.Sprintf("Credits: %v\n", a.Credits)
	s += fmt.Sprintf("HQ:      %v\n", a.Headquarters)
	s += fmt.Sprintf("Faction: %v\n", a.StartingFaction)
	s += fmt.Sprintf("Account: %v\n", a.AccountId)

	return s
}

func init() {
	rootCmd.AddCommand(agentCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// agentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// agentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

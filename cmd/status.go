/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/bgreen/space-traders-go/stapi"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get server status",
	Long:  `Get server status.`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := client.GetStatus()
		if err != nil {
			return
		}

		fmt.Print(statusInfoLong(resp))
	},
}

func statusInfoLong(status stapi.GetStatus200Response) string {
	var s string

	s += fmt.Sprintf("%v\n", status.Status)
	s += fmt.Sprintf("Version: %v\n", status.Version)
	s += fmt.Sprintf("Last Reset: %v\n", status.ResetDate)
	s += fmt.Sprintf("Next Reset: %v\n", status.ServerResets.Next)

	s += fmt.Sprintln()
	s += fmt.Sprintln("Stats")
	s += fmt.Sprintf("Agents:    %v\n", status.Stats.Agents)
	s += fmt.Sprintf("Systems:   %v\n", status.Stats.Systems)
	s += fmt.Sprintf("Ships:     %v\n", status.Stats.Ships)
	s += fmt.Sprintf("Waypoints: %v\n", status.Stats.Waypoints)

	s += fmt.Sprintln()
	for _, v := range status.Announcements {
		s += fmt.Sprintf("========== %v ==========\n%v\n", v.Title, v.Body)
	}

	s += fmt.Sprintln()
	s += fmt.Sprintln("Leaderboards (Credits)")
	for i, v := range status.Leaderboards.MostCredits {
		s += fmt.Sprintf("%2v: %-20v: %v\n", i+1, v.AgentSymbol, v.Credits)
	}

	s += fmt.Sprintln()
	s += fmt.Sprintln("Leaderboards (Charts)")
	for i, v := range status.Leaderboards.MostSubmittedCharts {
		s += fmt.Sprintf("%2v: %-20v: %v\n", i+1, v.AgentSymbol, v.ChartCount)
	}

	return s
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

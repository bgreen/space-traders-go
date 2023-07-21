/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/bgreen/space-traders-go/stapi"
	"github.com/spf13/cobra"
)

// systemCmd represents the system command
var systemCmd = &cobra.Command{
	Use:   "system [symbol]",
	Short: "Print system info",
	Long:  `Print system info`,
	Run: func(cmd *cobra.Command, args []string) {
		var systems []stapi.System
		c, _ := cmd.Flags().GetInt("count")

		if len(args) == 0 {
			resp, err := client.GetMoreSystems(0, c)
			if err != nil {
				fmt.Println(err)
				return
			}
			systems = append(systems, resp...)

		} else if len(args) == 1 {
			resp, err := client.GetSystem(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
			systems = append(systems, resp)
		}

		for i, s := range systems {
			if i == c {
				return
			}
			l, _ := cmd.Flags().GetBool("long")
			if l {
				wps, _ := client.GetSystemWaypoints(s.Symbol)
				fmt.Print(systemInfoLong(s, wps) + "\n")
			} else {
				fmt.Print(systemInfoShort(s))
			}
		}
	},
}

func systemInfoShort(sys stapi.System) string {
	return fmt.Sprintf("%-7v:\t%v\tX:%6v Y:%6v\n", sys.Symbol, sys.Type, sys.X, sys.Y)
}

func systemInfoLong(sys stapi.System, wps []stapi.Waypoint) string {
	var s string
	s += fmt.Sprintf("%-7v:\t%v\tX:%6v Y:%6v\n", sys.Symbol, sys.Type, sys.X, sys.Y)
	s += fmt.Sprintln("Waypoints:")
	for _, v := range wps {
		s += fmt.Sprintf("%-14v:\t%-14v\tX:%3v Y:%3v\t", v.Symbol, v.Type, v.X, v.Y)
		for _, t := range v.Traits {
			s += fmt.Sprintf("%v, ", t.Name)
		}
		s += "\n"
	}

	return s
}

func init() {
	rootCmd.AddCommand(systemCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// systemCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// systemCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	systemCmd.Flags().BoolP("long", "l", false, "print more info")
	systemCmd.Flags().IntP("count", "c", 10, "Maximum number of systems to list")
}

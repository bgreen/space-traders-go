/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new Space Traders agent",
	Long: `Register a new Space Traders agent

st register <agent name> [faction symbol]`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		faction := "COSMIC"
		if len(args) >= 2 {
			faction = args[1]
		}

		fmt.Printf("Register name \"%v\" faction \"%v\"\n", name, faction)
		r, err := client.Register(name, faction)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(r)
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// registerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

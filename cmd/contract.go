/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// contractCmd represents the contract command
var contractCmd = &cobra.Command{
	Use:   "contract",
	Short: "Describe and interact with contracts",
	Long:  `Describe and interact with contracts.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("contract called")
	},
}

func init() {
	rootCmd.AddCommand(contractCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// contractCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// contractCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

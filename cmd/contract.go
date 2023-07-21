/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/bgreen/space-traders-go/stapi"
	"github.com/spf13/cobra"
)

// contractCmd represents the contract command
var contractCmd = &cobra.Command{
	Use:   "contract [symbol]",
	Short: "Describe and interact with contracts",
	Long:  `Describe and interact with contracts.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var contracts []stapi.Contract
		c, _ := cmd.Flags().GetInt("count")

		if len(args) == 0 {
			resp, err := client.GetContracts()
			if err != nil {
				fmt.Println(err)
				return
			}
			contracts = append(contracts, resp...)

		} else if len(args) == 1 {
			resp, err := client.GetContract(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
			contracts = append(contracts, resp)
		}

		for i, s := range contracts {
			if i == c {
				return
			}
			l, _ := cmd.Flags().GetBool("long")
			if l {
				fmt.Print(contractInfoLong(s) + "\n")
			} else {
				fmt.Print(contractInfoShort(s))
			}
		}
	},
}

func contractInfoShort(c stapi.Contract) string {
	var s string
	s += fmt.Sprintf("%v: %v %v\n", c.Id, c.Type, c.Accepted)
	return s
}

func contractInfoLong(c stapi.Contract) string {
	var s string
	s += fmt.Sprintf("%v:\t%v\t%v\t%v\n", c.Id, c.Type, c.Terms.Deadline, c.Accepted)
	for _, v := range c.Terms.Deliver {
		s += fmt.Sprintf("%v:\t%v\t%v/%v", v.TradeSymbol, v.DestinationSymbol, v.UnitsFulfilled, v.UnitsRequired)
	}
	return s
}

var acceptContractCommand = &cobra.Command{
	Use:   "accept <id>",
	Short: "Accept a contract",
	Long:  "Accept a contract",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := client.AcceptContract(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%v\n", resp)
	},
}

var deliverContractCommand = &cobra.Command{
	Use:   "deliver <id> <ship> <trade symbol> <count>",
	Short: "Deliver a contract",
	Long:  "Deliver a contract",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		count, _ := strconv.Atoi(args[3])
		resp, err := client.DeliverContract(args[0], args[1], stapi.TradeSymbol(args[2]), int32(count))
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%v\n", resp)
	},
}

var fulfillContractCommand = &cobra.Command{
	Use:   "fulfill <id>",
	Short: "Fulfill a contract",
	Long:  "Fulfill a contract",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := client.FulfillContract(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%v\n", resp)
	},
}

var negotiateContractCommand = &cobra.Command{
	Use:   "negotiate <ship>",
	Short: "Create a new contract",
	Long:  "Create a new contract",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := client.NegotiateContract(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%v\n", resp)
	},
}

func init() {
	rootCmd.AddCommand(contractCmd)

	contractCmd.AddCommand(acceptContractCommand)
	contractCmd.AddCommand(deliverContractCommand)
	contractCmd.AddCommand(fulfillContractCommand)
	contractCmd.AddCommand(negotiateContractCommand)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// contractCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// contractCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	contractCmd.Flags().BoolP("long", "l", false, "print more info")
	contractCmd.Flags().IntP("count", "c", 10, "Maximum number of systems to list")
}

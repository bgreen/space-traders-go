/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"

	"github.com/bgreen/space-traders-go/stapi"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// shipCmd represents the ship command
var shipCmd = &cobra.Command{
	Use:   "ship [symbol]",
	Short: "Describe and interact with ships",
	Long:  `Describe and interact with ships.`,
	Run: func(cmd *cobra.Command, args []string) {
		var ships []stapi.Ship

		if len(args) == 0 {
			v, err := client.GetMyShips()
			if err != nil {
				fmt.Println(err)
				return
			}

			ships = append(ships, v...)
		} else if len(args) == 1 {
			v, err := client.GetMyShip(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}

			ships = append(ships, v)
		}

		for _, s := range ships {
			l, _ := cmd.Flags().GetBool("long")
			if l {
				fmt.Print(shipInfoLong(s))
				fmt.Print("\n")
			} else {
				fmt.Print(shipInfoShort(s))
			}
		}
	},
}

func shipInfoShort(ship stapi.Ship) string {
	return fmt.Sprintf("%v:\t%v\t%v\n", ship.Symbol, ship.Frame.Name, ship.Nav.WaypointSymbol)
}

func shipInfoLong(ship stapi.Ship) string {
	var s string
	s += fmt.Sprintf("%v:\t%v\t%v\tFuel: %v/%v\n",
		ship.Symbol,
		ship.Frame.Name,
		ship.Nav.WaypointSymbol,
		ship.Fuel.Current, ship.Fuel.Capacity)
	var modules string
	for i, v := range ship.Modules {
		modules += fmt.Sprintf("Module %02v: %v\n", i, v.Name)
	}
	var mounts string
	for i, v := range ship.Mounts {
		mounts += fmt.Sprintf("Mount %02v: %v\n", i, v.Name)
	}
	s += lg.JoinHorizontal(lg.Left, modules, mounts)
	return s
}

var shipOrbitCmd = &cobra.Command{
	Use:   "orbit <symbol>",
	Short: "Control a ship",
	Long:  `Control a ship`,
	Run: func(cmd *cobra.Command, args []string) {
		v, err := client.OrbitShip(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%v\n", v)
	},
}

var shipDockCmd = &cobra.Command{
	Use:   "dock <symbol>",
	Short: "Control a ship",
	Long:  `Control a ship`,
	Run: func(cmd *cobra.Command, args []string) {
		v, err := client.DockShip(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%v\n", v)
	},
}

var shipNavCmd = &cobra.Command{
	Use:   "nav <ship> <waypoint>",
	Short: "Pilot a ship to a destination",
	Long:  `Pilot a ship to a destination`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		v, err := client.NavigateShip(args[0], args[1])
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%v\n", v)
	},
}

var shipJumpCmd = &cobra.Command{
	Use:   "jump <ship> <system>",
	Short: "Jump a ship to a destination",
	Long:  `Jump a ship to a destination`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		v, err := client.JumpShip(args[0], args[1])
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%v\n", v)
	},
}

var shipWarpCmd = &cobra.Command{
	Use:   "warp <ship> <waypoint>",
	Short: "Warp a ship to a destination",
	Long:  `Warp a ship to a destination`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		v, err := client.WarpShip(args[0], args[1])
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%v\n", v)
	},
}

var shipRefuelCmd = &cobra.Command{
	Use:   "refuel <ship>",
	Short: "Refuel a ship",
	Long:  `Refuel a ship`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := client.RefuelShip(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%v\n", resp)
	},
}

var shipMineCmd = &cobra.Command{
	Use:   "mine <ship>",
	Short: "Mine resources with a ship",
	Long:  `Mine resources with a ship`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := client.ExtractResources(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%v\n", resp)
	},
}

var shipCargoCmd = &cobra.Command{
	Use:   "cargo <ship>",
	Short: "List ship cargo",
	Long:  `List ship cargo`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		v, err := client.GetMyShipCargo(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Print(shipCargoInfoLong(v))

	},
}

func shipCargoInfoLong(c stapi.ShipCargo) string {
	var s string
	s += fmt.Sprintf("Capacity: %v/%v\n", c.Units, c.Capacity)
	for _, v := range c.Inventory {
		s += fmt.Sprintf("%-20v Units:%-5v\n", v.Symbol, v.Units)
	}
	return s
}

var shipCargoSellCmd = &cobra.Command{
	Use:   "sell <ship> <trade symbol> <count>",
	Short: "Sell resources with a ship",
	Long:  `Sell resources with a ship`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		count, _ := strconv.Atoi(args[2])
		resp, err := client.SellCargo(args[0], stapi.TradeSymbol(args[1]), int32(count))
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%v\n", resp)
	},
}

var shipCargoBuyCmd = &cobra.Command{
	Use:   "buy <ship> <trade symbol> <count>",
	Short: "Buy resources with a ship",
	Long:  `Buy resources with a ship`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		count, _ := strconv.Atoi(args[2])
		resp, err := client.PurchaseCargo(args[0], stapi.TradeSymbol(args[1]), int32(count))
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%v\n", resp)
	},
}

func init() {
	rootCmd.AddCommand(shipCmd)

	shipCmd.AddCommand(shipOrbitCmd)
	shipCmd.AddCommand(shipDockCmd)
	shipCmd.AddCommand(shipNavCmd)
	shipCmd.AddCommand(shipJumpCmd)
	shipCmd.AddCommand(shipWarpCmd)
	shipCmd.AddCommand(shipRefuelCmd)
	shipCmd.AddCommand(shipMineCmd)
	shipCmd.AddCommand(shipCargoCmd)

	shipCargoCmd.AddCommand(shipCargoSellCmd)
	shipCargoCmd.AddCommand(shipCargoBuyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// shipCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// shipCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	shipCmd.Flags().BoolP("long", "l", false, "print more info")
}

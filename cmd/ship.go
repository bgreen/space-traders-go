/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"
	"time"

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
				return
			}

			ships = append(ships, v...)
		} else if len(args) == 1 {
			v, err := client.GetMyShip(args[0])
			if err != nil {
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
	s += fmt.Sprintf("%v: %v (%v)\n", ship.Symbol, ship.Frame.Name[6:], ship.Registration.Role)

	s += fmt.Sprintf("Location: %v\n", ship.Nav.WaypointSymbol)
	s += fmt.Sprintf("Engine:   %v (%v/%v)\n", ship.Engine.Name, ship.Fuel.Current, ship.Fuel.Capacity)
	s += fmt.Sprintf("Reactor:  %v\n", ship.Reactor.Name)

	s += fmt.Sprintf("Status:   %v %v\n", ship.Nav.Status, ship.Nav.FlightMode)
	if ship.Nav.Status == "IN_TRANSIT" {
		s += fmt.Sprintf("Departure: %14v %v\n", ship.Nav.Route.Departure.Symbol, ship.Nav.Route.DepartureTime.Local())
		s += fmt.Sprintf("Arrival:   %14v %v (%v)\n",
			ship.Nav.Route.Destination.Symbol,
			ship.Nav.Route.Arrival.Local(),
			time.Until(ship.Nav.Route.Arrival))
	}

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
		p, _ := cmd.Flags().GetString("patch")
		if p != "" {
			r, err := client.PatchShipNav(args[0], p)
			if err != nil {
				return
			}

			fmt.Println(r)
		}

		v, err := client.NavigateShip(args[0], args[1])
		if err != nil {
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
			return
		}

		fmt.Printf("%v\n", resp)
	},
}

var shipRefineCmd = &cobra.Command{
	Use:   "refine <ship> <product>",
	Short: "Refine resources with a ship",
	Long: `Refine resources with a ship.

Allowable products are:
IRON COPPER SILVER GOLD ALUMINUM PLATINUM URANITE MERITIUM FUEL`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := client.ShipRefine(args[0], args[1])
		if err != nil {
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
			return
		}

		fmt.Printf("%v\n", resp)
	},
}

var shipSurveyCmd = &cobra.Command{
	Use:   "survey <ship>",
	Short: "Perform surveys with a surveyor",
	Long:  `Perform surveys with a surveyor.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		v, err := client.CreateSurvey(args[0])
		if err != nil {
			return
		}

		fmt.Println(v)

	},
}

var shipScanCmd = &cobra.Command{
	Use:   "scan <ship>",
	Short: "Perform scans with a sensor array",
	Long:  `Perform scans with a sensor array.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var s any
		var err error
		if f, _ := cmd.Flags().GetBool("ship"); f {
			s, err = client.CreateShipShipScan(args[0])
		} else if f, _ := cmd.Flags().GetBool("waypoint"); f {
			s, err = client.CreateShipWaypointScan(args[0])
		} else if f, _ := cmd.Flags().GetBool("system"); f {
			s, err = client.CreateShipSystemScan(args[0])
		}

		if err != nil {
			return
		}

		switch sp := s.(type) {
		case stapi.CreateShipShipScan201ResponseData:
			fmt.Println(shipShipScanInfoLong(sp))
		case stapi.CreateShipWaypointScan201ResponseData:
			fmt.Println(shipWaypointScanInfoLong(sp))
		case stapi.CreateShipSystemScan201ResponseData:
			fmt.Println(shipSystemScanInfoLong(sp))
		default:
			return
		}
	},
}

func shipShipScanInfoLong(scan stapi.CreateShipShipScan201ResponseData) string {
	var s string
	for _, v := range scan.Ships {
		s += fmt.Sprintf("%-24v: %v\n", v.Frame.Symbol, v.Registration.Name)
	}
	return s
}

func shipWaypointScanInfoLong(scan stapi.CreateShipWaypointScan201ResponseData) string {
	var s string
	for _, v := range scan.Waypoints {
		s += fmt.Sprintf("%-14v: %v\n", v.Symbol, v.Type)
	}
	return s
}

func shipSystemScanInfoLong(scan stapi.CreateShipSystemScan201ResponseData) string {
	var s string
	for _, v := range scan.Systems {
		s += fmt.Sprintf("%-7v: %-12v Dist:%5v\n", v.Symbol, v.Type, v.Distance)
	}
	return s
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
	shipCmd.AddCommand(shipRefineCmd)
	shipCmd.AddCommand(shipCargoCmd)
	shipCmd.AddCommand(shipSurveyCmd)
	shipCmd.AddCommand(shipScanCmd)

	shipCargoCmd.AddCommand(shipCargoSellCmd)
	shipCargoCmd.AddCommand(shipCargoBuyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// shipCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// shipCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	shipNavCmd.Flags().StringP("patch", "p", "", "Set flight mode")

	shipScanCmd.Flags().Bool("system", false, "perform a system scan")
	shipScanCmd.Flags().Bool("waypoint", false, "perform a waypoint scan")
	shipScanCmd.Flags().Bool("ship", false, "perform a ship scan")
	shipScanCmd.MarkFlagsMutuallyExclusive("system", "waypoint", "ship")
}

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

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

func init() {
	rootCmd.AddCommand(shipCmd)

	shipCmd.AddCommand(shipOrbitCmd)
	shipCmd.AddCommand(shipDockCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// shipCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// shipCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	shipCmd.Flags().BoolP("long", "l", false, "print more info")
}

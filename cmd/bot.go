/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"math"
	"time"

	"github.com/bgreen/space-traders-go/stapi"
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
	Use:   "miner <ship>",
	Short: "An auto-mining bot",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for {
			mineUntilFull(args[0])
			sellAllCargo(args[0])
		}
	},
}

func mineUntilFull(shipSymbol string) {
	ship, err := client.GetMyShip(shipSymbol)
	if err != nil {
		fmt.Println(err)
		return
	}

	waypoints, err := client.GetSystemWaypoints(ship.Nav.SystemSymbol)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get the full waypoint info
	var wp stapi.Waypoint
	for _, v := range waypoints {
		if v.Symbol == ship.Nav.WaypointSymbol {
			wp = v
			break
		}
	}

	fmt.Printf("Current location: %v (%v, %v)\n", wp.Symbol, wp.X, wp.Y)

	closest := findClosestWaypoint(wp, waypoints, isWaypointMine)

	fmt.Printf("Closest mine: %v (%v, %v) at %v units away\n", closest.Symbol, closest.X, closest.Y, waypointDistance(wp, closest))

	navAndRefuel(ship, closest)

	// Check if there's a cooldown active
	cd, _ := client.GetShipCooldown(ship.Symbol)
	if cd.RemainingSeconds > 0 {
		t := time.Until(*cd.Expiration)
		fmt.Printf("Waiting %v for cooldown\n", t)
		time.Sleep(t)
		fmt.Printf("Cooldown complete\n")
	}

	for {
		r0, err := client.ExtractResources(shipSymbol)
		if err != nil {
			return
		}

		fmt.Printf("Mined %v units of %v\n", r0.Extraction.Yield.Units, r0.Extraction.Yield.Symbol)

		fmt.Printf("Cargo at %v/%v\n", r0.Cargo.Units, r0.Cargo.Capacity)

		if r0.Cargo.Units == r0.Cargo.Capacity {
			break
		}

		t := time.Until(*r0.Cooldown.Expiration)
		fmt.Printf("Waiting %v for cooldown\n", t)
		time.Sleep(t)
		fmt.Printf("Cooldown complete\n")
	}
	fmt.Println("Done mining")
}

func isWaypointMine(w stapi.Waypoint) bool {
	for _, v := range w.Traits {
		if (v.Symbol == "MINERAL_DEPOSITS") ||
			(v.Symbol == "COMMON_METAL_DEPOSITS") ||
			(v.Symbol == "PRECIOUS_METAL_DEPOSITS") ||
			(v.Symbol == "RARE_METAL_DEPOSITS") {
			return true
		}
	}
	return false
}

func waypointDistance(w1, w2 stapi.Waypoint) float64 {
	a := math.Abs(float64(w1.X - w2.X))
	b := math.Abs(float64(w1.Y - w2.Y))
	c := math.Sqrt(a*a + b*b)
	return c
}

func findClosestWaypoint(current stapi.Waypoint, waypoints []stapi.Waypoint, filter func(stapi.Waypoint) bool) stapi.Waypoint {

	var closest stapi.Waypoint
	closestDistance := 10000.0
	for _, v := range waypoints {
		if filter(v) {
			dist := waypointDistance(current, v)
			if dist < closestDistance {
				closestDistance = dist
				closest = v
			}
		}
	}
	return closest
}

func navAndRefuel(ship stapi.Ship, dest stapi.Waypoint) {
	// If ship is docked, undock
	if ship.Nav.Status == "DOCKED" {
		r0, _ := client.OrbitShip(ship.Symbol)
		fmt.Printf("Ship %v\n", r0.Nav.Status)
	}

	// If ship isn't already there, go there
	if ship.Nav.WaypointSymbol != dest.Symbol {

		// Set thrusters to slow
		if isLowFuel(ship.Fuel) && (ship.Nav.FlightMode != "DRIFT") {
			r1, _ := client.PatchShipNav(ship.Symbol, "DRIFT")
			fmt.Printf("Set flight mode to %v due to low fuel (%v/%v)\n", r1.FlightMode, ship.Fuel.Current, ship.Fuel.Capacity)
		} else if ship.Nav.FlightMode != "CRUISE" {
			r1, _ := client.PatchShipNav(ship.Symbol, "CRUISE")
			fmt.Printf("Set flight mode to %v\n", r1.FlightMode)
		}

		// Go to the destination
		r2, err := client.NavigateShip(ship.Symbol, dest.Symbol)
		if err != nil {
			return
		}
		fmt.Printf("Navigating to %v\n", r2.Nav.WaypointSymbol)

		// Wait to get there
		t := time.Until(r2.Nav.Route.Arrival)
		fmt.Printf("Waiting %v to arrive\n", t)
		time.Sleep(t)
		fmt.Println("Arrived")

		// Refuel
		if isWaypointMarket(dest) && isLowFuel(r2.Fuel) {
			fmt.Printf("Refueling due to low fuel (%v/%v)\n", r2.Fuel.Current, r2.Fuel.Capacity)
			client.DockShip(ship.Symbol)
			client.RefuelShip(ship.Symbol)
			client.OrbitShip(ship.Symbol)
		}
	}
}

func isLowFuel(fuel stapi.ShipFuel) bool {
	return (fuel.Current == 0) || (float64(fuel.Current)/float64(fuel.Capacity) <= 0.5)
}

func sellAllCargo(shipSymbol string) {
	ship, err := client.GetMyShip(shipSymbol)
	if err != nil {
		return
	}

	waypoints, err := client.GetSystemWaypoints(ship.Nav.SystemSymbol)
	if err != nil {
		return
	}

	// Get the full waypoint info
	var wp stapi.Waypoint
	for _, v := range waypoints {
		if v.Symbol == ship.Nav.WaypointSymbol {
			wp = v
			break
		}
	}

	fmt.Printf("Current location: %v (%v, %v)\n", wp.Symbol, wp.X, wp.Y)

	closest := findClosestWaypoint(wp, waypoints, isWaypointMarket)

	fmt.Printf("Closest market: %v (%v, %v) at %v units away\n", closest.Symbol, closest.X, closest.Y, waypointDistance(wp, closest))

	navAndRefuel(ship, closest)

	// Try to dock, ignore error if already docked
	client.DockShip(ship.Symbol)

	market, _ := client.GetMarket(closest.SystemSymbol, closest.Symbol)

	for _, v := range ship.Cargo.Inventory {
		if isMarketBuying(market, v.Symbol) {
			r0, err := client.SellCargo(ship.Symbol, stapi.TradeSymbol(v.Symbol), v.Units)
			if err != nil {
				fmt.Printf("Couldn't sell %v units of %v\n", v.Units, v.Symbol)
			} else {
				t := r0.Transaction
				fmt.Printf("Sold %v units of %v at %v for a total of %v\n", t.Units, t.TradeSymbol, t.PricePerUnit, t.TotalPrice)
			}
		} else {
			fmt.Printf("Market does not buy %v\n", v.Symbol)
		}
	}
	fmt.Println("Done selling")
}

func isMarketBuying(market stapi.Market, trade string) bool {
	for _, v := range market.Exports {
		if string(v.Symbol) == trade {
			return true
		}
	}
	for _, v := range market.Exchange {
		if string(v.Symbol) == trade {
			return true
		}
	}
	return false
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

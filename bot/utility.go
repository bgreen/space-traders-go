package bot

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/bgreen/space-traders-go/st"
	"github.com/bgreen/space-traders-go/stapi"
)

func isSectorSymbol(s string) bool {
	return strings.Count(s, "-") == 0
}

func isSystemSymbol(s string) bool {
	return strings.Count(s, "-") == 1
}

func isWaypointSymbol(s string) bool {
	return strings.Count(s, "-") == 2
}

func isWaypointMarket(w stapi.Waypoint) bool {
	for _, v := range w.Traits {
		if v.Symbol == "MARKETPLACE" {
			return true
		}
	}
	return false
}

func isWaypointShipyard(w stapi.Waypoint) bool {
	for _, v := range w.Traits {
		if v.Symbol == "SHIPYARD" {
			return true
		}
	}
	return false
}

func isWaypointJumpGate(w stapi.Waypoint) bool {
	return w.Type == stapi.WaypointType("JUMP_GATE")
}

func waypointSymbolToSystemSymbol(w string) string {
	ss := strings.Split(w, "-")
	return strings.Join(ss[:2], "-")
}

func isShipMiner(s stapi.Ship) bool {
	for _, v := range s.Mounts {
		if strings.Contains(v.Symbol, "MINING_LASER") {
			return true
		}
	}
	return false
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

func routeDistance(r stapi.ShipNavRoute) float64 {
	a := math.Abs(float64(r.Departure.X - r.Destination.X))
	b := math.Abs(float64(r.Departure.Y - r.Destination.Y))
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

func isLowFuel(fuel stapi.ShipFuel) bool {
	return (fuel.Current == 0) || (float64(fuel.Current)/float64(fuel.Capacity) <= 0.5)
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

func navAndRefuel(client *st.Client, ship stapi.Ship, dest stapi.Waypoint) error {

	// If ship is docked, undock
	if ship.Nav.Status == "DOCKED" {
		client.OrbitShip(ship.Symbol)
	} else if ship.Nav.Status == "IN_TRANSIT" {
		// Wait for arrival
		t := time.Until(ship.Nav.Route.Arrival)
		fmt.Printf("%v: In transit, waiting %v to arrive\n", ship.Symbol, t)
		time.Sleep(t)

		// Refresh ship data
		ship, _ = client.GetMyShip(ship.Symbol)
	}

	// If ship isn't already there, go there
	if ship.Nav.WaypointSymbol != dest.Symbol {

		// Set thrusters to slow
		if ship.Fuel.Current == 0 {
			fmt.Printf("%v: No fuel to navigate, bailing\n", ship.Symbol)
			return fmt.Errorf("no fuel")
		} else if isLowFuel(ship.Fuel) && (ship.Nav.FlightMode != "DRIFT") {
			r1, _ := client.PatchShipNav(ship.Symbol, "DRIFT")
			fmt.Printf("%v: Set flight mode to %v due to low fuel (%v/%v)\n", ship.Symbol, r1.FlightMode, ship.Fuel.Current, ship.Fuel.Capacity)
		} else if ship.Nav.FlightMode != "CRUISE" {
			r1, _ := client.PatchShipNav(ship.Symbol, "CRUISE")
			fmt.Printf("%v: Set flight mode to %v\n", ship.Symbol, r1.FlightMode)
		}

		// Go to the destination
		r2, err := client.NavigateShip(ship.Symbol, dest.Symbol)
		if err != nil {
			return err
		}
		from := r2.Nav.Route.Departure
		to := r2.Nav.Route.Destination

		// Wait to get there
		t := time.Until(r2.Nav.Route.Arrival)

		fmt.Printf("%v: Navigating from %v (%v, %v) to %v (%v, %v), total %v units. Arriving in %v\n",
			ship.Symbol,
			from.Symbol, from.X, from.Y,
			to.Symbol, to.X, to.Y,
			routeDistance(r2.Nav.Route), t)

		time.Sleep(t)

		// Refuel
		if isWaypointMarket(dest) && isLowFuel(r2.Fuel) {
			fmt.Printf("%v: Refueling due to low fuel (%v/%v)\n", ship.Symbol, r2.Fuel.Current, r2.Fuel.Capacity)
			client.DockShip(ship.Symbol)
			client.RefuelShip(ship.Symbol)
			client.OrbitShip(ship.Symbol)
		}
	}

	return nil
}

func sellAllCargo(client *st.Client, ship stapi.Ship) error {
	// Update ship data
	ship, err := client.GetMyShip(ship.Symbol)
	if err != nil {
		return err
	}

	waypoints, err := client.GetSystemWaypoints(ship.Nav.SystemSymbol)
	if err != nil {
		return err
	}

	// Get the full waypoint info
	var wp stapi.Waypoint
	for _, v := range waypoints {
		if v.Symbol == ship.Nav.WaypointSymbol {
			wp = v
			break
		}
	}

	closest := findClosestWaypoint(wp, waypoints, isWaypointMarket)

	err = navAndRefuel(client, ship, closest)
	if err != nil {
		return err
	}

	// Try to dock, ignore error if already docked
	client.DockShip(ship.Symbol)

	market, _ := client.GetMarket(closest.SystemSymbol, closest.Symbol)

	total := int32(0)
	leftover := int32(0)
	for _, v := range ship.Cargo.Inventory {
		if isMarketBuying(market, v.Symbol) {
			r0, err := client.SellCargo(ship.Symbol, stapi.TradeSymbol(v.Symbol), v.Units)
			if err != nil {
				fmt.Printf("%v: Couldn't sell %v units of %v\n", ship.Symbol, v.Units, v.Symbol)
			} else {
				t := r0.Transaction
				total += t.TotalPrice
				fmt.Printf("%v: Sold %2v units of %16v at %2v for a total of %3v\n", ship.Symbol, t.Units, t.TradeSymbol, t.PricePerUnit, t.TotalPrice)
			}
			leftover = r0.Cargo.Units
		} else {
			fmt.Printf("%v: Market does not buy %v\n", ship.Symbol, v.Symbol)
		}
	}
	fmt.Printf("%v: Market trip total %4v; %v units left over\n", ship.Symbol, total, leftover)

	return nil
}

package bot

import (
	"fmt"
	"time"

	"github.com/bgreen/space-traders-go/st"
	"github.com/bgreen/space-traders-go/stapi"
)

type MinerBot struct {
	Target stapi.Waypoint
	Ship   stapi.Ship
	Client *st.Client
}

func (b MinerBot) RunOnce() msg {

	err := mineUntilFull(b.Client, b.Ship)
	if err != nil {
		return botErrorMsg(err)
	}

	err = sellAllCargo(b.Client, b.Ship)
	if err != nil {
		return botErrorMsg(err)
	}

	return botDoneMsg(b)
}

func mineUntilFull(client *st.Client, ship stapi.Ship) error {

	// Update ship data
	ship, err := client.GetMyShip(ship.Symbol)
	if err != nil {
		fmt.Println(err)
		return err
	}

	waypoints, err := client.GetSystemWaypoints(ship.Nav.SystemSymbol)
	if err != nil {
		fmt.Println(err)
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

	closest := findClosestWaypoint(wp, waypoints, isWaypointMine)

	err = navAndRefuel(client, ship, closest)
	if err != nil {
		return err
	}

	// Check if there's a cooldown active
	cd, _ := client.GetShipCooldown(ship.Symbol)
	if cd.RemainingSeconds > 0 {
		t := time.Until(*cd.Expiration)
		fmt.Printf("%v: Waiting %v for cooldown\n", ship.Symbol, t)
		time.Sleep(t)
	}

	total := int32(0)
	cycles := int32(0)
	for {
		r0, err := client.ExtractResources(ship.Symbol)
		if err != nil {
			return err
		}

		t := time.Until(*r0.Cooldown.Expiration)

		fmt.Printf("%v: Mined %2v units of %16v; cargo at %2v/%2v; cooldown %v\n", ship.Symbol,
			r0.Extraction.Yield.Units, r0.Extraction.Yield.Symbol,
			r0.Cargo.Units, r0.Cargo.Capacity, t)

		total += r0.Extraction.Yield.Units
		cycles += 1

		// Leave at 90% full, overfilling on ore is wasteful
		if r0.Cargo.Units >= (r0.Cargo.Capacity*9)/10 {
			break
		}

		time.Sleep(t)
	}

	fmt.Printf("%v: Mining trip total %2v units in %2v cycles\n", ship.Symbol, total, cycles)

	return nil
}

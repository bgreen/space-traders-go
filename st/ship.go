package st

import "github.com/bgreen/space-traders-go/stapi"

type Ship stapi.Ship

type ShipAction struct {
	Ref    *Ship
	Name   string
	Params map[string]string
}

var (
	createChartAction            = ShipAction{}
	createShipShipScanAction     = ShipAction{}
	createShipSystemScanAction   = ShipAction{}
	createShipWaypointScanAction = ShipAction{}
	createSurveyAction           = ShipAction{}

	dockShipAction = ShipAction{
		Name: "Dock",
	}

	extractResourcesAction = ShipAction{
		Name: "Extract",
	}

	installMountAction = ShipAction{}
	jettisonAction     = ShipAction{}
	jumpAction         = ShipAction{}

	navigateAction = ShipAction{
		Name:   "Navigate",
		Params: map[string]string{"waypoint": ""},
	}

	negotiateContractAction = ShipAction{}

	orbitAction = ShipAction{
		Name: "Orbit",
	}

	patchNavAction      = ShipAction{}
	purchaseCargoAction = ShipAction{}

	refuelAction = ShipAction{
		Name: "Refuel",
	}

	removeMountAction = ShipAction{}

	sellCargoAction = ShipAction{
		Name:   "Sell Cargo",
		Params: map[string]string{"type": "", "units": "0"},
	}

	refineAction        = ShipAction{}
	transferCargoAction = ShipAction{}
	warpAction          = ShipAction{}
)

// TODO: Really needs the whole model to determine
// 		 which actions are currently valid
func (s Ship) GetShipActions() []ShipAction {
	a := []ShipAction{}

	a = append(a, navigateAction)
	a = append(a, dockShipAction)
	a = append(a, orbitAction)
	a = append(a, extractResourcesAction)
	a = append(a, refuelAction)
	a = append(a, sellCargoAction)

	return a
}

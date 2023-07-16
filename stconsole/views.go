package stconsole

import (
	"fmt"
	"strings"
	"time"

	st "github.com/bgreen/space-traders-go/stservice"
	lipgloss "github.com/charmbracelet/lipgloss"
)

/////////////////////////////
// Utility
/////////////////////////////

func lineCount(s string) int {
	ss := strings.Split(s, "\n")
	return len(ss)
}

func stringWidth(s string) int {
	ss := strings.Split(s, "\n")
	max := 0
	for _, v := range ss {
		l := len(v) - 2 // newline + nul
		if l > max {
			max = l
		}
	}
	return max
}

func truncate(s string, l int) string {
	ss := strings.Split(s, "\n")
	for i, v := range ss {
		if len(v)-1 > l {
			ss[i] = v[:l]
		}
	}
	return strings.Join(ss, "\n")
}

/////////////////////////////
// Bars
/////////////////////////////

func (m model) statusBarView() string {
	s := fmt.Sprintf("Name: %v    HQ: %v    Credits: %v", m.agent.Symbol, m.agent.Headquarters, m.agent.Credits)
	return m.style.statusBarStyle.Render(s)
}

func (m model) messageView() string {
	return m.style.msgBarStyle.Render(m.msg)
}

/////////////////////////////
// Ships
/////////////////////////////

func (m model) shipsView() string {
	rows := []string{}
	for i, v := range m.ships {
		s := fmt.Sprintf("%v: %v", v.Symbol, v.Frame.Name[6:])
		if i == m.shipSel {
			rows = append(rows, m.style.rowSelectedStyle.Render(s))
		} else {
			rows = append(rows, m.style.rowUnselectedStyle.Render(s))
		}
	}
	return m.style.paneStyle.Render(lipgloss.JoinVertical(lipgloss.Left, rows...))
}

func (m model) shipInfoView() string {
	var s string
	if m.shipSel < len(m.ships) {
		a := m.ships[m.shipSel]
		s += m.style.rowTitleStyle.Render(a.Symbol) + "\n"
		s += fmt.Sprintf("Frame:    %v\n", a.Frame.Name[6:])
		s += fmt.Sprintf("Engine:   %v\n", a.Engine.Name)
		s += fmt.Sprintf("Fuel:     %v/%v\n", a.Fuel.Current, a.Fuel.Capacity)

		s += m.style.rowTitleStyle.Render("Nav") + "\n"
		s += fmt.Sprintf("System:   %v\n", a.Nav.SystemSymbol)
		s += fmt.Sprintf("Waypoint: %v\n", a.Nav.WaypointSymbol)
		s += fmt.Sprintf("Status:   %v\n", a.Nav.Status)
		s += fmt.Sprintf("Route:    %v\n", a.Nav.Route.Destination.Symbol)
		s += fmt.Sprintf("Mode:     %v\n", a.Nav.FlightMode)

		s += m.style.rowTitleStyle.Render("More Info") + "\n"
		for i, v := range []string{"Actions", "Modules", "Mounts"} {
			if m.shipInfoSel == i {
				s += m.style.rowSelectedStyle.Render(v) + "\n"
			} else {
				s += m.style.rowUnselectedStyle.Render(v) + "\n"
			}
		}

	}
	return m.style.paneStyle.Render(s)
}

func (m model) shipActionsView() string {
	var s string
	if m.shipSel < len(m.ships) {
		a := m.ships[m.shipSel]
		s += m.style.rowTitleStyle.Render("Actions") + "\n"
		actions := st.Ship(a).GetShipActions()
		for i, v := range actions {
			if i == m.shipActionSel {
				s += m.style.rowSelectedStyle.Render(v.Name) + "\n"
			} else {
				s += m.style.rowUnselectedStyle.Render(v.Name) + "\n"
			}
		}
		s += m.style.rowTitleStyle.Render("Options") + "\n"
		for k, v := range actions[m.shipActionSel].Params {
			param := fmt.Sprintf("%v: %v", k, v)
			s += m.style.rowUnselectedStyle.Render(param) + "\n"
		}
	}
	return m.style.paneStyle.Render(s)
}

func (m model) shipModulesView() string {
	var s string
	if m.shipSel < len(m.ships) {
		a := m.ships[m.shipSel]
		s += m.style.rowTitleStyle.Render("Modules") + "\n"
		for _, v := range a.Modules {
			s += fmt.Sprintf("%v\n", v.Name)
		}
	}
	return m.style.paneStyle.Render(s)
}

func (m model) shipMountsView() string {
	var s string
	if m.shipSel < len(m.ships) {
		a := m.ships[m.shipSel]
		s += m.style.rowTitleStyle.Render("Mounts") + "\n"
		for _, v := range a.Mounts {
			s += fmt.Sprintf("%v\n", v.Name)
		}
	}
	return m.style.paneStyle.Render(s)
}

/////////////////////////////
// Contracts
/////////////////////////////

func (m model) contractsView() string {
	rows := []string{}
	for i, v := range m.contracts {
		s := fmt.Sprintf("%v %v", v.FactionSymbol, v.Type)
		if i == m.contractSel {
			rows = append(rows, m.style.rowSelectedStyle.Render(s))
		} else {
			rows = append(rows, m.style.rowUnselectedStyle.Render(s))
		}
	}
	return m.style.paneStyle.Render(lipgloss.JoinVertical(lipgloss.Left, rows...))
}

func (m model) contractsInfoView() string {
	var s string
	if m.contractSel < len(m.contracts) {

		a := m.contracts[m.contractSel]
		s += m.style.rowTitleStyle.Render(a.FactionSymbol, " ", a.Type) + "\n"

		s += fmt.Sprintf("Type:     %v\n", a.Type)
		s += fmt.Sprintf("Accepted: %v\n", a.Accepted)
		s += fmt.Sprintf("Deadline: %v\n", a.DeadlineToAccept.Local().Format(time.Stamp))

		var wp []string
		wp = append(wp, m.style.rowTitleStyle.Render("Terms"))
		for _, v := range a.Terms.Deliver {
			wp = append(wp, fmt.Sprintf("Trade:  %v", v.TradeSymbol))
			wp = append(wp, fmt.Sprintf("Amount: %v/%v", v.UnitsFulfilled, v.UnitsRequired))
			wp = append(wp, fmt.Sprintf("Dest:   %v", v.DestinationSymbol))
		}
		s += lipgloss.JoinVertical(lipgloss.Left, wp...)
	}
	return m.style.paneStyle.Render(s)
}

/////////////////////////////
// Systems
/////////////////////////////

func (m model) systemsView() string {
	rows := []string{}
	for i, v := range m.systems {
		s := fmt.Sprintf("%v: %v", v.Symbol, v.Type)
		if i == m.systemSel {
			rows = append(rows, m.style.rowSelectedStyle.Render(s))
		} else {
			rows = append(rows, m.style.rowUnselectedStyle.Render(s))
		}
	}
	return m.style.paneStyle.Render(lipgloss.JoinVertical(lipgloss.Left, rows...))
}

func (m model) systemInfoView() string {
	var s string
	if m.systemSel < len(m.systems) {

		a := m.systems[m.systemSel]
		s += m.style.rowTitleStyle.Render(a.Symbol) + "\n"

		s += fmt.Sprintf("Type:   %v\n", a.Type)
		s += fmt.Sprintf("Coords: X:%v Y:%v\n", a.X, a.Y)
		shipCount := 0
		for _, v := range m.ships {
			if v.Nav.SystemSymbol == a.Symbol {
				shipCount += 1
			}
		}
		s += fmt.Sprintf("Ships:  %v\n", shipCount)
		var wp []string

		wp = append(wp, m.style.rowTitleStyle.Render("Waypoints"))
		for i, v := range a.Waypoints {
			w := truncate(fmt.Sprintf("%v: %v", v.Symbol, v.Type), m.style.paneWidth-2)
			if i == m.wpListSel {
				wp = append(wp, m.style.rowSelectedStyle.Render(w))
			} else {
				wp = append(wp, m.style.rowUnselectedStyle.Render(w))
			}
		}
		s += lipgloss.JoinVertical(lipgloss.Left, wp...)
	}
	return m.style.paneStyle.Render(s)
}

func (m model) systemWaypointInfoView() string {
	var s string

	if m.systemSel < len(m.systems) {
		if m.wpListSel < len(m.systems[m.systemSel].Waypoints) {
			a := m.systems[m.systemSel].Waypoints[m.wpListSel]
			s += m.style.rowTitleStyle.Render(a.Symbol) + "\n"
			s += fmt.Sprintf("Type:   %v\n", a.Type)
			s += fmt.Sprintf("Coords: X:%v Y:%v\n", a.X, a.Y)
			if v, ok := m.waypoints[a.Symbol]; ok {
				s += m.style.rowTitleStyle.Render("Traits") + "\n"
				for _, w := range v.Traits {
					s += fmt.Sprintf("%v\n", w.Name)
				}
				s += m.style.rowTitleStyle.Render("Orbitals") + "\n"
				for _, w := range v.Orbitals {
					s += fmt.Sprintf("%v: %v\n", w.Symbol, m.waypoints[w.Symbol].Type)
				}
			}
		}
	}
	return m.style.paneStyle.Render(s)
}

/////////////////////////////
// Buttons
/////////////////////////////

func (m model) buttonsView() string {
	buttons := []string{}
	for i, v := range []string{"Ships", "Systems", "Contracts"} {
		if m.modeSel == i {
			buttons = append(buttons, m.style.buttonSelectedStyle.Render(v))
		} else {
			buttons = append(buttons, m.style.buttonUnselectedStyle.Render(v))
		}
	}
	return lipgloss.JoinVertical(lipgloss.Left, buttons...)
}

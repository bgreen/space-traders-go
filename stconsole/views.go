package main

import (
	"fmt"
	"time"

	lipgloss "github.com/charmbracelet/lipgloss"
)

func (m model) statusBarView() string {
	s := fmt.Sprintf("Name: %v    HQ: %v    Credits: %v", m.agent.Symbol, m.agent.Headquarters, m.agent.Credits)
	return m.style.statusBarStyle.Render(s)
}

func (m model) shipsView() string {
	rows := []string{}
	for i, v := range m.ships {
		s := fmt.Sprintf("%v: %v", v.Symbol, v.Frame.Name)
		if i == m.shipSel {
			rows = append(rows, m.style.rowSelectedStyle.Render(s))
		} else {
			rows = append(rows, m.style.rowUnselectedStyle.Render(s))
		}
	}
	return m.style.paneStyle.Render(lipgloss.JoinVertical(lipgloss.Left, rows...))
}

func (m model) contractsView() string {
	var s string
	for _, v := range m.contracts {
		timer := v.DeadlineToAccept.Sub(time.Now())
		s += fmt.Sprintf("%v: %v\n", v.Type, timer)
	}
	return m.style.paneStyle.Render(s)
}

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

func (m model) messageView() string {
	return m.style.msgBarStyle.Render(m.msg)
}

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

func (m model) shipInfoView() string {
	var s string
	if m.shipSel < len(m.ships) {
		a := m.ships[m.shipSel]
		s += m.style.rowTitleStyle.Render(a.Symbol) + "\n"
		s += fmt.Sprintf("Frame:    %v\n", a.Frame.Name)
		s += fmt.Sprintf("Engine:   %v\n", a.Engine.Name)
		s += fmt.Sprintf("Fuel:     %v/%v\n", a.Fuel.Current, a.Fuel.Capacity)

		s += m.style.rowTitleStyle.Render("Nav") + "\n"
		s += fmt.Sprintf("System:   %v\n", a.Nav.SystemSymbol)
		s += fmt.Sprintf("Waypoint: %v\n", a.Nav.WaypointSymbol)
		s += fmt.Sprintf("Status:   %v\n", a.Nav.Status)
		s += fmt.Sprintf("Route:    %v\n", a.Nav.Route.Destination.Symbol)
		s += fmt.Sprintf("Mode:     %v\n", a.Nav.FlightMode)

		s += m.style.rowTitleStyle.Render("Modules") + "\n"
		for _, v := range a.Modules {
			s += fmt.Sprintf("%v\n", v.Name)
		}

		s += m.style.rowTitleStyle.Render("Mounts") + "\n"
		for _, v := range a.Mounts {
			s += fmt.Sprintf("%v\n", v.Name)
		}

		s += m.style.rowTitleStyle.Render("Actions") + "\n"
	}
	return m.style.paneStyle.Render(s)
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
		for _, v := range a.Waypoints {
			wp = append(wp, fmt.Sprintf("%v: %v", v.Symbol, v.Type))
		}
		s += lipgloss.JoinVertical(lipgloss.Left, wp...)
	}
	return m.style.paneStyle.Render(s)
}

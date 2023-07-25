package stconsole

import (
	"fmt"
	"os"

	st "github.com/bgreen/space-traders-go/st"
	stapi "github.com/bgreen/space-traders-go/stapi"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

var client *st.Client

func Run() {

	m := model{win: coords{x: 120, y: 50},
		waypoints: make(map[string]stapi.Waypoint)}

	// Create API Client
	client = st.NewClient()
	client.Start()
	defer client.Stop()
	// Create the bubbletea app
	p := tea.NewProgram(m, tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

/////////////////////////
//	Model definition
/////////////////////////

type model struct {
	agent     stapi.Agent
	ships     []stapi.Ship
	contracts []stapi.Contract
	systems   []stapi.System
	waypoints map[string]stapi.Waypoint

	modeSel       int
	shipSel       int
	shipInfoSel   int
	shipActionSel int
	systemSel     int
	wpListSel     int
	contractSel   int

	activeX int
	activeY int

	win   coords
	style style

	msg string
}

func (m model) Init() tea.Cmd {
	return tea.Batch(getAgent, getShips, getContracts, getSystems)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case agentMsg:
		m.agent = stapi.Agent(msg)

	case shipsMsg:
		m.ships = []stapi.Ship(msg)
		var r []tea.Cmd
		for _, v := range m.ships {
			found := false
			for _, w := range m.systems {
				if w.Symbol == v.Nav.SystemSymbol {
					found = true
					break
				}
			}
			if !found {
				r = append(r, getSystem(v.Nav.SystemSymbol))
			}
		}

		return m, tea.Batch(r...)

	case contractsMsg:
		m.contracts = []stapi.Contract(msg)

	case systemsMsg:
		var r []tea.Cmd
		for _, v := range []stapi.System(msg) {
			found := false
			for _, w := range m.systems {
				if w.Symbol == v.Symbol {
					found = true
					break
				}
			}
			if !found {
				m.systems = append(m.systems, v)
				r = append(r, getSystemWaypoints(v.Symbol))
			}
		}
		return m, tea.Batch(r...)

	case systemMsg:
		var r []tea.Cmd
		found := false
		for _, v := range m.systems {
			if v.Symbol == msg.Symbol {
				found = true
				break
			}
		}
		if !found {
			m.systems = append(m.systems, stapi.System(msg))
			r = append(r, getSystemWaypoints(stapi.System(msg).Symbol))
		}
		return m, tea.Batch(r...)

	case waypointsMsg:
		for _, v := range []stapi.Waypoint(msg) {
			m.waypoints[v.Symbol] = v
		}

	case errMsg:
		m.msg = fmt.Sprint(msg.err)

	case tea.WindowSizeMsg:
		m.win.x = msg.Width
		m.win.y = msg.Height
		m.style = m.resetStyle()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			return m, tea.Quit
		case "esc":
			return m, tea.Quit
		case "up":
			m.activeY -= 1
			if m.activeY < 0 {
				m.activeY = 0
			}
		case "down":
			m.activeY += 1
			if m.activeX > 10 {
				m.activeX = 10
			}
		case "left":
			m.activeX -= 1
			if m.activeX < 0 {
				m.activeX = 0
			}
		case "right":
			m.activeX += 1
			if m.activeY > 3 {
				m.activeY = 3
			}
		}
		return m, nil

	case tea.MouseMsg:
		e := tea.MouseEvent(msg)
		switch e.Type {
		case tea.MouseLeft:
			m.msg = fmt.Sprintf("X:%v Y%v", e.X, e.Y)
			switch {
			case m.style.buttonLoc[0].contains(e.X, e.Y):
				m.modeSel = 0
			case m.style.buttonLoc[1].contains(e.X, e.Y):
				m.modeSel = 1
			case m.style.buttonLoc[2].contains(e.X, e.Y):
				m.modeSel = 2
			case m.style.paneLoc[0].contains(e.X, e.Y):
				sel := (e.Y - 2)
				if sel >= 0 {
					if (m.modeSel == 0) && (sel < len(m.ships)) {
						m.shipSel = sel
					}
					if (m.modeSel == 1) && (sel < len(m.systems)) {
						m.systemSel = sel
					}
					if (m.modeSel == 2) && (sel < len(m.contracts)) {
						m.contractSel = sel
					}
				}
			case m.style.paneLoc[1].contains(e.X, e.Y):
				if m.modeSel == 0 {
					sel := (e.Y - m.style.shipInfoListLoc.topLeft.y)
					if (sel >= 0) && (sel < 3) {
						m.shipInfoSel = sel
					}
				} else if m.modeSel == 1 {
					sel := (e.Y - m.style.wpListLoc.topLeft.y)
					if (sel >= 0) && (sel < len(m.systems[m.systemSel].Waypoints)) {
						m.wpListSel = sel
					}
				}
			case m.style.paneLoc[2].contains(e.X, e.Y):
				if m.modeSel == 0 {
					sel := (e.Y - 3)
					if (sel >= 0) && (m.shipSel < len(m.ships)) {
						actions := st.Ship(m.ships[m.shipSel]).GetShipActions()
						if sel < len(actions) {
							m.shipActionSel = sel
						}
					}
				}
			}
		}
	}

	return m, nil
}

func (m model) View() string {

	panes := []string{}

	panes = append(panes, m.buttonsView())

	if m.modeSel == 0 {
		panes = append(panes, m.shipsView())
		panes = append(panes, m.shipInfoView())
		switch m.shipInfoSel {
		case 0:
			panes = append(panes, m.shipActionsView())
		case 1:
			panes = append(panes, m.shipModulesView())
		case 2:
			panes = append(panes, m.shipMountsView())
		default:
		}
	} else if m.modeSel == 1 {
		panes = append(panes, m.systemsView())
		panes = append(panes, m.systemInfoView())
		panes = append(panes, m.systemWaypointInfoView())
	} else if m.modeSel == 2 {
		panes = append(panes, m.contractsView())
		panes = append(panes, m.contractsInfoView())
	}

	s := lipgloss.JoinHorizontal(lipgloss.Top, panes...)

	s = lipgloss.JoinVertical(lipgloss.Left,
		m.statusBarView(),
		s,
		m.messageView())

	return s
}

///////////////////////////
//	Msg Definitions
///////////////////////////

type agentMsg stapi.Agent

type shipsMsg []stapi.Ship

type contractsMsg []stapi.Contract

type systemsMsg []stapi.System

type systemMsg stapi.System

type waypointsMsg []stapi.Waypoint

type errMsg struct{ err error }

///////////////////////////
// 	 Cmd definitions
///////////////////////////

func getAgent() tea.Msg {
	a, err := client.GetMyAgent()
	if err != nil {
		return errMsg{err}
	}
	return agentMsg(a)
}

func getShips() tea.Msg {
	a, err := client.GetMyShips()
	if err != nil {
		return errMsg{err}
	}
	return shipsMsg(a)
}

func getContracts() tea.Msg {
	a, err := client.GetContracts()
	if err != nil {
		return errMsg{err}
	}
	return contractsMsg(a)
}

func getSystems() tea.Msg {
	a, err := client.GetSystems()
	if err != nil {
		return errMsg{err}
	}
	return systemsMsg(a)
}

func getSystem(symbol string) func() tea.Msg {
	return func() tea.Msg {
		a, err := client.GetSystem(symbol)
		if err != nil {
			return errMsg{err}
		}
		return systemMsg(a)
	}
}

func getSystemWaypoints(symbol string) func() tea.Msg {
	return func() tea.Msg {
		a, err := client.GetSystemWaypoints(symbol)
		if err != nil {
			return errMsg{err}
		}
		return waypointsMsg(a)
	}
}

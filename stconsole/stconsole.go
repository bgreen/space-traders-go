package main

import (
	"fmt"
	"os"

	stapi "github.com/bgreen/space-traders-go/stapi"
	st "github.com/bgreen/space-traders-go/sthandler"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

var client *st.Server

func main() {

	m := model{systems: make(map[string]stapi.System)}

	// Create API Client
	client = st.NewServer("PICKYPICKY")
	client.Start()
	defer client.Stop()
	// Create the bubbletea app
	p := tea.NewProgram(m, tea.WithMouseAllMotion())
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
	systems   map[string]stapi.System
	modeSel   int

	win window

	msg string
}

type window struct {
	x int
	y int
}

func (m model) Init() tea.Cmd {
	return tea.Batch(getAgent, getShips, getContracts)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case agentMsg:
		m.agent = stapi.Agent(msg)

	case shipsMsg:
		m.ships = []stapi.Ship(msg)
		var r []tea.Cmd
		for _, v := range m.ships {
			if _, ok := m.systems[v.Nav.SystemSymbol]; !ok {
				r = append(r, getSystem(v.Nav.SystemSymbol))
			}
		}
		return m, tea.Batch(r...)

	case contractsMsg:
		m.contracts = []stapi.Contract(msg)

	case systemsMsg:
		for _, v := range []stapi.System(msg) {
			m.systems[v.Symbol] = v
		}

	case systemMsg:
		s := stapi.System(msg)
		m.systems[s.Symbol] = stapi.System(s)

	case errMsg:
		m.msg = fmt.Sprint(msg.err)

	case tea.WindowSizeMsg:
		m.win.x = msg.Width
		m.win.y = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			return m, tea.Quit
		case "esc":
			return m, tea.Quit
		case " ":
			return m, getAgent
		}

	case tea.MouseMsg:
		e := tea.MouseEvent(msg)
		e.Y = (m.win.y - 1) - e.Y
		switch e.Type {
		case tea.MouseLeft:
			m.msg = fmt.Sprintf("X:%v Y%v", e.X, e.Y)
			switch {
			case e.X >= 0 && e.X < buttonWidth:
				switch {
				case (e.Y >= 1+buttonHeight*0) && (e.Y < 1+buttonHeight*1):
					m.modeSel = 2
				case (e.Y >= 1+buttonHeight*1) && (e.Y < 1+buttonHeight*2):
					m.modeSel = 1
				case (e.Y >= 1+buttonHeight*2) && (e.Y < 1+buttonHeight*3):
					m.modeSel = 0
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
	} else if m.modeSel == 1 {
		panes = append(panes, m.systemsView())
	} else if m.modeSel == 2 {
		panes = append(panes, m.contractsView())
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

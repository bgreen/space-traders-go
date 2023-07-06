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

	m := model{}

	// Create API Client
	client = st.NewServer("PICKYPICKY")
	client.Start()
	defer client.Stop()
	// Create the bubbletea app
	p := tea.NewProgram(m)
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

	err error
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

	case contractsMsg:
		m.contracts = []stapi.Contract(msg)

	case systemsMsg:
		m.systems = []stapi.System(msg)

	case errMsg:
		m.err = msg.err
		return m, tea.Quit

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeySpace:
			return m, getAgent
		}
	}

	return m, nil
}

func (m model) View() string {

	s := lipgloss.JoinHorizontal(lipgloss.Top,
		m.shipsView(),
		m.contractsView(),
		m.systemsView())

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

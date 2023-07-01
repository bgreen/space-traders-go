package main

import (
	"context"
	"fmt"
	"os"

	stapi "github.com/bgreen/space-traders-go/stapi"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

var client *stapi.APIClient
var ctx context.Context

func main() {
	// Read Bearer token from token.txt
	token, err := os.ReadFile("token.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "No token file found")
		return
	}

	m := model{}

	// Add bearer token to the context
	ctx = context.WithValue(context.Background(), stapi.ContextAccessToken, string(token))

	// Create API Client
	configuration := stapi.NewConfiguration()
	client = stapi.NewAPIClient(configuration)

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

	err error
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

	case contractsMsg:
		m.contracts = []stapi.Contract(msg)

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

var (
	paneStyle = lipgloss.NewStyle().
			Width(30).
			Height(12).
			Align(lipgloss.Top, lipgloss.Left).
			BorderStyle(lipgloss.NormalBorder())

	errStyle = lipgloss.NewStyle().
			Height(2).
			Align(lipgloss.Bottom, lipgloss.Left)
)

func (m model) View() string {

	agentPane := fmt.Sprintf("Agent:\nName: %v\nHQ: %v\nCredits: %v\n", m.agent.Symbol, m.agent.Headquarters, m.agent.Credits)

	shipPane := "Ships:\n"
	for i, v := range m.ships {
		shipPane += fmt.Sprintf("%v: %v\n", i, v.Symbol)
	}

	contractPane := "Contracts:\n"
	for i, v := range m.contracts {
		contractPane += fmt.Sprintf("%v: %v %v\n", i, v.Id, v.Type)
	}

	var errPane = ""
	if m.err != nil {
		errPane = fmt.Sprintf("%s", m.err)
	}

	s := lipgloss.JoinHorizontal(lipgloss.Top,
		paneStyle.Render(agentPane),
		paneStyle.Render(shipPane),
		paneStyle.Render(contractPane))

	s = lipgloss.JoinVertical(lipgloss.Left,
		s,
		errStyle.Render(errPane))

	return s
}

///////////////////////////
//	Msg Definitions
///////////////////////////

type agentMsg stapi.Agent

type shipsMsg []stapi.Ship

type contractsMsg []stapi.Contract

type errMsg struct{ err error }

///////////////////////////
// 	 Cmd definitions
///////////////////////////

func getAgent() tea.Msg {
	// Get Agent from API
	resp, _, err := client.AgentsApi.GetMyAgent(ctx).Execute()
	if err != nil {
		return errMsg{err}
	}

	return agentMsg(resp.GetData())
}

func getShips() tea.Msg {
	resp, _, err := client.FleetApi.GetMyShips(ctx).Execute()
	if err != nil {
		return errMsg{err}
	}

	return shipsMsg(resp.GetData())
}

func getContracts() tea.Msg {
	resp, _, err := client.ContractsApi.GetContracts(ctx).Execute()
	if err != nil {
		return errMsg{err}
	}

	return contractsMsg(resp.GetData())
}

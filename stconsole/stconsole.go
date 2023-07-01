package main

import (
	"context"
	"fmt"
	"os"

	stapi "github.com/bgreen/space-traders-go/stapi"
	tea "github.com/charmbracelet/bubbletea"
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
	agent stapi.Agent
}

func (m model) Init() tea.Cmd {
	return getAgent
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case agentMsg:
		m.agent = stapi.Agent(msg)

	case errMsg:
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

	s := fmt.Sprintf("Name: %v\nHQ: %v\nCredits: %v\n", m.agent.Symbol, m.agent.Headquarters, m.agent.Credits)

	return s
}

///////////////////////////
//	Msg Definitions
///////////////////////////

type agentMsg stapi.Agent

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

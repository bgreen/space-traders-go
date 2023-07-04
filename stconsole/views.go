package main

import (
	"fmt"

	lipgloss "github.com/charmbracelet/lipgloss"
)

func (m model) statusBarView() string {
	s := fmt.Sprintf("Name: %v    HQ: %v    Credits: %v", m.agent.Symbol, m.agent.Headquarters, m.agent.Credits)
	return barStyle.AlignHorizontal(lipgloss.Center).Render(s)
}

func (m model) shipsView() string {
	s := "Ships:\n"
	for _, v := range m.ships {
		s += fmt.Sprintf("%v: %v\n", v.Symbol, v.Frame.Name)
	}
	return paneStyle.Render(s)
}

func (m model) contractsView() string {
	s := "Contracts:\n"
	for _, v := range m.contracts {
		s += fmt.Sprintf("%v: %v\n", v.Id, v.Type)
	}
	return paneStyle.Render(s)
}

func (m model) systemsView() string {
	s := "Systems:\n"
	for _, v := range m.systems {
		s += fmt.Sprintf("%v: %v\n", v.Symbol, v.Type)
	}
	return paneStyle.Render(s)
}

func (m model) messageView() string {
	s := ""
	if m.err != nil {
		s = fmt.Sprintf("%s", m.err)
	}
	return barStyle.Render(s)
}

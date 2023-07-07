package main

import (
	"fmt"

	lipgloss "github.com/charmbracelet/lipgloss"
)

func (m model) statusBarView() string {
	s := fmt.Sprintf("Name: %v    HQ: %v    Credits: %v", m.agent.Symbol, m.agent.Headquarters, m.agent.Credits)
	return statusBarStyle.Render(s)
}

func (m model) shipsView() string {
	var s string
	for _, v := range m.ships {
		s += fmt.Sprintf("%v: %v\n", v.Symbol, v.Frame.Name)
	}
	return paneStyle.Render(s)
}

func (m model) contractsView() string {
	var s string
	for _, v := range m.contracts {
		s += fmt.Sprintf("%v: %v\n", v.Id, v.Type)
	}
	return paneStyle.Render(s)
}

func (m model) systemsView() string {
	var s string
	for _, v := range m.systems {
		s += fmt.Sprintf("%v: %v\n", v.Symbol, v.Type)
	}
	return paneStyle.Render(s)
}

func (m model) messageView() string {
	return msgBarStyle.Render(m.msg)
}

func (m model) buttonsView() string {
	buttons := []string{}
	for i, v := range []string{"Ships", "Systems", "Contracts"} {
		if m.modeSel == i {
			buttons = append(buttons, buttonSelectedStyle.Render(v))
		} else {
			buttons = append(buttons, buttonUnselectedStyle.Render(v))
		}
	}
	return lipgloss.JoinVertical(lipgloss.Left, buttons...)
}

package stconsole

import (
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

/////////////////////////////
// row
/////////////////////////////

type row struct {
	style lipgloss.Style
	s     string
}

func (m row) Init() tea.Cmd {
	return nil
}

func (m row) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m row) View() string {
	return m.style.Render(m.s)
}

/////////////////////////////
// selectList
/////////////////////////////

type selectList struct {
	style    lipgloss.Style
	selStyle lipgloss.Style
	selected int
	rows     []row
	bounds   box
}

func (m selectList) Init() tea.Cmd {
	return nil
}

func (m selectList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: should MouseEvents/Bounds coordinates be relative or absolute?
	// Should be relative to parent container
	return m, nil
}

func (m selectList) View() string {
	var s string
	for _, r := range m.rows {
		s = lipgloss.JoinVertical(lipgloss.Left, s, r.View())
	}
	return s
}

/////////////////////////////
// pane
/////////////////////////////

type pane struct {
	style  lipgloss.Style
	rows   []tea.Model
	bounds box
}

func (m pane) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, v := range m.rows {
		cmds = append(cmds, v.Init())
	}
	return tea.Batch(cmds...)
}

func (m pane) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.MouseMsg:
		/*
			e := tea.MouseEvent(msg)
			for i, v := range m.rows {
				if v.bounds.contains(e.X, e.Y) {
					e_rel := e
					e_rel.X = v.bounds.X
					var c tea.Cmd
					m.rows[i], c = v.Update(msg)
					cmds = append(cmds, c)
				}
			}
		*/
	default:
		for i, v := range m.rows {
			var c tea.Cmd
			m.rows[i], c = v.Update(msg)
			cmds = append(cmds, c)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m pane) View() string {
	var s string
	for _, v := range m.rows {
		s = lipgloss.JoinVertical(lipgloss.Left, s, v.View())
	}
	return m.style.Render(s)
}

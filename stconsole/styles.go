package main

import lipgloss "github.com/charmbracelet/lipgloss"

type style struct {
	totalWidth   int
	totalHeight  int
	buttonWidth  int
	buttonHeight int
	paneHeight   int
	paneWidth    int
	paneCount    int

	color []lipgloss.Color

	paneStyle             lipgloss.Style
	statusBarStyle        lipgloss.Style
	msgBarStyle           lipgloss.Style
	buttonSelectedStyle   lipgloss.Style
	buttonUnselectedStyle lipgloss.Style
	rowSelectedStyle      lipgloss.Style
	rowUnselectedStyle    lipgloss.Style
	rowTitleStyle         lipgloss.Style
}

func (m model) resetStyle() style {
	var s style = m.style
	s.totalWidth = m.win.x
	s.totalHeight = m.win.y
	s.buttonWidth = 11
	s.buttonHeight = 5
	s.paneCount = 3
	s.paneHeight = s.totalHeight - 2
	s.paneWidth = (s.totalWidth - s.buttonWidth) / s.paneCount

	// https://colorhunt.co/palette/0000004e4feb068fffeeeeee
	s.color = []lipgloss.Color{lipgloss.Color("#000000"), // Black
		lipgloss.Color("#4E4FEB"), // Dark Blue
		lipgloss.Color("#068FFF"), // Light Blue
		lipgloss.Color("#EEEEEE")} // Light Gray

	s.paneStyle = lipgloss.NewStyle().
		Width(s.paneWidth - 2).
		Height(s.paneHeight - 2).
		AlignHorizontal(lipgloss.Left).
		AlignVertical(lipgloss.Top).
		BorderStyle(lipgloss.NormalBorder()).
		Background(s.color[0]).
		Foreground(s.color[3]).
		BorderBackground(s.color[0]).
		BorderForeground(s.color[3])

	barStyle := lipgloss.NewStyle().
		Height(1).
		Width(s.totalWidth).
		AlignHorizontal(lipgloss.Left).
		AlignVertical(lipgloss.Bottom).
		Background(s.color[3]).
		Foreground(s.color[0])

	s.statusBarStyle = barStyle.Copy().
		AlignHorizontal(lipgloss.Center)

	s.msgBarStyle = barStyle.Copy().
		AlignHorizontal(lipgloss.Left)

	buttonStyle := lipgloss.NewStyle().
		Height(s.buttonHeight - 2).
		Width(s.buttonWidth - 2).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		BorderStyle(lipgloss.NormalBorder())

	s.buttonSelectedStyle = buttonStyle.Copy().
		Background(s.color[1]).
		Foreground(s.color[3]).
		BorderBackground(s.color[0]).
		BorderForeground(s.color[3])

	s.buttonUnselectedStyle = buttonStyle.Copy().
		Background(s.color[0]).
		Foreground(s.color[3]).
		BorderBackground(s.color[0]).
		BorderForeground(s.color[3])

	s.rowSelectedStyle = lipgloss.NewStyle().
		Width(s.paneWidth - 2).
		Background(s.color[1]).
		Foreground(s.color[3])

	s.rowUnselectedStyle = lipgloss.NewStyle().
		Width(s.paneWidth - 2).
		Background(s.color[0]).
		Foreground(s.color[3])

	s.rowTitleStyle = lipgloss.NewStyle().
		Width(s.paneWidth - 2).
		AlignHorizontal(lipgloss.Center).
		Background(s.color[3]).
		Foreground(s.color[0])

	return s
}

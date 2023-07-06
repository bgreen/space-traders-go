package main

import lipgloss "github.com/charmbracelet/lipgloss"

var (
	paneStyle = lipgloss.NewStyle().
			Width(30).
			Height(12).
			AlignHorizontal(lipgloss.Left).
			AlignVertical(lipgloss.Top).
			BorderStyle(lipgloss.NormalBorder())

	barStyle = lipgloss.NewStyle().
			Height(1).
			Width(90).
			AlignHorizontal(lipgloss.Left).
			AlignVertical(lipgloss.Bottom).
			Background(lipgloss.Color("7")).
			Foreground(lipgloss.Color("0"))

	buttonStyle = lipgloss.NewStyle().
			Height(4).
			Width(8).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Padding(1)

	buttonUnselectedStyle = buttonStyle.
				Background(lipgloss.Color("7")).
				Foreground(lipgloss.Color("0"))

	buttonSelectedStyle = buttonStyle.
				Background(lipgloss.Color("6")).
				Foreground(lipgloss.Color("0"))
)

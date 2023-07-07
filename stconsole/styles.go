package main

import lipgloss "github.com/charmbracelet/lipgloss"

var (
	totalWidth   = 99
	buttonWidth  = 11
	buttonHeight = 5
	paneCount    = 3

	paneStyle = lipgloss.NewStyle().
			Width(((totalWidth - buttonWidth) / paneCount) - 2).
			Height(13).
			AlignHorizontal(lipgloss.Left).
			AlignVertical(lipgloss.Top).
			BorderStyle(lipgloss.NormalBorder())

	statusBarStyle = lipgloss.NewStyle().
			Height(1).
			Width(totalWidth).
			AlignHorizontal(lipgloss.Left).
			AlignVertical(lipgloss.Bottom).
			Background(lipgloss.Color("7")).
			Foreground(lipgloss.Color("0")).
			AlignHorizontal(lipgloss.Center)

	msgBarStyle = lipgloss.NewStyle().
			Height(1).
			Width(totalWidth).
			AlignHorizontal(lipgloss.Left).
			AlignVertical(lipgloss.Bottom).
			Background(lipgloss.Color("7")).
			Foreground(lipgloss.Color("0")).
			AlignHorizontal(lipgloss.Left)

	buttonSelectedStyle = lipgloss.NewStyle().
				Height(buttonHeight - 2).
				Width(buttonWidth - 2).
				AlignHorizontal(lipgloss.Center).
				AlignVertical(lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				Background(lipgloss.Color("6")).
				Foreground(lipgloss.Color("0"))

	buttonUnselectedStyle = lipgloss.NewStyle().
				Height(buttonHeight - 2).
				Width(buttonWidth - 2).
				AlignHorizontal(lipgloss.Center).
				AlignVertical(lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				Background(lipgloss.Color("7")).
				Foreground(lipgloss.Color("0"))
)

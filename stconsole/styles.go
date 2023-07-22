package stconsole

import lipgloss "github.com/charmbracelet/lipgloss"

type style struct {
	totalWidth   int
	totalHeight  int
	buttonWidth  int
	buttonHeight int
	paneHeight   int
	paneWidth    int
	paneCount    int

	buttonLoc       [3]box
	paneLoc         [3]box
	shipInfoListLoc box
	wpListLoc       box

	color []lipgloss.Color

	paneStyle             lipgloss.Style
	paneActiveStyle       lipgloss.Style
	statusBarStyle        lipgloss.Style
	msgBarStyle           lipgloss.Style
	buttonSelectedStyle   lipgloss.Style
	buttonUnselectedStyle lipgloss.Style
	rowSelectedStyle      lipgloss.Style
	rowUnselectedStyle    lipgloss.Style
	rowActiveStyle        lipgloss.Style
	rowTitleStyle         lipgloss.Style
}

type coords struct {
	x int
	y int
}

type box struct {
	topLeft     coords
	bottomRight coords
	contents    string
	orientation int
	more        []box
	style       lipgloss.Style
}

func (b box) contains(x, y int) bool {
	if (x >= b.topLeft.x) && (x <= b.bottomRight.x) && (y >= b.topLeft.y) && (y <= b.bottomRight.y) {
		return true
	}
	return false
}

func (m model) resetStyle() style {
	var s style = m.style
	//s.totalWidth = m.win.x
	s.totalWidth = 95 // TODO: fix bubbletea mouse width limit
	s.totalHeight = m.win.y
	s.buttonWidth = 11
	s.buttonHeight = 5
	s.paneCount = 3
	s.paneHeight = s.totalHeight - 2
	s.paneWidth = (s.totalWidth - s.buttonWidth) / s.paneCount

	for i := range s.buttonLoc {
		var b box
		b.topLeft.x = 0
		b.topLeft.y = 1 + s.buttonHeight*i
		b.bottomRight.x = b.topLeft.x + s.buttonWidth - 1
		b.bottomRight.y = b.topLeft.y + s.buttonHeight - 1
		s.buttonLoc[i] = b
	}

	for i := range s.paneLoc {
		var b box
		b.topLeft.x = 0 + s.buttonWidth + s.paneWidth*i
		b.topLeft.y = 1
		b.bottomRight.x = b.topLeft.x + s.paneWidth - 1
		b.bottomRight.y = b.topLeft.y + s.paneHeight - 1
		s.paneLoc[i] = b
	}

	s.shipInfoListLoc.topLeft.x = s.paneLoc[1].topLeft.x
	s.shipInfoListLoc.topLeft.y = s.paneLoc[1].topLeft.y + 12
	s.shipInfoListLoc.bottomRight.x = s.paneLoc[1].bottomRight.x
	s.shipInfoListLoc.bottomRight.y = s.paneLoc[1].bottomRight.y - 1

	s.wpListLoc.topLeft.x = s.paneLoc[1].topLeft.x
	s.wpListLoc.topLeft.y = s.paneLoc[1].topLeft.y + 6
	s.wpListLoc.bottomRight.x = s.paneLoc[1].bottomRight.x
	s.wpListLoc.bottomRight.y = s.paneLoc[1].bottomRight.y - 1

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

	s.paneActiveStyle = s.paneStyle.Copy().
		BorderBackground(s.color[0]).
		BorderForeground(s.color[2])

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

	s.rowActiveStyle = lipgloss.NewStyle().
		Width(s.paneWidth - 2).
		Background(s.color[2]).
		Foreground(s.color[0])

	s.rowTitleStyle = lipgloss.NewStyle().
		Width(s.paneWidth - 2).
		AlignHorizontal(lipgloss.Center).
		Background(s.color[3]).
		Foreground(s.color[0])

	return s
}

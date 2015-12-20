package main

type menu struct {
	items         []menuItem
	selectedIndex int
}

type menuItem int

const (
	menuContinueGame menuItem = iota
	menuNewGame
)

var menuItemText = map[menuItem]string{
	menuContinueGame: "C O N T I N U E  G A M E",
	menuNewGame:      "N E W  G A M E",
}

func newMenu() *menu {
	return &menu{
		items: []menuItem{
			menuNewGame,
		},
	}
}

func (m *menu) moveDown() {
	m.selectedIndex = (m.selectedIndex + 1) % len(m.items)
}

func (m *menu) moveUp() {
	if m.selectedIndex -= 1; m.selectedIndex < 0 {
		m.selectedIndex = len(m.items) - 1
	}
}

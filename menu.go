package main

type menu struct {
	items         []menuItem
	selectedIndex int
}

type menuItem int

const (
	menuContinueGame menuItem = iota
	menuNewGame
	menuExit
)

var menuItemText = map[menuItem]string{
	menuContinueGame: "C O N T I N U E  G A M E",
	menuNewGame:      "N E W  G A M E",
	menuExit:         "E X I T",
}

func newMenu() *menu {
	return &menu{
		items: []menuItem{
			menuNewGame,
			menuExit,
		},
	}
}

func (m *menu) addContinueGame() {
	m.items = []menuItem{
		menuContinueGame,
		menuNewGame,
		menuExit,
	}
	m.selectedIndex = 0
}

func (m *menu) removeContinueGame() {
	m.items = []menuItem{
		menuNewGame,
		menuExit,
	}
	m.selectedIndex = 0
}

func (m *menu) moveDown() {
	m.selectedIndex = (m.selectedIndex + 1) % len(m.items)
}

func (m *menu) moveUp() {
	if m.selectedIndex -= 1; m.selectedIndex < 0 {
		m.selectedIndex = len(m.items) - 1
	}
}

func (m *menu) selectedItem() menuItem {
	return m.items[m.selectedIndex]
}

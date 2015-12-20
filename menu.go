package main

type menu struct {
	selectedItem menuItem
}

type menuItem int

const (
	menuContinueGame menuItem = iota
	menuNewGame
)

const mainMenuItemCount = 2

func newMenu() *menu {
	return &menu{}
}

func (m *menu) moveDown() {
	m.selectedItem = (m.selectedItem + 1) % mainMenuItemCount
}

func (m *menu) moveUp() {
	m.selectedItem -= 1
	if m.selectedItem < 0 {
		m.selectedItem = mainMenuItemCount - 1
	}
}

package game

type Menu struct {
	Items         []MenuItem
	SelectedIndex int
}

type MenuItem int

const (
	menuContinueGame MenuItem = iota
	menuNewGame
	menuExit
)

var MenuItemText = map[MenuItem]string{
	menuContinueGame: "C O N T I N U E  G A M E",
	menuNewGame:      "N E W  G A M E",
	menuExit:         "E X I T",
}

func newMenu() *Menu {
	return &Menu{
		Items: []MenuItem{
			menuNewGame,
			menuExit,
		},
	}
}

func (m *Menu) addContinueGame() {
	m.Items = []MenuItem{
		menuContinueGame,
		menuNewGame,
		menuExit,
	}
	m.SelectedIndex = 0
}

func (m *Menu) removeContinueGame() {
	m.Items = []MenuItem{
		menuNewGame,
		menuExit,
	}
	m.SelectedIndex = 0
}

func (m *Menu) moveDown() {
	m.SelectedIndex = (m.SelectedIndex + 1) % len(m.Items)
}

func (m *Menu) moveUp() {
	if m.SelectedIndex -= 1; m.SelectedIndex < 0 {
		m.SelectedIndex = len(m.Items) - 1
	}
}

func (m *Menu) selectedItem() MenuItem {
	return m.Items[m.SelectedIndex]
}

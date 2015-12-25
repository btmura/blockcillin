package game

type Menu struct {
	Items        []MenuItem
	FocusedIndex int
	Selected     bool
	Pulse        float32
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
	m.FocusedIndex = 0
	m.Selected = false
}

func (m *Menu) removeContinueGame() {
	m.Items = []MenuItem{
		menuNewGame,
		menuExit,
	}
	m.FocusedIndex = 0
	m.Selected = false
}

func (m *Menu) moveDown() {
	m.FocusedIndex = (m.FocusedIndex + 1) % len(m.Items)
	m.Selected = false
}

func (m *Menu) moveUp() {
	if m.FocusedIndex -= 1; m.FocusedIndex < 0 {
		m.FocusedIndex = len(m.Items) - 1
		m.Selected = false
	}
}

func (m *Menu) focused() MenuItem {
	return m.Items[m.FocusedIndex]
}

func (m *Menu) update() {
	m.Pulse++
}

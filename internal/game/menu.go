package game

type Menu struct {
	Items        []MenuItem
	FocusedIndex int
	Selected     bool
	Pulse        float32
}

type MenuItem int

const (
	MenuItemContinueGame MenuItem = iota
	MenuItemNewGame
	MenuItemExit
)

var MenuItemText = map[MenuItem]string{
	MenuItemContinueGame: "C O N T I N U E  G A M E",
	MenuItemNewGame:      "N E W  G A M E",
	MenuItemExit:         "E X I T",
}

func newMenu() *Menu {
	return &Menu{
		Items: []MenuItem{
			MenuItemNewGame,
			MenuItemExit,
		},
	}
}

func (m *Menu) addContinueGame() {
	m.Items = []MenuItem{
		MenuItemContinueGame,
		MenuItemNewGame,
		MenuItemExit,
	}
	m.FocusedIndex = 0
	m.Selected = false
}

func (m *Menu) removeContinueGame() {
	m.Items = []MenuItem{
		MenuItemNewGame,
		MenuItemExit,
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

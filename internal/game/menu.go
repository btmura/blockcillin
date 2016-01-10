package game

type Menu struct {
	Title        MenuTitle
	Items        []MenuItem
	FocusedIndex int
	Selected     bool
}

//go:generate stringer -type=MenuTitle
type MenuTitle int32

const (
	MenuTitleInitial MenuTitle = iota
	MenuTitlePaused
	MenuTitleGameOver
)

var MenuTitleText = map[MenuTitle]string{
	MenuTitleInitial:  "b l o c k c i l l i n",
	MenuTitlePaused:   "P A U S E D",
	MenuTitleGameOver: "G A M E  O V E R",
}

//go:generate stringer -type=MenuItem
type MenuItem int

const (
	MenuItemContinueGame MenuItem = iota
	MenuItemNewGame
	MenuItemStats
	MenuItemOptions
	MenuItemCredits
	MenuItemExit
)

var MenuItemText = map[MenuItem]string{
	MenuItemContinueGame: "C O N T I N U E  G A M E",
	MenuItemNewGame:      "N E W  G A M E",
	MenuItemStats:        "S T A T S",
	MenuItemOptions:      "O P T I O N S",
	MenuItemCredits:      "C R E D I T S",
	MenuItemExit:         "E X I T",
}

func newMenu() *Menu {
	return &Menu{
		Items: []MenuItem{
			MenuItemNewGame,
			MenuItemStats,
			MenuItemOptions,
			MenuItemCredits,
			MenuItemExit,
		},
	}
}

func (m *Menu) pause() {
	m.Title = MenuTitlePaused
	m.Items = []MenuItem{
		MenuItemContinueGame,
		MenuItemNewGame,
		MenuItemStats,
		MenuItemOptions,
		MenuItemCredits,
		MenuItemExit,
	}
	m.FocusedIndex = 0
	m.Selected = false
}

func (m *Menu) gameOver() {
	m.Title = MenuTitleGameOver
	m.Items = []MenuItem{
		MenuItemNewGame,
		MenuItemStats,
		MenuItemOptions,
		MenuItemCredits,
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

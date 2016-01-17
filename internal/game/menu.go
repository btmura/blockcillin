package game

import "github.com/btmura/blockcillin/internal/audio"

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
	MenuTitleNewGame
	MenuTitlePaused
	MenuTitleGameOver
)

var MenuTitleText = [...]string{
	MenuTitleInitial:  "b l o c k c i l l i n",
	MenuTitleNewGame:  "N E W  G A M E",
	MenuTitlePaused:   "P A U S E D",
	MenuTitleGameOver: "G A M E  O V E R",
}

//go:generate stringer -type=MenuItem
type MenuItem int

const (
	MenuItemNewGame MenuItem = iota
	MenuItemStats
	MenuItemOptions
	MenuItemCredits
	MenuItemExit

	// Pause menu
	MenuItemContinueGame
	MenuItemQuit

	// New game menu
	MenuItemDifficulty
	MenuItemSpeed
	MenuItemOK
)

var MenuItemText = [...]string{
	MenuItemNewGame: "N E W  G A M E",
	MenuItemStats:   "S T A T S",
	MenuItemOptions: "O P T I O N S",
	MenuItemCredits: "C R E D I T S",
	MenuItemExit:    "E X I T",

	// Pause menu
	MenuItemContinueGame: "C O N T I N U E  G A M E",
	MenuItemQuit:         "Q U I T",

	// New game menu
	MenuItemDifficulty: "D I F F I C U L T Y",
	MenuItemSpeed:      "S P E E D",
	MenuItemOK:         "O K",
}

var (
	mainMenu = &Menu{
		Title: MenuTitleInitial,
		Items: []MenuItem{
			MenuItemNewGame,
			MenuItemStats,
			MenuItemOptions,
			MenuItemCredits,
			MenuItemExit,
		},
	}

	newGameMenu = &Menu{
		Title: MenuTitleNewGame,
		Items: []MenuItem{
			MenuItemSpeed,
			MenuItemDifficulty,
			MenuItemOK,
		},
	}

	pauseMenu = &Menu{
		Title: MenuTitlePaused,
		Items: []MenuItem{
			MenuItemContinueGame,
			MenuItemOptions,
			MenuItemQuit,
		},
	}

	gameOverMenu = &Menu{
		Title: MenuTitleGameOver,
		Items: []MenuItem{
			MenuItemQuit,
		},
	}
)

func (m *Menu) reset() {
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

func (m *Menu) selectItem() {
	m.Selected = true
	audio.Play(audio.SoundSelect)
}

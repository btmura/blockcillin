package game

import "github.com/btmura/blockcillin/internal/audio"

type Menu struct {
	ID           MenuID
	Items        []*MenuItem
	FocusedIndex int
	Selected     bool
}

//go:generate stringer -type=MenuID
type MenuID byte

const (
	MenuMain MenuID = iota
	MenuNewGame
	MenuPaused
	MenuGameOver
)

var MenuTitleText = map[MenuID]string{
	MenuMain:     "b l o c k c i l l i n",
	MenuNewGame:  "N E W  G A M E",
	MenuPaused:   "P A U S E D",
	MenuGameOver: "G A M E  O V E R",
}

type MenuItem struct {
	ID       MenuItemID
	Selector *MenuSelector
	Slider   *MenuSlider
}

//go:generate stringer -type=MenuItemID
type MenuItemID byte

const (
	MenuNewGameItem MenuItemID = iota
	MenuStats
	MenuOptions
	MenuCredits
	MenuExit

	MenuSpeed
	MenuDifficulty
	MenuOK

	MenuContinueGame
	MenuQuit
)

var MenuItemText = map[MenuItemID]string{
	MenuNewGameItem: "N E W  G A M E",
	MenuStats:       "S T A T S",
	MenuOptions:     "O P T I O N S",
	MenuCredits:     "C R E D I T S",
	MenuExit:        "E X I T",

	MenuSpeed:      "S P E E D",
	MenuDifficulty: "D I F F I C U L T Y",
	MenuOK:         "O K",

	MenuContinueGame: "C O N T I N U E  G A M E",
	MenuQuit:         "Q U I T",
}

func (i *MenuItem) SingleChoice() bool {
	return i.Selector == nil && i.Slider == nil
}

type MenuSelector struct {
	Choices       []MenuChoiceID
	selectedIndex int
}

//go:generate stringer -type=MenuChoiceID
type MenuChoiceID byte

const (
	MenuEasy MenuChoiceID = iota
	MenuMedium
	MenuHard
)

var MenuChoiceText = map[MenuChoiceID]string{
	MenuEasy:   "E A S Y",
	MenuMedium: "M E D I U M",
	MenuHard:   "H A R D",
}

func (s *MenuSelector) Value() MenuChoiceID {
	return s.Choices[s.selectedIndex]
}

type MenuSlider struct {
	Min   int
	Max   int
	Value int
}

var (
	mainMenu = &Menu{
		ID: MenuMain,
		Items: []*MenuItem{
			{ID: MenuNewGameItem},
			{ID: MenuStats},
			{ID: MenuOptions},
			{ID: MenuCredits},
			{ID: MenuExit},
		},
	}

	speedItem = &MenuItem{
		ID: MenuSpeed,
		Slider: &MenuSlider{
			Min:   1,
			Max:   100,
			Value: 1,
		},
	}

	difficultyItem = &MenuItem{
		ID: MenuDifficulty,
		Selector: &MenuSelector{
			Choices: []MenuChoiceID{
				MenuEasy,
				MenuMedium,
				MenuHard,
			},
		},
	}

	newGameMenu = &Menu{
		ID: MenuNewGame,
		Items: []*MenuItem{
			speedItem,
			difficultyItem,
			{ID: MenuOK},
		},
	}

	pausedMenu = &Menu{
		ID: MenuPaused,
		Items: []*MenuItem{
			{ID: MenuContinueGame},
			{ID: MenuOptions},
			{ID: MenuQuit},
		},
	}

	gameOverMenu = &Menu{
		ID: MenuGameOver,
		Items: []*MenuItem{
			{ID: MenuQuit},
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
	audio.Play(audio.SoundMove)
}

func (m *Menu) moveUp() {
	if m.FocusedIndex -= 1; m.FocusedIndex < 0 {
		m.FocusedIndex = len(m.Items) - 1
		m.Selected = false
	}
	audio.Play(audio.SoundMove)
}

func (m *Menu) moveLeft() {
	item := m.Items[m.FocusedIndex]
	switch {
	case item.Selector != nil:
		if item.Selector.selectedIndex--; item.Selector.selectedIndex < 0 {
			item.Selector.selectedIndex = len(item.Selector.Choices) - 1
		}
		audio.Play(audio.SoundMove)

	case item.Slider != nil:
		if item.Slider.Value--; item.Slider.Value < item.Slider.Min {
			item.Slider.Value = item.Slider.Max
		}
		audio.Play(audio.SoundMove)
	}
}

func (m *Menu) moveRight() {
	item := m.Items[m.FocusedIndex]
	switch {
	case item.Selector != nil:
		item.Selector.selectedIndex = (item.Selector.selectedIndex + 1) % len(item.Selector.Choices)
		audio.Play(audio.SoundMove)

	case item.Slider != nil:
		if item.Slider.Value++; item.Slider.Value > item.Slider.Max {
			item.Slider.Value = item.Slider.Min
		}
		audio.Play(audio.SoundMove)
	}
}

func (m *Menu) focused() MenuItemID {
	return m.Items[m.FocusedIndex].ID
}

func (m *Menu) selectItem() {
	m.Selected = true
	audio.Play(audio.SoundSelect)
}

package game

import "github.com/btmura/blockcillin/internal/audio"

type Menu struct {
	ID           MenuID
	Items        []*MenuItem
	FocusedIndex int
	Selected     bool
}

//go:generate stringer -type=MenuID
type MenuID int32

const (
	MenuIDMain MenuID = iota
	MenuIDNewGame
	MenuIDPaused
	MenuIDGameOver
)

var MenuText = [...]string{
	MenuIDMain:     "b l o c k c i l l i n",
	MenuIDNewGame:  "N E W  G A M E ",
	MenuIDPaused:   "P A U S E D",
	MenuIDGameOver: "G A M E  O V E R ",
}

type MenuItem struct {
	ID   MenuItemID
	Type MenuItemType

	Choices        []MenuChoice
	SelectedChoice int

	MinSliderValue int
	MaxSliderValue int
	SliderValue    int
}

//go:generate stringer -type=MenuItemID
type MenuItemID byte

const (
	MenuItemIDNewGame MenuItemID = iota
	MenuItemIDStats
	MenuItemIDOptions
	MenuItemIDCredits
	MenuItemIDExit

	MenuItemIDSpeed
	MenuItemIDDifficulty
	MenuItemIDOK

	MenuItemIDContinueGame
	MenuItemIDQuit
)

var MenuItemText = [...]string{
	MenuItemIDNewGame: "N E W  G A M E",
	MenuItemIDStats:   "S T A T S",
	MenuItemIDOptions: "O P T I O N S",
	MenuItemIDCredits: "C R E D I T S",
	MenuItemIDExit:    "E X I T",

	MenuItemIDSpeed:      "S P E E D",
	MenuItemIDDifficulty: "D I F F I C U L T Y",
	MenuItemIDOK:         "O K",

	MenuItemIDContinueGame: "C O N T I N U E  G A M E",
	MenuItemIDQuit:         "Q U I T",
}

//go:generate stringer -type=MenuItemType
type MenuItemType byte

const (
	MenuItemTypeChoice MenuItemType = iota
	MenuItemTypeSlider
)

//go:generate stringer -type=MenuChoice
type MenuChoice byte

const (
	MenuChoiceEasy MenuChoice = iota
	MenuChoiceMedium
	MenuChoiceHard
)

var MenuChoiceText = [...]string{
	MenuChoiceEasy:   "E A S Y",
	MenuChoiceMedium: "M E D I U M",
	MenuChoiceHard:   "H A R D",
}

var (
	mainMenu = &Menu{
		ID: MenuIDMain,
		Items: []*MenuItem{
			{ID: MenuItemIDNewGame},
			{ID: MenuItemIDStats},
			{ID: MenuItemIDOptions},
			{ID: MenuItemIDCredits},
			{ID: MenuItemIDExit},
		},
	}

	speedItem = &MenuItem{
		ID:             MenuItemIDSpeed,
		Type:           MenuItemTypeSlider,
		MinSliderValue: 1,
		MaxSliderValue: 99,
		SliderValue:    1,
	}

	difficultyItem = &MenuItem{
		ID:   MenuItemIDDifficulty,
		Type: MenuItemTypeChoice,
		Choices: []MenuChoice{
			MenuChoiceEasy,
			MenuChoiceMedium,
			MenuChoiceHard,
		},
	}

	newGameMenu = &Menu{
		ID: MenuIDNewGame,
		Items: []*MenuItem{
			speedItem,
			difficultyItem,
			{ID: MenuItemIDOK},
		},
	}

	pausedMenu = &Menu{
		ID: MenuIDPaused,
		Items: []*MenuItem{
			{ID: MenuItemIDContinueGame},
			{ID: MenuItemIDOptions},
			{ID: MenuItemIDQuit},
		},
	}

	gameOverMenu = &Menu{
		ID: MenuIDGameOver,
		Items: []*MenuItem{
			{ID: MenuItemIDQuit},
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
	switch item.Type {
	case MenuItemTypeChoice:
		if len(item.Choices) == 0 {
			return // nothing to choose
		}
		if item.SelectedChoice -= 1; item.SelectedChoice < 0 {
			item.SelectedChoice = len(item.Choices) - 1
		}
		audio.Play(audio.SoundMove)

	case MenuItemTypeSlider:
		if item.SliderValue--; item.SliderValue < item.MinSliderValue {
			item.SliderValue = item.MinSliderValue
		}
		audio.Play(audio.SoundMove)
	}
}

func (m *Menu) moveRight() {
	item := m.Items[m.FocusedIndex]
	switch item.Type {
	case MenuItemTypeChoice:
		if len(item.Choices) == 0 {
			return // nothing to choose
		}
		item.SelectedChoice = (item.SelectedChoice + 1) % len(item.Choices)
		audio.Play(audio.SoundMove)

	case MenuItemTypeSlider:
		if item.SliderValue++; item.SliderValue > item.MaxSliderValue {
			item.SliderValue = item.MaxSliderValue
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

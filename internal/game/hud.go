package game

type HUD struct {
	Speed   int
	TimeSec int
	Score   int

	timeUpdates int
}

func newHUD(speed int) *HUD {
	return &HUD{
		Speed: speed,
	}
}

//go:generate stringer -type=HUDItem
type HUDItem int32

const (
	HUDItemSpeed HUDItem = iota
	HUDItemTime
	HUDItemScore
)

var HUDItemText = [...]string{
	HUDItemSpeed: "S P E E D",
	HUDItemTime:  "T I M E",
	HUDItemScore: "S C O R E",
}

func (h *HUD) update() {
	if h.timeUpdates++; h.timeUpdates == updatesPerSec {
		h.timeUpdates = 0
		h.TimeSec++
	}
}

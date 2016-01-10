package game

type HUD struct {
	Speed   int
	TimeSec int
	Score   int

	timeUpdates int
}

func newHUD() *HUD {
	return &HUD{
		Speed: 1,
	}
}

//go:generate stringer -type=HUDItem
type HUDItem int32

const (
	HUDItemSpeed HUDItem = iota
	HUDItemTime
	HUDItemScore
)

var HUDItemText = [3]string{
	"S P E E D",
	"T I M E",
	"S C O R E",
}

func (h *HUD) update() {
	if h.timeUpdates++; h.timeUpdates == updatesPerSec {
		h.timeUpdates = 0
		h.TimeSec++
	}
}

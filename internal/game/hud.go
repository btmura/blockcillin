package game

type HUD struct {
	Speed   int
	TimeSec int
	Score   int
}

func newHUD() *HUD {
	return &HUD{
		Speed: 1,
	}
}

type HUDItem int

const (
	HUDItemSpeed HUDItem = iota
	HUDItemTime
	HUDItemScore
)

var HUDItemText = map[HUDItem]string{
	HUDItemSpeed: "S P E E D",
	HUDItemTime:  "T I M E",
	HUDItemScore: "S C O R E",
}

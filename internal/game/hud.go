package game

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

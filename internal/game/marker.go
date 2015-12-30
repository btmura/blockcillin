package game

type Marker struct {
	State      MarkerState
	ComboLevel int
	ChainLevel int
	step       float32
}

//go:generate stringer -type=MarkerState
type MarkerState int32

const (
	// MarkerNone is an empty marker.
	MarkerNone MarkerState = iota

	// MarkeShowing is a marker being shown.
	MarkerShowing
)

var markerStateSteps = map[MarkerState]float32{
	MarkerShowing: 2.0 / SecPerUpdate,
}

func (m *Marker) show(comboLevel, chainLevel int) {
	m.ComboLevel = comboLevel
	m.ChainLevel = chainLevel
	m.setState(MarkerShowing)
}

func (m *Marker) update() {
	advance := func(nextState MarkerState) {
		if m.step++; m.step >= markerStateSteps[m.State] {
			m.setState(nextState)
		}
	}

	switch m.State {
	case MarkerShowing:
		advance(MarkerNone)
	}
}

func (m *Marker) StateProgress(fudge float32) float32 {
	totalSteps := markerStateSteps[m.State]
	if totalSteps == 0 {
		return 1
	}

	if p := (m.step + fudge) / totalSteps; p < 1 {
		return p
	}
	return 1
}

func (m *Marker) setState(state MarkerState) {
	m.State = state
	m.step = 0
}

package game

import "sort"

type match struct {
	color BlockColor
	cells []*matchCell
}

type matchCell struct {
	x int
	y int
}

func findMatches(b *Board) []*match {
	var matches []*match

	hm := findHorizontalMatches(b)
	vm := findVerticalMatches(b)

	// Combine intersecting horizontal and vertical matches.
	for len(hm) > 0 {
		// Pop the first match off as the candidate.
		m := hm[0]
		hm = hm[1:]

		// Add the candidate to the final list.
		matches = append(matches, m)

		// Keep appending intersections to the candidate.
		needSort := false
		for {
			// Check for new vertical intersections.
			intersected := false
			for i := 0; i < len(vm); i++ {
				if intersects(m, vm[i]) {
					// Add all cells except the one that intersected.
					for _, c := range vm[i].cells {
						if !contains(m, c) {
							m.cells = append(m.cells, c)
						}
					}
					// Remove the match since it's now part of the candidate.
					vm = append(vm[:i], vm[i+1:]...)
					intersected = true
					needSort = true
					i--
				}
			}

			// Break if no intersections. Candidate can't grow larger.
			if !intersected {
				break
			}

			// Check for new horizontal intersections.
			intersected = false
			for i := 0; i < len(hm); i++ {
				if intersects(m, hm[i]) {
					// Add all cells except the one that intersected.
					for _, c := range hm[i].cells {
						if !contains(m, c) {
							m.cells = append(m.cells, c)
						}
					}
				}
				// Remove the match since it's now part of the candidate.
				hm = append(hm[:i], hm[i+1:]...)
				intersected = true
				needSort = true
				i--
			}

			// Break if no intersections. Candidate can't grow any larger.
			if !intersected {
				break
			}
		}

		// Sort the match's cells by row so they disappear row by row.
		// Use stable sort to preserve horizontal ordering due to wrapping.
		if needSort {
			sort.Stable(byRowAndIndex(m.cells))
		}
	}

	// Add any vertical matches that never intersected.
	for _, m := range vm {
		matches = append(matches, m)
	}

	return matches
}

func findHorizontalMatches(b *Board) []*match {
	var matches []*match

	for y, r := range b.Rings {
		// Find the initial position where the colors change
		// to handle a matching chain that wraps around.

		var bc BlockColor
		var initX int
	init:
		for x, c := range r.Cells {
			switch {
			case x == 0:
				bc = c.Block.Color

			case bc != c.Block.Color:
				break init
			}
			initX++
		}

		// Now find matching matches starting from the initial position
		// which may require wrapping around.

		var startX int
		var numMatches int

		startChain := func(x int, c *Cell) {
			bc = c.Block.Color
			startX = x
			numMatches = 1
		}

		continueChain := func() {
			numMatches++
		}

		endChain := func() {
			if numMatches >= 3 {
				m := &match{color: bc}
				for i := 0; i < numMatches; i++ {
					x := (startX + i) % b.CellCount
					m.cells = append(m.cells, &matchCell{x, y})
				}
				matches = append(matches, m)
			}
			numMatches = 0
		}

		for i := 0; i < b.CellCount; i++ {
			x := (initX + i) % b.CellCount
			c := r.Cells[x]
			switch {
			case c.Block.State != BlockStatic:
				endChain()

			case numMatches == 0:
				startChain(x, c)

			case bc == c.Block.Color:
				continueChain()

			case bc != c.Block.Color:
				endChain()
				startChain(x, c)
			}
		}

		endChain()
	}

	return matches
}

func findVerticalMatches(b *Board) []*match {
	var matches []*match

	for x := 0; x < b.CellCount; x++ {
		var bc BlockColor
		var startY int
		var numMatches int

		startChain := func(y int, c *Cell) {
			bc = c.Block.Color
			startY = y
			numMatches = 1
		}

		continueChain := func() {
			numMatches++
		}

		endChain := func() {
			if numMatches >= 3 {
				m := &match{color: bc}
				for i := 0; i < numMatches; i++ {
					y := startY + i
					m.cells = append(m.cells, &matchCell{x, y})
				}
				matches = append(matches, m)
			}
			numMatches = 0
		}

		for y, r := range b.Rings {
			c := r.Cells[x]
			switch {
			case c.Block.State != BlockStatic:
				endChain()

			case numMatches == 0:
				startChain(y, c)

			case bc == c.Block.Color:
				continueChain()

			case bc != c.Block.Color:
				endChain()
				startChain(y, c)
			}
		}

		endChain()
	}
	return matches
}

func intersects(m1, m2 *match) bool {
	if m1.color != m2.color {
		return false
	}

	for _, c1 := range m1.cells {
		for _, c2 := range m2.cells {
			if c1.x == c2.x && c1.y == c2.y {
				return true
			}
		}
	}
	return false
}

func contains(m *match, mc *matchCell) bool {
	for _, c := range m.cells {
		if c.x == mc.x && c.y == mc.y {
			return true
		}
	}
	return false
}

// byRowAndIndex is a matchCell slice that can be sorted.
type byRowAndIndex []*matchCell

// Len implements sort.Interface
func (c byRowAndIndex) Len() int {
	return len(c)
}

// Less implements sort.Interface
func (c byRowAndIndex) Less(i, j int) bool {
	ydiff := c[i].y - c[j].y
	if ydiff != 0 {
		return ydiff < 0
	}
	// Horizontal chain cells may wrap around the cylinder's seam.
	// So leave them in the order they were inserted.
	return i < j
}

// Swap implements sort.Interface
func (c byRowAndIndex) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

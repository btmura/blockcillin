package main

import "sort"

type chain struct {
	color blockColor
	cells []*chainCell
}

type chainCell struct {
	x int
	y int
}

func findChains(b *board) []*chain {
	var chains []*chain

	hc := findHorizontalChains(b)
	vc := findVerticalChains(b)

	// Combine intersecting horizontal and vertical chains.
	for len(hc) > 0 {
		// Pop the first chain off and make it the candidate.
		ch := hc[0]
		hc = hc[1:]

		// Add the candidate to the final list.
		chains = append(chains, ch)

		// Keep appending intersections to the candidate.
		for {
			// Check for new vertical intersections.
			intersected := false
			for i := 0; i < len(vc); i++ {
				if intersects(ch, vc[i]) {
					// Add all cells except the one that intersected.
					for _, c := range vc[i].cells {
						if !contains(ch, c) {
							ch.cells = append(ch.cells, c)
						}
					}
					// Remove the chain since it's now part of the candidate.
					vc = append(vc[:i], vc[i+1:]...)
					intersected = true
				}
			}

			// Break if no intersections. Candidate can't grow larger.
			if !intersected {
				break
			}

			// Check for new horizontal intersections.
			intersected = false
			for i := 0; i < len(hc); i++ {
				if intersects(ch, hc[i]) {
					// Add all cells except the one that intersected.
					for _, c := range hc[i].cells {
						if !contains(ch, c) {
							ch.cells = append(ch.cells, c)
						}
					}
				}
				// Remove the chain since it's now part of the candidate.
				hc = append(hc[:i], hc[i+1:]...)
				intersected = true
			}

			// Break if no intersections. Candidate can't grow any larger.
			if !intersected {
				break
			}
		}
	}

	// Add any vertical chains that never intersected.
	for _, ch := range vc {
		chains = append(chains, ch)
	}

	// Sort the cells within each chain by row so they disappear orderly.
	for _, ch := range chains {
		sort.Sort(byRowAndIndex(ch.cells))
	}

	return chains
}

func findHorizontalChains(b *board) []*chain {
	var chains []*chain

	for y, r := range b.rings {
		// Find the initial position where the colors change
		// to handle a matching chain that wraps around.

		var cc blockColor
		var initX int
	init:
		for x, c := range r.cells {
			switch {
			case x == 0:
				cc = c.block.color

			case cc != c.block.color:
				break init
			}
			initX++
		}

		// Now find matching chains starting from the initial position
		// which may require wrapping around.

		var startX int
		var numMatches int

		startChain := func(x int, c *cell) {
			cc = c.block.color
			startX = x
			numMatches = 1
		}

		continueChain := func() {
			numMatches++
		}

		endChain := func() {
			if numMatches >= 3 {
				ch := &chain{color: cc}
				for i := 0; i < numMatches; i++ {
					x := (startX + i) % b.cellCount
					ch.cells = append(ch.cells, &chainCell{x, y})
				}
				chains = append(chains, ch)
			}
			numMatches = 0
		}

		for i := 0; i < b.cellCount; i++ {
			x := (initX + i) % b.cellCount
			c := r.cells[x]
			switch {
			case !c.block.isClearable():
				endChain()

			case numMatches == 0:
				startChain(x, c)

			case cc == c.block.color:
				continueChain()

			case cc != c.block.color:
				endChain()
				startChain(x, c)
			}
		}

		endChain()
	}

	return chains
}

func findVerticalChains(b *board) []*chain {
	var chains []*chain

	for x := 0; x < b.cellCount; x++ {
		var cc blockColor
		var startY int
		var numMatches int

		startChain := func(y int, c *cell) {
			cc = c.block.color
			startY = y
			numMatches = 1
		}

		continueChain := func() {
			numMatches++
		}

		endChain := func() {
			if numMatches >= 3 {
				ch := &chain{color: cc}
				for i := 0; i < numMatches; i++ {
					y := startY + i
					ch.cells = append(ch.cells, &chainCell{x, y})
				}
				chains = append(chains, ch)
			}
			numMatches = 0
		}

		for y, r := range b.rings {
			c := r.cells[x]
			switch {
			case !c.block.isClearable():
				endChain()

			case numMatches == 0:
				startChain(y, c)

			case cc == c.block.color:
				continueChain()

			case cc != c.block.color:
				endChain()
				startChain(y, c)
			}
		}

		endChain()
	}
	return chains
}

func intersects(ch1, ch2 *chain) bool {
	if ch1.color != ch2.color {
		return false
	}

	for _, c1 := range ch1.cells {
		for _, c2 := range ch2.cells {
			if c1.x == c2.x && c1.y == c2.y {
				return true
			}
		}
	}
	return false
}

func contains(ch *chain, cc *chainCell) bool {
	for _, c := range ch.cells {
		if c.x == cc.x && c.y == cc.y {
			return true
		}
	}
	return false
}

// byRowAndIndex is a chainCell slice that can be sorted.
type byRowAndIndex []*chainCell

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
	return i-j < 0
}

// Swap implements sort.Interface
func (c byRowAndIndex) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

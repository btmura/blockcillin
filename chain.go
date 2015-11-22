package main

type chain struct {
	cells []*chainCell
}

type chainCell struct {
	x int
	y int
}

func findChains(b *board) []*chain {
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
				ch := &chain{}
				for i := 0; i < numMatches; i++ {
					x := (startX + i) % b.cellCount
					ch.cells = append(ch.cells, &chainCell{x, y})
				}
				chains = append(chains, ch)
			}
		}

		for i := 0; i < b.cellCount; i++ {
			x := (initX + i) % b.cellCount
			c := r.cells[x]
			switch {
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

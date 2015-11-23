package main

type chain struct {
	cells []*chainCell
}

type chainCell struct {
	x int
	y int
}

func findChains(b *board) []*chain {
	return findVerticalChains(b)
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
				ch := &chain{}
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
				ch := &chain{}
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

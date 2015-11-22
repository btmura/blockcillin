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
		var cc blockColor
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
					ch.cells = append(ch.cells, &chainCell{startX + i, y})
				}
				chains = append(chains, ch)
			}
		}

		for x, c := range r.cells {
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

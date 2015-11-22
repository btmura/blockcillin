package main

type drop struct {
	x int
	y int
}

func findDrops(b *board) []*drop {
	var drops []*drop

	for y, r := range b.rings {
		for x, c := range r.cells {
			if !c.block.isDroppable() {
				continue
			}

			if y+1 == b.ringCount {
				continue
			}

			nr := b.rings[y+1]
			nc := nr.cells[x]

			if !nc.block.canReceiveDrop() {
				continue
			}

			drops = append(drops, &drop{x, y})
		}
	}

	return drops
}

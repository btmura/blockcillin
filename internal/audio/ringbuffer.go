package audio

import (
	"fmt"
	"sync"
)

type ringBuffer struct {
	sync.Mutex
	data  []int16
	start int
	end   int
	count int
}

func newRingBuffer(size int) *ringBuffer {
	return &ringBuffer{
		data: make([]int16, size),
	}
}

// push adds the values from values to the buffer.
// Oldest values are overwritten if the buffer is full.
func (b *ringBuffer) push(values ...int16) {
	b.Lock()
	defer b.Unlock()

	for _, v := range values {
		// Remove the first element to make room if we're full.
		if b.count == len(b.data) {
			b.start = (b.start + 1) % len(b.data)
			b.count--
		}

		// Add new value to the end.
		b.data[b.end] = v
		b.end = (b.end + 1) % len(b.data)
		b.count++
	}
}

// pop removes the values from the buffer and puts them into values slice.
func (b *ringBuffer) pop(values []int16) {
	b.Lock()
	defer b.Unlock()

	for i := 0; i < len(values); i++ {
		// Set value of zero if we're empty.
		if b.count == 0 {
			values[i] = 0
			continue
		}

		// Remove the first element.
		values[i] = b.data[b.start]
		b.start = (b.start + 1) % len(b.data)
		b.count--
	}
}

func (b *ringBuffer) String() string {
	return fmt.Sprintf("start: %d end: %d count: %d", b.start, b.end, b.count)
}

package lib

import "sync"

type RingBuffer[T any] struct {
	l sync.RWMutex

	buffer   []*T
	readPos  uint64
	writePos uint64
	size     uint64
}

func NewRingBuffer[T any](size uint64) *RingBuffer[T] {
	return &RingBuffer[T]{
		l:        sync.RWMutex{},
		buffer:   make([]*T, size),
		readPos:  0,
		writePos: 0,
		size:     size,
	}
}

func (rb *RingBuffer[T]) Enqueue(item *T) bool {
	if item == nil {
		return false
	}

	rb.l.Lock()
	defer rb.l.Unlock()

	nextPos := (rb.writePos + 1) % rb.size
	if nextPos == rb.readPos {
		return false // buffer is full
	}

	rb.buffer[rb.writePos] = item
	rb.writePos = nextPos
	return true
}

func (rb *RingBuffer[T]) Dequeue() (*T, bool) {
	rb.l.RLock()
	defer rb.l.RUnlock()

	if rb.readPos == rb.writePos {
		return nil, false // buffer is empty
	}
	item := rb.buffer[rb.readPos]
	rb.buffer[rb.readPos] = nil
	rb.readPos = (rb.readPos + 1) % rb.size
	return item, true
}

func (rb *RingBuffer[T]) Length() uint64 {
	rb.l.RLock()
	defer rb.l.RUnlock()

	return (rb.writePos - rb.readPos + rb.size) % rb.size
}

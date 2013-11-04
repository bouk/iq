package priority_iq

import (
	"container/heap"
)

type Object struct {
	Priority int
	Value    string
}

type ObjectHeap []Object

func (h ObjectHeap) Len() int {
	return len(h)
}

// It's a min heap so we do a gt to make the highest priority pop first
func (h ObjectHeap) Less(i, j int) bool {
	return h[i].Priority > h[j].Priority
}

func (h ObjectHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *ObjectHeap) Push(x interface{}) {
	*h = append(*h, x.(Object))
}

func (h *ObjectHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// HeapIQ takes in objects with a priority and gives back just the objects ordered by priority
func HeapIQ(in <-chan Object, next chan<- string) {
	defer close(next)

	// pending events (this is the "infinite" part)
	pending := &ObjectHeap{}

recv:
	for {
		// Ensure that pending always has values so the select can
		// multiplex between the receiver and sender properly
		if pending.Len() == 0 {
			v, ok := <-in
			if !ok {
				// in is closed, flush values
				break
			}

			// We now have something to send
			heap.Push(pending, v)
		}

		select {
		// Queue incoming values
		case v, ok := <-in:
			if !ok {
				// in is closed, flush values
				break recv
			}
			heap.Push(pending, v)

		// Send queued values
		case next <- heap.Pop(pending).(Object).Value:
		}
	}

	// After in is closed, we may still have events to send
	for pending.Len() != 0 {
		next <- heap.Pop(pending).(Object).Value
	}
}

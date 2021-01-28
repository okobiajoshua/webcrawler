package util

// FIFO queue
type FIFO struct {
	queue []string
}

// NewFifo creates new FIFO and returns it
func NewFifo() *FIFO {
	return &FIFO{
		queue: make([]string, 0),
	}
}

// Push pushed node to the back of the queue
func (f *FIFO) Push(node string) {
	f.queue = append(f.queue, node)
}

// Length of queue
func (f *FIFO) Length() int {
	return len(f.queue)
}

// IsEmpty return true if queue is empty else, it returns false
func (f *FIFO) IsEmpty() bool {
	return len(f.queue) == 0
}

// Front takes a value from the front of the queue and returns it
func (f *FIFO) Front() string {
	if len(f.queue) == 0 {
		return ""
	}

	node := f.queue[0]
	f.queue = f.queue[1:]

	return node
}

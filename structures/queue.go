package structures

// Queue is a convenient data structure used for FIFO
type Queue[T any] interface {
	Enqueue(T)
	Dequeue() (T, bool)
	Len() int
}

func NewQueue[T any]() Queue[T] {
	return &queue[T]{
		l: &LinkedList[T]{},
	}
}

type queue[T any] struct {
	l *LinkedList[T]
}

func (q *queue[T]) Enqueue(value T) {
	q.l.Append(value)
}

func (q *queue[T]) Dequeue() (T, bool) {
	return q.l.Shift()
}

func (q *queue[T]) Len() int {
	return q.l.length
}

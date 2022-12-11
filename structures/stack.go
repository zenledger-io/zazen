package structures

// Stack is a convenient data structure used for LIFO
type Stack[T any] interface {
	Push(T)
	Pop() (T, bool)
	Len() int
}

func NewStack[T any]() Stack[T] {
	return &stack[T]{
		l: &LinkedList[T]{},
	}
}

type stack[T any] struct {
	l *LinkedList[T]
}

func (s *stack[T]) Push(value T) {
	s.l.Unshift(value)
}

func (s *stack[T]) Pop() (T, bool) {
	return s.l.Shift()
}

func (s *stack[T]) Len() int {
	return s.l.Len()
}

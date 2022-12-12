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

// NewStackFromSlice accepts a Slice and returns a Stack with
// the first item in the slice at the top of the stack
func NewStackFromSlice[T any](slc []T) Stack[T] {
	s := NewStack[T]()
	for i := len(slc) - 1; i >= 0; i-- {
		s.Push(slc[i])
	}
	return s
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

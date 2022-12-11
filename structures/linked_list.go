package structures

type LinkedList[T any] struct {
	head   *LinkedListNode[T]
	tail   *LinkedListNode[T]
	length int
}

func (l *LinkedList[T]) Append(value T) {
	n := &LinkedListNode[T]{
		Value: value,
	}
	if l.head == nil {
		l.head = n
	} else {
		l.tail.Next = n
	}
	l.tail = n
	l.length += 1
}

func (l *LinkedList[T]) Unshift(value T) {
	n := &LinkedListNode[T]{
		Value: value,
	}
	if l.head == nil {
		l.tail = n
	} else {
		n.Next = l.head
	}
	l.head = n
	l.length += 1
}

func (l *LinkedList[T]) Shift() (T, bool) {
	if l.head == nil {
		var noop T
		return noop, false
	}

	n := l.head
	l.head = n.Next
	n.Next = nil
	l.length -= 1
	return n.Value, true
}

func (l *LinkedList[T]) Len() int {
	return l.length
}

func (l *LinkedList[T]) ForEach(f func(T)) {
	n := l.head
	for n != nil {
		f(n.Value)
		n = n.Next
	}
}

func (l *LinkedList[T]) Copy() *LinkedList[T] {
	cpy := &LinkedList[T]{}
	l.ForEach(func(value T) {
		cpy.Append(value)
	})
	return cpy
}

func (l *LinkedList[T]) ReversedCopy() *LinkedList[T] {
	cpy := &LinkedList[T]{}
	l.ForEach(func(value T) {
		cpy.Unshift(value)
	})
	return cpy
}

func (l *LinkedList[T]) Reverse() {
	var prev *LinkedListNode[T]
	cur := l.head
	l.head = l.tail
	l.tail = cur
	for cur != nil {
		nxt := cur.Next
		cur.Next = prev
		prev = cur
		cur = nxt
	}
}

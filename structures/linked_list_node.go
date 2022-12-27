package structures

type LinkedListNode[T any] struct {
	Value T
	Next  *LinkedListNode[T]
}

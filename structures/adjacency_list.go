package structures

import (
	"golang.org/x/exp/constraints"
)

type AdjacencyList[T comparable, O constraints.Ordered] interface {
	AddEdge(n1, n2 T, weight O)
	AddUndirectedEdge(n1, n2 T, weight O)
	// Edges returns all edges for a given node, return true to stop
	Edges(node T, fn func(node T, weight O) bool)
	Len() int
}

func NewAdjacencyList[T comparable, O constraints.Ordered]() AdjacencyList[T, O] {
	return &adjacencyList[T, O]{
		m: make(map[T][]adjacencyListEdge[T, O], 0),
	}
}

type adjacencyList[T comparable, O constraints.Ordered] struct {
	m map[T][]adjacencyListEdge[T, O]
}

func (l *adjacencyList[T, O]) AddUndirectedEdge(n1, n2 T, weight O) {
	l.AddEdge(n1, n2, weight)
	l.AddEdge(n2, n1, weight)
}

func (l *adjacencyList[T, O]) AddEdge(n1, n2 T, weight O) {
	if len(l.m[n1]) == 0 {
		l.m[n1] = make([]adjacencyListEdge[T, O], 0)
	}

	l.m[n1] = append(l.m[n1], adjacencyListEdge[T, O]{
		Node:   n2,
		Weight: weight,
	})
}

// Edges returns if the passed function returns true
func (l *adjacencyList[T, O]) Edges(t T, fn func(T, O) bool) {
	for _, n := range l.m[t] {
		if fn(n.Node, n.Weight) {
			return
		}
	}
}

func (l *adjacencyList[T, O]) Len() int {
	return len(l.m)
}

// adjacencyListEdge holds a node and weight
type adjacencyListEdge[T comparable, O constraints.Ordered] struct {
	Node   T
	Weight O
}

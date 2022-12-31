package algorithms

import (
	"github.com/zenledger-io/zazen/structures"
	"golang.org/x/exp/constraints"
)

func Dijkstra[T comparable, O constraints.Ordered](al structures.AdjacencyList[T, O], origin T) map[T]O {
	visited := make(map[T]bool)
	distances := make(map[T]O)
	var zero O
	distances[origin] = zero
	findMinV := func() (T, bool) {
		var setVal bool
		var minv T
		for v, d := range distances {
			if visited[v] {
				continue
			}

			if setVal && distances[minv] < d {
				continue
			}

			minv = v
			setVal = true
		}

		return minv, setVal
	}

	minv, ok := findMinV()
	for ok {
		visited[minv] = true
		al.Edges(minv, func(v T, w O) {
			if visited[v] {
				return
			}

			d := distances[minv] + w
			if _, ok := distances[v]; !ok || distances[v] > d {
				distances[v] = d
			}
		})
		minv, ok = findMinV()
	}

	return distances
}

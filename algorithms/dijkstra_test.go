package algorithms

import (
	"github.com/stretchr/testify/require"
	"github.com/zenledger-io/zazen/structures"
	"math/rand"
	"testing"
)

func TestDijkstra(t *testing.T) {
	tcs := map[string]struct {
		AddEdges func(list structures.AdjacencyList[string, int])
		Origin   string
		Expected map[string]int
	}{
		"3 nodes cyclic": {
			AddEdges: func(list structures.AdjacencyList[string, int]) {
				list.AddEdge("B", "A", 1)
				list.AddEdge("A", "B", 3)

				list.AddEdge("A", "C", 1)
				list.AddEdge("C", "A", 2)

				list.AddEdge("C", "B", 1)
				list.AddEdge("B", "C", 5)
			},
			Origin: "A",
			Expected: map[string]int{
				"A": 0,
				"B": 2,
				"C": 1,
			},
		},
		"origin does not exist": {
			AddEdges: func(list structures.AdjacencyList[string, int]) {
				list.AddEdge("B", "A", 1)
				list.AddEdge("A", "B", 3)

				list.AddEdge("A", "C", 1)
				list.AddEdge("C", "A", 2)

				list.AddEdge("C", "B", 1)
				list.AddEdge("B", "C", 5)
			},
			Origin:   "Z",
			Expected: map[string]int{"Z": 0},
		},
		"5 nodes": {
			AddEdges: func(list structures.AdjacencyList[string, int]) {
				list.AddEdge("B", "A", 1)
				list.AddEdge("A", "B", 3)

				list.AddEdge("A", "C", 1)
				list.AddEdge("C", "A", 2)

				list.AddEdge("C", "B", 1)
				list.AddEdge("B", "C", 5)

				list.AddEdge("C", "D", 2)

				list.AddEdge("B", "E", 6)

				list.AddEdge("D", "E", 6)
			},
			Origin: "A",
			Expected: map[string]int{
				"A": 0,
				"B": 2,
				"C": 1,
				"D": 3,
				"E": 8,
			},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			al := structures.NewAdjacencyList[string, int]()
			tc.AddEdges(al)
			distances := Dijkstra(al, tc.Origin)
			require.Equal(t, len(tc.Expected), len(distances))
			for k, v := range tc.Expected {
				require.Equal(t, v, distances[k])
			}
		})
	}
}

func BenchmarkDijkstra(b *testing.B) {
	al := structures.NewAdjacencyList[int, int]()
	for i := 0; i < 10_000; i++ {
		al.AddUndirectedEdge(i, i+1, rand.Int())
		if i >= 2 {
			al.AddUndirectedEdge(i-2, i+1, rand.Int())
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Dijkstra(al, al.Len()/2)
	}
}

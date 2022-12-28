package structures

// BinaryHeap is a data structure that could speed up hifo calculations
type BinaryHeap[T any] interface {
	Len() int
	Push(T)
	Pop() T
}

func NewBinaryHeap[T any](compFn func(T, T) int) BinaryHeap[T] {
	return &binaryHeap[T]{
		compFn: compFn,
	}
}

type binaryHeap[T any] struct {
	slc    []T
	compFn func(T, T) int
}

func (bh *binaryHeap[T]) Len() int {
	return len(bh.slc)
}

func (bh *binaryHeap[T]) Push(t T) {
	bh.slc = append(bh.slc, t)
	bh.bubbleUp(bh.Len() - 1)
}

func (bh *binaryHeap[T]) Pop() T {
	bh.swap(0, bh.Len()-1)
	t := bh.slc[bh.Len()-1]
	bh.slc = bh.slc[:bh.Len()-1]
	bh.bubbleDown(0)
	return t
}

func (bh *binaryHeap[T]) bubbleUp(i int) {
	for i > 0 {
		pi := i - 1
		if i%2 == 0 {
			pi -= 1
		}
		pi = pi / 2
		if bh.compFn(bh.slc[pi], bh.slc[i]) >= 0 {
			return
		}

		bh.swap(i, pi)
		i = pi
	}
}

func (bh *binaryHeap[T]) bubbleDown(i int) {
	for i < bh.Len()-1 {
		li := 2*i + 1
		if li >= bh.Len() {
			return
		}

		ri := 2*i + 2
		ci := li
		if ri < bh.Len() && bh.compFn(bh.slc[ri], bh.slc[li]) > 0 {
			ci = ri
		}
		if bh.compFn(bh.slc[i], bh.slc[ci]) >= 0 {
			return
		}

		bh.swap(i, ci)
		i = ci
	}
}

func (bh *binaryHeap[T]) swap(i, j int) {
	bh.slc[i], bh.slc[j] = bh.slc[j], bh.slc[i]
}

package sortx

// Quick is good for random data but do not use if the slice is almost sorted
func Quick[T any](slc []T, isBefore func(T, T) bool) {
	quick(slc, 0, len(slc)-1, isBefore)
}

func quick[T any](slc []T, left, right int, isBefore func(T, T) bool) {
	if left >= right {
		return
	}

	pi := pivot(slc, left, right, isBefore)
	quick(slc, left, pi-1, isBefore)
	quick(slc, pi+1, right, isBefore)
}

func pivot[T any](slc []T, left, right int, isBefore func(T, T) bool) int {
	swap := left
	for i := left + 1; i <= right; i++ {
		if !isBefore(slc[i], slc[left]) {
			continue
		}

		swap += 1
		slc[swap], slc[i] = slc[i], slc[swap]
	}
	slc[swap], slc[left] = slc[left], slc[swap]
	return swap
}

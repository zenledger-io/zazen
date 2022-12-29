package sortx

// Insertion sort is useful if the slice is mostly sorted
func Insertion[T any](slc []T, isBefore func(T, T) bool) {
	for i := 1; i < len(slc); i++ {
		for j := i; j > 0 && isBefore(slc[j], slc[j-1]); j-- {
			slc[j-1], slc[j] = slc[j], slc[j-1]
		}
	}
}

// Insert is a convenience method that appends an item to a slice that
// should already be sorted and then puts the new item in the correct place
// O(n) if the array is already sorted
func Insert[T any](slc []T, t T, isBefore func(T, T) bool) []T {
	slc = append(slc, t)
	Insertion(slc, isBefore)
	return slc
}

// InsertInSortedSlice inserts an item only into a slice that has been sorted
func InsertInSortedSlice[T any](slc []T, t T, isBefore func(T, T) bool) []T {
	slc = append(slc, t)
	for i := len(slc) - 1; i > 0; i-- {
		if isBefore(slc[i], slc[i-1]) {
			slc[i], slc[i-1] = slc[i-1], slc[i]
		} else {
			return slc
		}
	}
	return slc
}

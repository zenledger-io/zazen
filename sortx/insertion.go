package sortx

// Insertion sort is useful if the slice is mostly sorted
func Insertion[T any](slc []T, isBefore func(T, T) bool) {
	for i := 1; i < len(slc); i++ {
		temp := slc[i]
		j := i - 1
		for j >= 0 && isBefore(temp, slc[j]) {
			slc[j+1] = slc[j]
			j -= 1
		}
		slc[j+1] = temp
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

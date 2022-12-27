package sortx

func Merge[T any](slc []T, isBefore func(T, T) bool) []T {
	if len(slc) <= 1 {
		return slc
	}
	mid := len(slc) / 2
	return merge(Merge(slc[:mid], isBefore), Merge(slc[mid:], isBefore), isBefore)
}

func merge[T any](slc1, slc2 []T, isBefore func(T, T) bool) []T {
	slc := make([]T, len(slc1)+len(slc2))
	var i, j int
	for i < len(slc1) && j < len(slc2) {
		if isBefore(slc1[i], slc2[j]) {
			slc[i+j] = slc1[i]
			i += 1
		} else {
			slc[i+j] = slc2[j]
			j += 1
		}
	}
	for i < len(slc1) {
		slc[i+j] = slc1[i]
		i += 1
	}
	for j < len(slc2) {
		slc[i+j] = slc2[j]
		j += 1
	}
	return slc
}

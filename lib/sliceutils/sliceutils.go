package sliceutils

func FlattenSlice[T any](slice [][]T) []T {
	var result = make([]T, 0)
	for _, v := range slice {
		result = append(result, v...)
	}
	return result
}
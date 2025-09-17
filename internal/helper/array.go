package helper

func Filter[T any](callback func(T) bool, arr []T) []T {
	for i, v := range arr {
		res := callback(v)
		if res {
			arr = append(arr[:i], arr[i+1:]...)
		}
	}

	return arr
}

func Map[T any](callback func(T) T, arr []T) []T {
	for i, v := range arr {
		res := callback(v)
		arr[i] = res
	}

	return arr
}

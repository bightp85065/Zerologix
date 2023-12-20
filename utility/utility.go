package utility

func First[T comparable](v ...T) T {
	var empty T

	for i := 0; i < len(v); i++ {
		if v[i] != empty {
			return v[i]
		}
	}

	return empty
}

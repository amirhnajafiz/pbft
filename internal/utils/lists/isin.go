package lists

// IsInList checks if an input item is inside a list.
func IsInList[T comparable](item T, list []T) bool {
	for _, key := range list {
		if key == item {
			return true
		}
	}

	return false
}

package former

func FormExpresses[T any](selections []T, systemSizes []int) [][]T {
	expresses := make([][]T, 0)

	for _, systemSize := range systemSizes {
		expresses = append(expresses, formSystemSizeExpresses(selections, systemSize)...)
	}

	return expresses
}

func formSystemSizeExpresses[T any](selections []T, systemSize int) [][]T {
	expresses := make([][]T, 0)

	expressIterator := newExpressIterator(systemSize, len(selections))
	for expressIterator.next() {
		expresses = append(expresses, get(selections, expressIterator.expressSelIndexes, systemSize))
	}

	return expresses
}

func get[T any](selections []T, expressSelIndexes []int, systemSize int) []T {
	express := make([]T, systemSize)
	for i, selIndex := range expressSelIndexes {
		express[i] = selections[selIndex]
	}

	return express
}

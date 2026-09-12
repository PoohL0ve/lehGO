package arrSlice

func Sum(figures [5]int) int {
	addFigures := 0

	for i := 0; i < len(figures); i++ {
		addFigures += figures[i]
	}
	return addFigures
}

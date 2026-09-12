package arrSlice

func Sum(figures [5]int) int {
	addFigures := 0

	for _, figure := range figures {
		addFigures += figure
	}
	return addFigures
}

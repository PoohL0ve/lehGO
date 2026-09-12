package arrSlice

func Sum(figures []int) int {
	addFigures := 0

	for _, figure := range figures {
		addFigures += figure
	}
	return addFigures
}

package arrSlice

func Sum(figures []int) int {
	addFigures := 0

	for _, figure := range figures {
		addFigures += figure
	}
	return addFigures
}

func SumAll(numbersToSum ...[]int) []int {
	var addAll []int

	for _, numbers := range numbersToSum {
		addAll = append(addAll, Sum(numbers))
	}

	return addAll
}

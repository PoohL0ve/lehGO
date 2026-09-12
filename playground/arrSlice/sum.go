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

func SumAllTails(figuresToSum ...[]int) []int {
	var tailsSum []int

	for _, figures := range figuresToSum {
		if len(figures) == 0 {
			tailsSum = append(tailsSum, 0)
		} else {
			extract := figures[1:]
			tailsSum = append(tailsSum, Sum(extract))
		}
	}
	return tailsSum
}

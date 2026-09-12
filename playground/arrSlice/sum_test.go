package arrSlice

import (
	"slices"
	"testing"
)

func TestSum(t *testing.T) {
	t.Run("Dynamic collection", func(t *testing.T) {
		numbers := []int{1, 2, 4}
		actual := Sum(numbers)
		expected := 7

		if actual != expected {
			t.Errorf("Wanted: %d but Received: %d from %v", expected, actual, numbers)
		}
	})
}

func TestSumAll(t *testing.T) {
	actual := SumAll([]int{1, 2, 4}, []int{7, 3})
	expected := []int{7, 10}

	if !slices.Equal(actual, expected) {
		t.Errorf("Wanted: %d but Got: %d", expected, actual)
	}
}

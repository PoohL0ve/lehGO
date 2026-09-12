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

func TestSumAllTails(t *testing.T) {
	// Hides a variables and functions that do not need to be exposed
	// ...also adds some type-safety
	checksums := func(t *testing.T, actual, expected []int) {
		t.Helper()
		if !slices.Equal(actual, expected) {
			t.Errorf("Wanted: %d but Got: %d", expected, actual)
		}
	}

	t.Run("Sum of tails for slices with size more than on", func(t *testing.T) {
		actual := SumAllTails([]int{1, 2, 4}, []int{7, 3})
		expected := []int{6, 3}
		checksums(t, actual, expected)
	})

	t.Run("Handling empty slices with 0", func(t *testing.T) {
		actual := SumAllTails([]int{}, []int{3, 5, 7}, []int{9, 8})
		expected := []int{0, 12, 8}
		checksums(t, actual, expected)
	})

}

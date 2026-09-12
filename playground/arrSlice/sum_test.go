package arrSlice

import "testing"

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

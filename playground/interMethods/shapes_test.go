package intermethods

import "testing"

func TestPerimeter(t *testing.T) {
	rectangle := Rectangle{12.2, 4.5}
	actual := Perimeter(rectangle)
	expected := 33.4

	if actual != expected {
		t.Errorf("Wanted Perimeter of %.2f but got a perimeter of %.2f instead", expected, actual)
	}
}

func TestArea(t *testing.T) {
	checksums := func(t *testing.T, actual, expected float64) {
		t.Helper()
		if actual != expected {
			t.Errorf("Needed an area of %g but got %g instead", expected, actual)
		}
	}

	t.Run("Area for Rectangle", func(t *testing.T) {
		rectangle := Rectangle{10.5, 2.0}
		actual := rectangle.Area()
		expected := 21.0
		checksums(t, actual, expected)
	})

	t.Run("Area for Circle", func(t *testing.T) {
		circle := Circle{10.0}
		actual := circle.Area()
		expected := 314.1592653589793

		checksums(t, actual, expected)
	})
}

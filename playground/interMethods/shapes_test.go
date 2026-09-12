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
	checkArea := func(t *testing.T, shape Shape, expected float64) {
		t.Helper()
		actual := shape.Area()
		if actual != expected {
			t.Errorf("Needed an area of %g but got %g instead", expected, actual)
		}
	}

	t.Run("Area for Rectangle", func(t *testing.T) {
		rectangle := Rectangle{10.5, 2.0}
		checkArea(t, rectangle, 21.0)
	})

	t.Run("Area for Circle", func(t *testing.T) {
		circle := Circle{10.0}
		checkArea(t, circle, 314.1592653589793)
	})
}

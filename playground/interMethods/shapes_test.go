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
	// Using Table-Driven Tests: For creating a list of test cases
	// ... to be tested the same way
	areaTests := []struct {
		shape    Shape
		expected float64
	}{
		{Rectangle{10.5, 2.0}, 21.0},
		{Circle{10.0}, 314.1592653589793},
		{Triangle{12.0, 6.0}, 36.0},
	}

	// Iterate over the struct object
	for _, areat := range areaTests {
		actual := areat.shape.Area()
		if actual != areat.expected {
			t.Errorf("Needed an area of %g but got %g instead", areat.expected, actual)
		}
	}
}

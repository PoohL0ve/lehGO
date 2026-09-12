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
		name    string
		shape   Shape
		hasArea float64
	}{
		// Added named fields
		{name: "Rectangle", shape: Rectangle{Width: 10.5, Height: 2.0}, hasArea: 21.0},
		{name: "Circle", shape: Circle{Radius: 10.0}, hasArea: 314.1592653589793},
		{name: "Triangle", shape: Triangle{Base: 12.0, Height: 6.0}, hasArea: 36.0},
	}

	// Iterate over the struct object
	for _, areat := range areaTests {
		t.Run(areat.name, func(t *testing.T) {
			actual := areat.shape.Area()
			if actual != areat.hasArea {
				t.Errorf("%#v got %g want %g", areat.shape, actual, areat.hasArea)
			}
		})
	}
}

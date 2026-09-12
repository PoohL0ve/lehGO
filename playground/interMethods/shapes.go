package intermethods

import "math"

type Rectangle struct {
	// Structs are collections of fields
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * math.Pow(c.Radius, 2)
}

func Perimeter(rectangle Rectangle) float64 {
	perimeter := 2 * (rectangle.Height + rectangle.Width)
	return perimeter
}

func Area(rectangle Rectangle) float64 {
	area := rectangle.Height * rectangle.Width
	return area
}

package integers

import (
	"fmt"
	"testing"
)

func TestAdder(t *testing.T) {
	sum := Add(2, 2)
	expected := 4

	if sum != expected {
		t.Errorf("expect %d but got %d", expected, sum)
	}
}

func ExampleAdd() {
	// Documentation will always reflect current code behaviour
	// The comment at the button ensure the ExampleAdd() is executed in the test
	sum := Add(1, 5)
	fmt.Println(sum)
	// Output: 6
}

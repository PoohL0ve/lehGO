package iteration

import (
	"fmt"
	"testing"
)

func TestRepeat(t *testing.T) {
	actual := Repeat("c", 7)
	expected := "ccccccc"

	if actual != expected {
		t.Errorf("Expected: %q\nGot: %q", expected, actual)
	}
}

func BenchmarkReapeat(b *testing.B) {
	for b.Loop() {
		Repeat("c", 5)
	}
}

func ExampleRepeat() {
	repeatChar := Repeat("d", 9)
	fmt.Println(repeatChar)
	// Output: ddddddddd
}

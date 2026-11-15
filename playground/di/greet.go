package di

import (
	"fmt"
	"io"
)

// Greet accepts an io.Writer interface instead of hardcoding os.Stdout
// The function can be called from anywhere that imports the package
// ... Greet(os.Stdout, string)
func Greet(writer io.Writer, name string) {
	fmt.Fprintf(writer, "Hello, %s", name)
}

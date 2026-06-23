# All About GO
__Golang__ is an open-sourced _statically_ (code is checked before runtime) typed, _compiled_ (translates code to binary code) programming language. It compiles faster than _Rust_, _C_, and _Zig_ but runs slower than them. It was created by system engineers at __Google__ (Ken Thompson, Rob Pike, Robert Griesemer) to provide better performance through speed, scaling, trivial concurrency, and readability than their previous infrastructure. It was designed to be used in network and backend infrastructures but can be used anywhere really.

Although Go is lightweight, its programs produce extra memory that is stored in the _executable binary_ that the __Go Runtime__ cleans up, as it includes a _garbage collector_.

While the `func main() {}` is the main function of the program and the entry point, the __package main__ at the top tells the Go compiler that the code has to be run and compiled as a standalone package. The `import fmt` imports the formatting package from the _standard library_ which allows things like `fmt.Println()` to be used.

```bash
brew install go

go version
which go
```
Go __comments__ serve a dual technical purpose: they are structural documentation. The built-in tool `go doc` parses your comments automatically to generate documentation for your code packages.
```go
// In-line: use for functions and most other things
/* Block: for package-specific */
```
What and Why GO
History
Setup
Comments

__Contents__ to master:
- [Basics](#understand-the-basics)
- [Composite Types](#composite-types)
- [Conditionals and Loops](#making-decisions-and-loops)
- [Functions](#functions)
- [Pointers](#pointers)

## Understand the Basics
### 📇 Variables and Constants
For memory safety variables are explicitly bound to a specific data and there are two syntax used for declaring them:
- __Standard Declaration__: Uses the `var` keyword, explicitly state the name, type, and optional value. This allows a developer to initialise a variable without a value or create one at the package level (outside functions).
- __Shorthand Declaration__: Used inside blocks like functions and loops to provide clean and readable code. It does not specify a type as the compiler would infer the type based on the value at the right of walrus operator`:=`.
```go
package main
import "fmt"

// Package-level
var name string = "Someone Like You"
var (
    age int = 40
    available bool = True
)

func main() {
    // Standard declaration inside a function (often used when value is unknown yet)
	var operationalLimit int
	operationalLimit = 500

	// Shorthand declaration (Type inferred as string)
	environmentName := "Staging_Cluster"
}
```
If a variable is declared with `var` without a value the GO compiler assigns a __Zero Value__ to it instead of undefined like most languages.

Go uses lexical block scoping enforced by curly braces {}. Variables declared inside an inner block (like an if statement, for loop, or function) are invisible to the outer blocks. __Shadowing__ occurs when you declare a variable in an inner block with the exact same name as a variable in an outer block. The inner variable "shadows" (hides) the outer variable, meaning changes inside that block do not affect the outer scope.

___Critical Rules of Go Variables___:
- The Unused Variable Constraint: The Go compiler strictly treats unused local variables as a build error. If you declare a variable inside a function and do not read from it, your program will refuse to compile. This keeps binary sizes lean and stops dead code accumulation.
- The Blank Identifier (_): If a function returns multiple values (such as data and an error) and you do not need one of those variables, you must assign it to the blank identifier _ to discard it safely without triggering the unused variable compiler error.

__Constants__ are static primitive values that cannot be changed or re-assigned. They are declared using the `const` keyword and the assignment operator not the _walrus_. They can also be grouped with `()`. They have to be computed at compile time, therefore using things like `time.Now()` would not work for a constant.
```go
package main
func main() {
	const hello = "Hello World!"
} 
```

When it comes to formatting in go there are two ways in which it can be done using the `fmt` package:
- `fmt.Printf()`: Prints to the standard out
- `fmt.Sprintf()`: returns the formatted string

The default symbol is the `%v` which returns any value, while `%s`, `%d`, and `%f` are for string, integers, and floats respectively. The `%T` tells the type of a value.
```go
s := fmt.Sprintf("You have %.2f points", 98.4658)
// You have 98.47 points

fmt.Printf("The type of penniesPerText is %T\n", penniesPerText)
```
[_learn more formatting_](https://pkg.go.dev/fmt#hdr-Printing)

### 🗑️ Basic Data Types
__Signed__ and __Unsigned__ integers are numeric values that do not have decimals where unsigned are only positive values. The figure after the `int` or `uint` respresents the number of bits which is the size of the value:
- __Signed__: int (32 or 64 bits based on user environment), int8, int16, int32, int64
- __Unsigned__: unit, uint8, uint16, uint32, uint64
- __Float__: float32, float64 (use)
- __Complex__: complex64, complex128 (use)
- __Byte__: byte an alias of uint8
- __Rune__: rune an alias of int32

```go
// Convert from float64 to int64
temperatureFloat := 88.26
temperatureInt := int64(temperatureFloat)
```
With number types, only differ from the defaults if its absolutely necessary like for performance and memory.



### 🎥 Documentation and Commands in GO

## Composite Types
### Arrays
### Strings
Two strings can be concatenated with the `+` operator but string and an int or float64 cannot. 
### Slices
### Maps

## Making Decisions and Loops
### ❓ Conditionals
The code of a condition is executed only if the statement is true. In _Go_ the `if/else` keywords do not need parentheses but they need brackets on the same line as the condition. Variables can also be declared in the statement.
```go
if age := 41; age < 30 {
	fmt.Println(age, "is not in the dating bracket")
} else if age > 27 && age < 35 {
	fmt.Println(age, "not too sure")
} else {
	fmt.PrintLn(age, "is perfect")
}
```
For __switch__ statements in Go, the _break_ keyword is not required at the end of the __case__, as its implicit. However, the `fallthrough` statement can be used if a case has to fall through to another.
```go
switch someone {
case "No one":
	message = "Keep looking"
case "Her":
	message = "You deserve happiness"
default:
	message = "I guess not in this lifetime"
}
```
### ⭕️ Iterations

## Functions

## Pointers

## Methods and Interfaces
### Interfaces
### Generics

## Dealing with Errors

## Organising Code
### Module and Dependencies
### Packages
### Publishing Modules

## Resources
- [_Effective GO_](https://go.dev/doc/effective_go)
- [_GO By Example_](https://gobyexample.com/)
- [_Golang tutiral series_](https://golangbot.com/learn-golang-series/)
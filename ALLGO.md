# All About GO
__Golang__ is an open-sourced _statically_ (code is checked before runtime) typed, _compiled_ (translates code to binary code) programming language. It compiles faster than _Rust_, _C_, and _Zig_ but runs slower than them. It was created by system engineers at __Google__ (Ken Thompson, Rob Pike, Robert Griesemer) to provide better performance through speed, scaling, trivial concurrency, and readability than their previous infrastructure. It was designed to be used in network and backend infrastructures but can be used anywhere really.

Although Go is lightweight, its programs produce extra memory that is stored in the _executable binary_ that the __Go Runtime__ cleans up, as it includes a _garbage collector_.

While the `func main() {}` is the main function of the program and the entry point, the __package main__ at the top tells the Go compiler that the code has to be run and compiled as a standalone package. The `import fmt` imports the formatting package from the _standard library_ which allows things like `fmt.Println()` to be used.

```bash
brew install go

go version
which go
```
To run Go a folder has to be initialised with a module file tracking dependencies:
```bash
go mod init structs-demo
```
The `go.mod` file is a plain text file that declares the unique path of your module and states the version of Go your code requires. Historically, Go forced developers to put all their code inside one giant global system folder (the old $GOPATH). Go Modules changed this. A go.mod file tells the Go compiler: "Treat this specific folder as an isolated, self-contained workspace. Any files inside this folder belong to this project." Without it, the compiler cannot track imports or handle third-party libraries. Scripts cannot have _test_ in them as in `struct_test.go`.

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

## Structs
__Structs__ are used to represent structured data where data that needs to be grouped into a single unit can be placed there. __Go__ rejects OOP, avoiding inheritance hierachies and heavy virtual-table objects, therefore, structs are used. Structs operate at both ends of the abstraction spectrum:
1. __High-level__: It acts as an object where methods and complex behaviours can be attached using embedding
2. __Low-level__: It maps directly to physical CPU memory. You control the exact byte layout, stack vs. heap allocation, and how the data aligns with hardware caches.
Simple declaration:
```go
type Her struct {
	name string
	age int
	available bool
}
// Use a literal to create an instant
eba := Her {name: Eba, age: 40, available: false}
``` 
It is similar to dictionaries and objects in other programming languages. Structs can be created without specifying the fields, however, each value must be in the correct order. Structs can be nested where the fields can be accessed using the dot operator:
```go
type Her struct {
	name string
	age int
	available bool
	home house
}

type building struct {
	rooms int
	lightType string
}

tara := Her {tara.home.rooms = 5}
```
While _Java or Python_ use complex class hierarchies, Go champions __composition__ over inheritance. Instead of a class inheriting from a parent class, Go structs use struct __embedding__ (anonymous fields) to compose behaviors. If you embed struct `A` inside struct `B`, all of A's fields and methods are "promoted" to B. It looks and acts like inheritance to the caller, but under the hood, it is pure, __flat composition__. This means the extra dot notation needed to access fields like in the nested struct is not need; you directly access like normal. There are two ways to attach behaviour using __method receivers__:
1. __Value Receivers__ (func (u User)): The method receives a complete copy of the struct. Safe from side effects, but incurs copy overhead for large structs.
2. __Pointer Receivers__ (func (u *User)): The method receives a pointer to the original struct. This allows you to modify the original data and avoids copying memory.

A receiver is a function where the parameter goes before the name of the function. It is used to provide interfaces that structs and other data types can implement:
```go
func (r rect) area() int {
	return r.width * r.height
}
// struct
var r = rect{width: 5, height: 10,}
```
Structs can be anonymous which can be used when they are intended to be used once. To create an anonymous struct, just instantiate the instance immediately using a second pair of brackets after declaring the type:
```go
myCar := struct {
  brand string
  model string
} {
  brand: "Toyota",
  model: "Camry",
}
```
Anonymous structs can also be used in other structs.

There are a few ways in which structs can be broken:
1. __The Value Receiver (Silent Bug)__: A method that is supposed to modify the state of a struct which forget the asterisk will compile fine but no changes will be applied.
2. __Method Set/ Interface Trap__: Go's compiler has strict rules about "Method Sets". If an interface requires a method, and you implement that method with a pointer receiver, a value type cannot satisfy that interface.
3. __Ambiguity Collissions in Embedding__: If you embed two different structs that happen to share a field name, Go compiles fine. However, the moment you attempt to access the promoted field, the compiler panics due to ambiguity. Use dot notation to be explicit.

### Structs at the Hardware Level
At the hardware level, the operating system and CPU don't read memory 1 byte at a time. Modern 64-bit CPUs (like Apple Silicon or Intel/AMD x86-64) read memory in words of 8 bytes. To optimize access speed, every primitive type has an alignment guarantee. For example, a 64-bit integer (int64) must start at a memory address that is a multiple of 8. If you place a 1-byte boolean right before an 8-byte integer, the compiler will insert 7 "invisible" empty bytes of _padding_ to ensure the integer starts at an 8-byte boundary. If it didn't pad, the integer would be "misaligned" (straddling two 8-byte words), forcing the CPU to make two memory cycles instead of one just to read a single integer.

Structs sit in memory in a contiguous block, with fields placed one after the other, where the order of fields is important:
```go
type stats struct {
	NumPosts uint8
	Reach    uint16
	NumLikes uint8
}
```
Go will add padding to make up for the size difference which is done for execution speed. The reflect package can be used if there is concerns to debug the memory layout of a struct or manually align fields from largest to smallest:
```go
typ := reflect.TypeOf(stats{})
fmt.Printf("Struct is %d bytes\n", typ.Size())
```

Go's compiler lets you attach literal metadata string __tags__ to fields (e.g., json:"id"). Go frameworks use the standard __reflect__ library at runtime to parse these tags. This is how JSON serialization, database ORMs, and configuration binders inspect and map structural data.

Empty structs are the smallest possible type in Go as they take up 0 bytes of memory. They are used as _unary_ values useful for:
1. __Set Implementations__: Using a map where you only care about keys: `map[string]struct{}`. This consumes far less memory than `map[string]bool`.
2. __Channel Signaling__: Sending a pure signal without allocating memory: `chan struct{}`.

```go

// anonymous empty struct type
empty := struct{}{}

// named empty struct type
type emptyStruct struct{}
empty := emptyStruct{}
```


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
Functions are created for __reusability__ where any number of inputs are given, some calculation is performed on them and an output is returned. The parameter and return type are not necessary but makes it easier to know what values the function accepts and returns. When multiple arguments are used of the same type, the type only needs to be declared after the last one.
```go
func functionName(parameter parameterType) returnType {
	// function body
}

func getFullName(firstName, lastName string) string {}
```
In Go, a function cannot pass data back to the caller unless it __explicitly declares__ what type of data it is returning inside its signature. There can also be __named__ return values, which will be declared in the function signature. It allows the _return_ keyword to be used by itself, where the function would know what values to return. _Naked_ return statements should only be used in short funtions as they can harm readability.
```go
func getNumber() (name string, number int) {
	return // naked return
	// They can still be explicitly returned or over written in the return statement
}
```
However, if the variable/s are shadowed (re-assigned) in the function it would break the function. They (_named returns_) are good for longer functions and documentation which eliminates the need for comments as what the function returns would already be known. Some return values can be ignored by using the __blank identifier__ `_`.
```go
func getPoint() (x int, y int) {
	return 3, 4
}
x, _ : getPoint()
```
A good practice is to use _guard clauses_ with _if_ statements to return early, which is can also be a simple form of error handling. It makes the code simpler and more linear logic is understood and written.

As functions are another data type, they can be used as arguments to a function, this is called __first-class__ and __higher-order__ functions. 
```go
func aggregate(a, b, c int, arithmetic func(int, int) int) int {
  firstResult := arithmetic(a, b)
  secondResult := arithmetic(firstResult, c)
  return secondResult
}
// It can call any function that takes two arguments
```
### Anonymous Functions
They are simply functions without a name. They are vital for localizing logic. Instead of polluting your package space with a global function that is only ever used in one specific spot, you write it inline. They are heavily leveraged when launching concurrent tasks (Goroutines), deferring cleanup tasks (defer), or creating Closures.

__Mental Model 🧠__: An anonymous function is a disposable zip-top bag. You materialize it right on the spot to hold a specific set of ingredients (variables from the local scope). You mix the ingredients inside it, use it for that exact meal, and then throw it away. If it's a closure, that bag keeps holding onto those specific local ingredients even if you carry the bag into a completely different room (different execution scope).

```go
newZ := func(amount float64) float64 {
		return amount * 0.15
	}

// Anonymous function inside another function
printCostReport(func(msg string) int {
		return len(msg) * 4
	},outro,)
```

__How Can It Be Broken__: They linked themselves to variables outside their scope which is common for subtle bugs to be introduced:
1. __The Concurrency Pointer/Race Condition Trap__: When you spin up an anonymous function inside a concurrent loop (Goroutine), it binds to the reference of the outer variable, not a snapshot of its value. To fix it, you must explicitly pass the data into the anonymous function as a hard parameter, forcing Go to make a safe local copy of the value for that specific execution block.
2. __Hidden State Corruption__: Because closures can modify outer variables, passing a closure into deep infrastructure layers can cause code miles away to change your local state without you realizing it, making debugging a nightmare

The __defer__ keyword allows a function to be executed automatically just before its enclosing function returns. The deferred call's arguments are evaluated immediately, but the function call is not executed until the surrounding function returns. Deferred functions are generally used to clean up resources that are no longer being used like file handlers and database connections.
```go
func GetUsername(dstName, srcName string) (username string, err error) {
	// Open a connection to a database
	conn, _ := db.Open(srcName)

	// Close the connection *anywhere* the GetUsername function returns
	defer conn.Close()

	username, err = db.FetchUser()
	if err != nil {
		// The defer statement is auto-executed if we return here
		return "", err // defer is called here
	}

	// The defer statement is auto-executed if we return here
	return username, nil // defer is also called here
}
```
Multiple defer statements are executed in LIFO (last-in/firstout) format.

__N.B__: GO is __blocked scoped__ often denoted by `{}`. Variables can be hidden inside blocks, which implies a new scope and are not accessible by things outside of it:
```go
{
	age: 19
}
fmt.Println(age) // Not accessible
```

A __closure__ occurs when an anonymous function captures and wrapts itself around variables defined outside its own body, which retains access to the variables even when the outer block finishes executing. In other words, it mutates the variable.

__Currying__ is a functional programming technique where a function that takes multiple arguments is transformed into a chain of nested functions, each taking exactly one argument (or a smaller subset of arguments). Why do this in Go?
- __Partial Application / Configuration Presets__: It allows you to lock in configuration parameters or dependencies early in your pipeline, generating highly specialized helper functions from a generic base function.
- __Lazy Evaluation__: You can pass data through your system step-by-step, delaying final execution until the absolute last argument becomes available.
```go
func getLogger(formatter func(string, string) string) func(string, string) {
	return func(str1, str2 string) {
		confused := formatter(str1, str2)
		fmt.Println(confused)
	}
}
```
__How Can It Be Broken__:
1. Hidden Memory Leaks: Because currying heavily relies on closures, each nested function retains pointers to the variables in its outer scope. If your first argument is a massive data structure (like a large byte slice or an entire database configuration struct), that massive chunk of memory cannot be cleared by Go's garbage collector as long as your curried function variables (dbErrorLogger in the example above) remain alive in your program's memory.
2. Complete Code Unreadability: Go was intentionally designed to be hyper-explicit and simple to read. When you curry heavily, function signatures turn into a nested, confusing mess of anonymous declarations: `func(string) func(int) func(bool) func(error)`. If a teammate tries to read your code, they will have to trace multiple execution hops just to understand what a function actually outputs. Over-currying forces idiomatic Go to look like unmaintainable academic code, breaking the readability philosophy of the language.



## Pointers

## Methods and Interfaces
### Interfaces
An __interface__ is a type that specifies a set of methods that declares a protocol for behaviour. It is a contractual agreement that does not care about the type implementing it. In _Go_, interfaces are implicit where types using them, like structs do not use the `implements` keyword.
```go
// Define the Interface (Contract)
type DBRepository interface {
	GetUserName(id int) (string, error)
}

type SQLDatabase struct {
	ConnectionString string
}

func (db SQLDatabase) GetUserName(id int) (string, error) {
	// Pretend this connects to an actual physical SQL database server
	return fmt.Sprintf("User_%d_From_Postgres", id), nil
}
```
Now to use the method, simply type `db.GetUserName()`.
Implicit interfaces decouple the definition of an interface from its implementation. You may add methods to a type and in the process be unknowingly implementing various interfaces, and that's okay. A type can implement any number of interfaces. The __empty interface__ `interface{}` is always implemented because it hs no requirement.

Parameters of an interface can be named so it can be clear what type of data they are:
```go
type Copier interface {
  Copy(sourceFile string, destinationFile string) (bytesCopied int)
}
```
__Type assertions__ can be used to check the underlying type of an interface.
```go
func getExpenseReport(e expense) (string, float64) {
	em, echeck := e.(email)
	sm, scheck := e.(sms)

	if echeck {
		return em.toAddress, em.cost()
	}

	if scheck {
		return sm.toPhoneNumber, sm.cost()
	}

	return "", 0.0
}
// e is the interface and email and sms are the struct types
```
If multiple types can be a the potential, then __type switch__ should be used where the case values are types instead of actual values:
```go
func getExpenseReport(e expense) (string, float64) {
	switch guess := e.(type) {
	case email:
		return guess.toAddress, guess.cost()
	case sms:
		return guess.toPhoneNumber, guess.cost()
	default:
		return "", 0.0
	}
}
```
[Best Practices for Interfaces in Go](https://www.boot.dev/blog/golang/golang-interfaces)

At the hardware/runtime layer, an interface variable is not a plain pointer. It is a two-word wide structure containing two memory addresses:
1. `itab` pointer: Points to an Interface Table. The `itab` tracks the concrete type of the value metadata and lists pointers to the actual executable functions matching the interface contract.
2. `data` pointer: Points to the physical location of the concrete value (the data) allocated on the stack or heap.
### Generics

## Dealing with Errors
An error in Go is any value that satisfies the built-in error interface:
```go
type error interface {
    Error() string
}
```
__Why Go does this__: Languages like Java, Python, or JavaScript use exceptions that jump up the call stack and break execution. This creates invisible control paths—you never really know which line of code might randomly throw an exception and crash the thread. Go forces errors to be explicit, local, and predictable. If a function can fail, it returns a `(result, error)` tuple, and you handle it right there.

Because errors are just interfaces, you can build your own custom types that implement the error interface:
```go
type userError struct {
    name string
}

func (e userError) Error() string {
    return fmt.Sprintf("%v has a problem with their account", e.name)
}
```
Using the `errors.New()` function:
```go
package main

import (
	"errors"
)

func divide(x, y float64) (float64, error) {
	if y == 0 {
		return 0.0, errors.New("no dividing by 0")
	}
	return x / y, nil
}
```
 

## Organising Code
### Module and Dependencies
### Packages
### Publishing Modules

## Resources
- [_Effective GO_](https://go.dev/doc/effective_go)
- [_GO By Example_](https://gobyexample.com/)
- [_Golang tutiral series_](https://golangbot.com/learn-golang-series/)


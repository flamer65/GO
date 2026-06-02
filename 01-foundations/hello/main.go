// ============================================================================
// MODULE 1.1: Hello World & Project Setup
// ============================================================================
//
// KEY CONCEPTS:
//   - Every Go file starts with a "package" declaration
//   - The "main" package is special — it defines an executable program
//   - The "main" function is the entry point (like main() in C or Java)
//   - "import" brings in other packages you want to use
//   - Go enforces clean imports — unused imports cause compile errors!
//
// HOW TO RUN:
//   go run main.go
//   go run main.go Naman
//
// HOW TO BUILD (creates a binary):
//   go build -o greeter .
//   ./greeter Naman
//
// ============================================================================

package main // Every executable program must be in package "main"

import (
	"fmt" // fmt = "format" — Go's primary I/O formatting package
	"flag"
)

func main() {
	name := flag.String("name", "World", "A name to greet")
	count := flag.Int("count", 1, "Number of times to greet")
	shout := flag.Bool("shout", false, "Shout the greeting")

	flag.Parse()
	fmt.Println("Name is:", *name)
	fmt.Println("Count is:", *count)
	fmt.Println("Shout is:", *shout)
	
	
}

// ============================================================================
// 🏋️ EXERCISE: Build on this!
//
// 1. Add a "--shout" flag: if passed, print the greeting in UPPERCASE
//    Hint: look at the "strings" package — strings.ToUpper()
//
// 2. Add a "--count N" flag: repeat the greeting N times
//    Hint: you'll need the "strconv" package — strconv.Atoi() to convert
//    a string to an integer
//
// 3. [BONUS] Use the "flag" package instead of manually parsing os.Args
//    The flag package handles argument parsing elegantly. Try:
//      name := flag.String("name", "World", "person to greet")
//      flag.Parse()
//
// ============================================================================


// -----------------------------------------------------------------------
	// PART 1: The simplest Go program
	// -----------------------------------------------------------------------

	// fmt.Println prints a line to stdout (adds a newline at the end)
	// age := 23
	// fmt.Println("Hello, World!", "I am", age, "years old")

	// fmt.Printf is like C's printf — uses format verbs:
	//   %s = string, %d = integer, %f = float, %v = any value, %T = type
	//fmt.Printf("This is Go version: %s and I am %d years old\n", "1.22.5", age)

	// -----------------------------------------------------------------------
	// PART 2: Command-line arguments
	// -----------------------------------------------------------------------

	// os.Args is a slice (dynamic array) of strings
	// os.Args[0] is always the program name itself
	// os.Args[1:] are the actual arguments

	// fmt.Printf("\nProgram name: %s\n", os.Args[0])
	// fmt.Printf("Number of arguments: %d\n", len(os.Args)-1)

	// Check if user provided a name argument
	// if len(os.Args) > 1 {
	// 	name := os.Args[1] // := is "short variable declaration" (type inferred)
	// 	fmt.Printf("Hello, %s! Welcome to Go. 🚀\n", name)
	// } else {
	// 	fmt.Println("Hello, stranger! Try: go run main.go YourName")
	// }

	// -----------------------------------------------------------------------
	// PART 3: Multiple arguments
	// -----------------------------------------------------------------------

	// If more than one argument, greet all of them
	// if len(os.Args) > 2 && os.Args[1] == "--shout" {
	// 	fmt.Println(strings.ToUpper("\nGreeting everyone:"))
	// 	// "range" iterates over slices. It returns (index, value).
	// 	// We use "_" to discard the index since we don't need it.
	// 	for _, name := range os.Args[1:] {
	// 		fmt.Printf("  👋 Hi, %s!\n", name)
	// 	}
	// }

	// -----------------------------------------------------------------------
	// EXERCISE: Shout Flag
	// -----------------------------------------------------------------------
	// if len(os.Args) > 2 && os.Args[1] == "--shout" {
	// 	greeting := "Greetings everyone"
	// 	fmt.Println("\n", strings.ToUpper(greeting))
	// }
	// if len(os.Args) > 2 && os.Args[1] == "--count" {
	// 	// strconv.Atoi returns TWO values: the integer, and an error.
	// 	// We have to capture both!
	// 	count, err := strconv.Atoi(os.Args[2])
	// 	if err != nil {
	// 		fmt.Println("Error: count must be a number!")
	// 		return // exit if it's not a valid number
	// 	}
	// 	greeting := "Greetings everyone"
	// 	for range count {
	// 		fmt.Println("\n", strings.ToUpper(greeting))
	// 	}
	// }
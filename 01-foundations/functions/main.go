// ============================================================================
// MODULE 1.4: Functions
// ============================================================================
//
// KEY CONCEPTS:
//   - Functions can return MULTIPLE values (killer feature for error handling)
//   - Named return values act as documentation and allow "naked returns"
//   - Functions are first-class: assign to variables, pass as arguments
//   - Closures capture and remember their surrounding scope
//   - "defer" schedules cleanup that runs when the function exits
//   - "panic/recover" are Go's exception mechanism (use sparingly!)
//
// RUN: go run main.go
// ============================================================================

package main

import (
	"fmt"
	"math"
	"strings"
)

func main() {
	// -----------------------------------------------------------------------
	// PART 1: Basic Functions
	// -----------------------------------------------------------------------

	fmt.Println("=== Basic Functions ===")
    
	result := add(3, 5)
	fmt.Printf("add(3, 5) = %d\n", result)

	fmt.Printf("greet(\"Naman\") = %s\n", greet("Naman"))

	// -----------------------------------------------------------------------
	// PART 2: Multiple Return Values
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Multiple Return Values ===")

	// Go functions can return multiple values — THIS IS HUGE
	// The most common pattern: (result, error)
	quotient, remainder := divide(17, 5)
	fmt.Printf("17 ÷ 5 = %d remainder %d\n", quotient, remainder)

	// The error pattern — we'll explore this deeply in Module 3
	result2, err := safeDivide(10, 3)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("10 / 3 = %.4f\n", result2)
	}

	// What happens with division by zero?
	_, err = safeDivide(10, 0)
	if err != nil {
		fmt.Println("Caught error:", err)
	}

	// -----------------------------------------------------------------------
	// PART 3: Named Return Values
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Named Return Values ===")

	// Named returns serve as documentation — they tell you what's returned
	area, perim := rectangleProps(5.0, 3.0)
	fmt.Printf("Rectangle 5×3: area=%.1f, perimeter=%.1f\n", area, perim)

	// -----------------------------------------------------------------------
	// PART 4: Variadic Functions (variable number of arguments)
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Variadic Functions ===")

	// The ... means "zero or more arguments of this type"
	fmt.Printf("sum() = %d\n", sum())
	fmt.Printf("sum(1) = %d\n", sum(1))
	fmt.Printf("sum(1, 2, 3, 4, 5) = %d\n", sum(1, 2, 3, 4, 5))

	// You can "spread" a slice into a variadic function
	numbers := []int{10, 20, 30}
	fmt.Printf("sum(10, 20, 30) = %d\n", sum(numbers...)) // Note the ...

	// -----------------------------------------------------------------------
	// PART 5: Functions as Values (First-Class Functions)
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Functions as Values ===")

	// Functions can be assigned to variables
	var mathOp func(int, int) int // Declare a function-typed variable
	mathOp = add
	fmt.Printf("mathOp(10, 20) = %d\n", mathOp(10, 20))

	mathOp = multiply
	fmt.Printf("mathOp(10, 20) = %d\n", mathOp(10, 20))

	// Functions can be passed as arguments
	fmt.Println("\nApplying operations:")
	applyAndPrint("add", 10, 20, add)
	applyAndPrint("multiply", 10, 20, multiply)
	applyAndPrint("subtract", 10, 20, func(a, b int) int { return a - b })

	// -----------------------------------------------------------------------
	// PART 6: Closures
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Closures ===")

	// A closure is a function that "captures" variables from its outer scope
	// The captured variables survive even after the outer function returns

	// Counter factory — each call creates an independent counter
	counter1 := makeCounter()
	counter2 := makeCounter()

	fmt.Printf("counter1: %d, %d, %d\n", counter1(), counter1(), counter1())
	fmt.Printf("counter2: %d, %d\n", counter2(), counter2())
	fmt.Printf("counter1: %d (continues from 3!)\n", counter1())

	// Practical closure: create a multiplier
	double := makeMultiplier(2)
	triple := makeMultiplier(3)
	fmt.Printf("double(5) = %d\n", double(5))
	fmt.Printf("triple(5) = %d\n", triple(5))

	// -----------------------------------------------------------------------
	// PART 7: defer
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Defer ===")

	// "defer" schedules a function call to run when the enclosing function exits
	// Think of it as a "cleanup" guarantee — like "finally" in other languages
	// Deferred calls execute in LIFO (stack) order

	fmt.Println("Start")
	defer fmt.Println("Deferred 1 (runs last)")
	defer fmt.Println("Deferred 2 (runs second)")
	defer fmt.Println("Deferred 3 (runs first)")
	fmt.Println("End")
	// Output order: Start, End, Deferred 3, Deferred 2, Deferred 1

	// Common use case: resource cleanup
	// defer file.Close()
	// defer db.Disconnect()
	// defer mutex.Unlock()

	// -----------------------------------------------------------------------
	// PART 8: panic and recover
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Panic & Recover ===")

	// "panic" stops normal execution (like throwing an exception)
	// "recover" catches a panic (only works inside a deferred function)
	// Rule: prefer returning errors. Use panic only for truly unrecoverable situations.

	fmt.Println("Before safeDivision")
	result3 := safeDivision(10, 0) // Would panic, but we recover
	fmt.Printf("safeDivision result: %d\n", result3)
	fmt.Println("After safeDivision — program continues!")

	// -----------------------------------------------------------------------
	// PART 9: Calculator Exercise (SOLVED)
	// -----------------------------------------------------------------------

	fmt.Println("\n=== 🧮 Calculator with First-Class Functions ===")

	// Store operations in a map of functions
	operations := map[string]func(float64, float64) float64{
		"+": func(a, b float64) float64 { return a + b },
		"-": func(a, b float64) float64 { return a - b },
		"*": func(a, b float64) float64 { return a * b },
		"/": func(a, b float64) float64 {
			if b == 0 {
				fmt.Println("  ⚠️  Cannot divide by zero")
				return math.NaN()
			}
			return a / b
		},
		"^": func(a, b float64) float64 { return math.Pow(a, b) },
	}

	// Use the calculator
	testCases := []struct {
		a, b float64
		op   string
	}{
		{10, 3, "+"},
		{10, 3, "-"},
		{10, 3, "*"},
		{10, 3, "/"},
		{2, 10, "^"},
		{10, 0, "/"},
	}

	for _, tc := range testCases {
		if fn, exists := operations[tc.op]; exists {
			result := fn(tc.a, tc.b)
			fmt.Printf("  %.0f %s %.0f = %.4f\n", tc.a, tc.op, tc.b, result)
		} else {
			fmt.Printf("  Unknown operation: %s\n", tc.op)
		}
	}

	// -----------------------------------------------------------------------
	// BONUS: Functional programming patterns
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Functional Patterns ===")

	// Decorator pattern with closures
	loudGreet := withPrefix("🔊 ", greet)
	quietGreet := withSuffix(" 🤫", greet)

	fmt.Println(loudGreet("World"))
	fmt.Println(quietGreet("World"))
	addOne := func(x int) int { return x + 1 }
	doubles := func(x int) int { return x * 2 }
	fmt.Println("addOne", addOne(2))
	fmt.Println("double", doubles(2))
	fmt.Println(doubles(10))
	results := composer(addOne, doubles)
	fmt.Println("result", results(10))

}

// -----------------------------------------------------------------------
// Function Definitions
// -----------------------------------------------------------------------

func composer(add, double func(x int) int) func(x int) int{
	return func(x int) int {
		return double(add(x))
	}
}

// Basic function with parameters and a return value
func add(a int, b int) int {
	return a + b
}

// When consecutive parameters share a type, you can combine them
func multiply(a, b int) int {
	return a * b
}

// Function returning a string
func greet(name string) string {
	return "Hello, " + name + "!"
}

// Multiple return values
func divide(a, b int) (int, int) {
	return a / b, a % b
}

// Multiple return values with error (THE Go pattern)
func safeDivide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("cannot divide %.1f by zero", a)
	}
	return a / b, nil // nil means "no error"
}

// Named return values — serve as documentation
func rectangleProps(length, width float64) (area float64, perimeter float64) {
	area = length * width
	perimeter = 2 * (length + width)
	return // "naked return" — returns the named values
}

// Variadic function — accepts any number of ints
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// Higher-order function — takes a function as parameter
func applyAndPrint(name string, a, b int, op func(int, int) int) {
	result := op(a, b)
	fmt.Printf("  %s(%d, %d) = %d\n", name, a, b, result)
}

// Closure — returns a function that captures "count"
func makeCounter() func() int {
	count := 0 // This variable is "captured" by the returned function
	return func() int {
		count++ // Modifies the captured variable
		return count
	}
}

// Closure factory — returns a multiplier function
func makeMultiplier(factor int) func(int) int {
	return func(x int) int {
		return x * factor
	}
}

// Panic + recover example
func safeDivision(a, b int) int {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("  Recovered from panic: %v\n", r)
		}
	}()

	// This will panic if b == 0
	return a / b
}

// Decorator pattern: add a prefix to a function's output
func withPrefix(prefix string, fn func(string) string) func(string) string {
	return func(s string) string {
		return prefix + fn(s)
	}
}

// Decorator pattern: add a suffix
func withSuffix(suffix string, fn func(string) string) func(string) string {
	return func(s string) string {
		return fn(s) + suffix
	}
}

func memoize(fn func(int) int) func(int) int{
  cache := make(map[int]int)
  return func(n int) int{
    if value, ok := cache[n]; ok {
	  return value
	}	
    value := fn(n)
	cache[n] = value
	return value
  }
}
func retry(attempts int, fn func() error) error {
	var err error
	for i := 1; i < attempts; i++ {
		fmt.Printf("Attempt %d...\n", i)
		err = fn()

		if(err == nil){
			return nil
		}
		fmt.Printf("  Failed: %v\n", err)
	}
	return fmt.Errorf("Failed after %d attempts. last error: %v", attempts, err)

}
func pipeline(fns ...func(int)int) func(int) int {
	return func(x int) int {
		result := x
		for _, fn := range fns {
			result = fn(result)
		}
		return result
	}
}
// Ensure strings package is used
var _ = strings.ToUpper

// ============================================================================
// 🏋️ EXERCISES:
//
// 1. COMPOSE FUNCTION: Write a function compose(f, g func(int) int)
//    that returns a new function h where h(x) = f(g(x))
//    Test it with double and addOne functions.
//
// 2. MEMOIZE: Write a function memoize(fn func(int) int) func(int) int
//    that caches results. When called with the same argument, return
//    the cached value instead of recomputing.
//    Test it with an expensive Fibonacci function.
//
// 3. RETRY: Write a function retry(attempts int, fn func() error)
//    that calls fn up to "attempts" times until it succeeds (returns nil).
//    This is a real-world pattern used in network calls!
//
// 4. [BONUS] PIPELINE: Write a function pipeline(fns ...func(int) int)
//    that chains multiple functions together. Example:
//      transform := pipeline(double, addOne, square)
//      transform(3) → square(addOne(double(3))) = square(7) = 49
//
// ============================================================================

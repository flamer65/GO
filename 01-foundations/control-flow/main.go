// ============================================================================
// MODULE 1.3: Control Flow
// ============================================================================
//
// KEY CONCEPTS:
//   - Go has only ONE loop keyword: "for" (but it covers while/do-while/foreach)
//   - "if" can have an init statement (unique to Go!)
//   - "switch" doesn't need "break" — cases don't fall through by default
//   - No parentheses around conditions, but braces are ALWAYS required
//   - No ternary operator (? :) — Go keeps it simple
//
// RUN: go run main.go
// ============================================================================

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// -----------------------------------------------------------------------
	// PART 1: if / else
	// -----------------------------------------------------------------------

	fmt.Println("=== If/Else ===")

	age := 20

	// Basic if/else — note: NO parentheses needed, braces ARE required
	if age >= 18 {
		fmt.Println("You're an adult")
	} else if age >= 13 {
		fmt.Println("You're a teenager")
	} else {
		fmt.Println("You're a child")
	}
	// IF WITH INIT STATEMENT — This is a Go superpower!
	// The variable declared in the init statement is scoped to the if/else block
	// This is THE idiomatic way to handle errors in Go
	if remainder := 17 % 5; remainder == 0 {
		fmt.Println("17 is divisible by 5")
	} else {
		fmt.Printf("17 %% 5 = %d (not divisible)\n", remainder)
	}
	// "remainder" is NOT accessible here — it's scoped to the if block

	// The classic Go error-handling pattern (preview):
	// if err := someFunction(); err != nil {
	//     fmt.Println("Error:", err)
	//     return
	// }

	// -----------------------------------------------------------------------
	// PART 2: for loops (the ONLY loop in Go)
	// -----------------------------------------------------------------------

	fmt.Println("\n=== For Loops ===")

	// FORM 1: Classic C-style for loop
	fmt.Print("Classic:  ")
	for i := 0; i < 5; i++ {
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// FORM 2: While-style (condition only)
	fmt.Print("While:    ")
	count := 0
	for count < 5 {
		fmt.Printf("%d ", count)
		count++
	}
	fmt.Println()

	// FORM 3: Infinite loop (break to exit)
	fmt.Print("Infinite: ")
	n := 0
	for {
		if n >= 5 {
			break // Exit the loop
		}
		fmt.Printf("%d ", n)
		n++
	}
	fmt.Println()
	// FORM 4: Range-based (for iterating over slices, maps, strings, channels)
	fruits := []string{"apple", "banana", "cherry"}
	fmt.Println("Range (index + value):")
	for i, fruit := range fruits {
		fmt.Printf("  [%d] %s\n", i, fruit)
	}
	fmt.Print("Range (value only): ")
	for _, fruit := range fruits {
		fmt.Printf("%s ", fruit)
	}
	fmt.Println()
	fmt.Print("Range (index only):  ")
	for i := range fruits {
		fmt.Printf("%d ", i)
	}
	fmt.Println()
	// continue — skip current iteration
	fmt.Print("\nSkip evens: ")
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			continue // Jump to next iteration
		}
		fmt.Printf("%d ", i)
	}
	fmt.Println()
	// LABELED BREAK — break out of nested loops
	fmt.Println("\nLabeled break (find first pair summing to 10):")
     outer: // This is a label
	for i := 1; i <= 9; i++ {
		for j := 1; j <= 9; j++ {
			if i+j == 10 {
				fmt.Printf("  Found: %d + %d = 10\n", i, j)
				break outer // Breaks out of BOTH loops
			}
		}
	}

	// -----------------------------------------------------------------------
	// PART 3: switch
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Switch ===")

	// Basic switch — NO fallthrough by default (unlike C/Java)
	day := "Wednesday"
	switch day {
	case "Monday", "Tuesday", "Wednesday", "Thursday", "Friday":
		fmt.Printf("%s is a weekday\n", day)
	case "Saturday", "Sunday":
		fmt.Printf("%s is a weekend\n", day)
	default:
		fmt.Println("Unknown day")
	}

	// Switch with init statement
	switch hour := time.Now().Hour(); {
	case hour < 12:
		fmt.Println("Good morning!")
	case hour < 17:
		fmt.Println("Good afternoon!")
	default:
		fmt.Println("Good evening!")
	}

	// Switch with NO condition (acts like if/else chain, but cleaner)
	score := 85
	switch {
	case score >= 90:
		fmt.Println("Grade: A")
	case score >= 80:
		fmt.Println("Grade: B")
	case score >= 70:
		fmt.Println("Grade: C")
	default:
		fmt.Println("Grade: F")
	}

	// TYPE SWITCH — check the type of an interface value (preview of interfaces)
	var value interface{} = 42.0
	switch v := value.(type) {
	case int:
		fmt.Printf("Integer: %d\n", v)
	case float64:
		fmt.Printf("Float: %f\n", v)
	case string:
		fmt.Printf("String: %s\n", v)
	default:
		fmt.Printf("Unknown type: %T\n", v)
	}

	// -----------------------------------------------------------------------
	// PART 4: FizzBuzz (Classic Exercise — SOLVED)
	// -----------------------------------------------------------------------

	fmt.Println("\n=== FizzBuzz (1-30) ===")
	for i := 1; i <= 30; i++ {
		switch {
		case i%15 == 0:
			fmt.Printf("FizzBuzz ")
		case i%3 == 0:
			fmt.Printf("Fizz ")
		case i%5 == 0:
			fmt.Printf("Buzz ")
		default:
			fmt.Printf("%d ", i)
		}
	}
	fmt.Println()

	// -----------------------------------------------------------------------
	// PART 5: Number Guessing Game (SOLVED EXERCISE)
	// -----------------------------------------------------------------------

	fmt.Println("\n=== 🎲 Number Guessing Game ===")
	playGuessingGame()

	fmt.Println("\n=== Prime Numbers (1-100) ===")
	for nums := 1; nums <= 100; nums++ {
		if isPrime(nums) {
			fmt.Printf("%d ", nums)
		}
	}
	fmt.Println()
}

func playGuessingGame() {
	// rand.Intn(n) returns a random int in [0, n)
	target := rand.Intn(100) + 1 // Random number 1-100
	maxAttempts := 7
	fmt.Printf("I'm thinking of a number between 1 and 100.\n")
	fmt.Printf("You have %d attempts.\n\n", maxAttempts)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		fmt.Printf("Attempt %d/%d — Enter your guess: ", attempt, maxAttempts)

		var guess int
		// fmt.Scan reads from stdin
		_, err := fmt.Scan(&guess) // & gives a pointer (we'll learn this later)
		if err != nil {
			fmt.Println("Please enter a valid number!")
			continue
		}

		switch {
		case guess < target:
			fmt.Println("📈 Too low!")
		case guess > target:
			fmt.Println("📉 Too high!")
		default:
			fmt.Printf("🎉 Correct! You got it in %d attempts!\n", attempt)
			return
		}

		// Give a hint on the last attempt
		if attempt == maxAttempts-1 {
			if target%2 == 0 {
				fmt.Println("💡 Hint: The number is even")
			} else {
				fmt.Println("💡 Hint: The number is odd")
			}
		}
	}

	fmt.Printf("😢 Out of attempts! The number was %d\n", target)
}

func isPrime(n int) bool {
	// 1 is not a prime number
	if n <= 1 {
		return false
	}

	// Start checking from 2.
	// Hint optimization: we only need to check up to the square root of n (i * i <= n)
	for i := 2; i*i <= n; i++ {
		// If n is perfectly divisible by i, it's not a prime number
		if n%i == 0 {
			return false
		}
	}

	// If we checked all numbers and none divided evenly, it IS prime!
	return true
}

func fibonacci(n int) {
	// The first two numbers of the Fibonacci sequence are always 0 and 1
	a := 0
	b := 1

	for i := 0; i < n; i++ {
		fmt.Printf("%d ", a)
		// Calculate the next number by adding the two previous numbers together
		next := a + b
		
		// Shift the numbers down
		a = b
		b = next
	}
	fmt.Println()
}
func collatz(n int) int{
	fmt.Println("\n=== Collatz Conjecture ===")
	if n <= 0 {
		fmt.Println("Please enter a positive number!")
		return 0
	}
	steps := 0
	for n > 1 {
		fmt.Printf("%d ", n)
		if n%2 == 0 {
			n = n / 2
			steps++
		} else {
			n = 3*n + 1
			steps++
		}
	}
	fmt.Println(1) // Print the final 1
	fmt.Printf("Steps: %d\n", steps)
	return steps
}
func pattern(n int){
	fmt.Println("\n=== Pattern Printer ===")
   for i := 1; i <= n; i++{
	for space := 0; space < n - i; space++{
		fmt.Printf(" ")
	}	
	for star := 0; star < 2 * i - 1; star++{
		fmt.Printf("*")
	}
	fmt.Println()
   }
}
// ============================================================================
// 🏋️ EXERCISES:
//
// 1. PRIME CHECKER: Write a function isPrime(n int) bool that checks
//    if a number is prime. Then print all primes from 1 to 100.
//    Hint: only check divisors up to sqrt(n)
//
// 2. FIBONACCI: Write a function that prints the first N Fibonacci
//    numbers using a for loop (no recursion).
//
// 3. PATTERN PRINTER: Print this pattern for N=5:
//        *
//       ***
//      *****
//     *******
//    *********
//
// 4. [BONUS] COLLATZ CONJECTURE: For any positive integer n:
//    - If n is even, divide by 2
//    - If n is odd, multiply by 3 and add 1
//    - Repeat until n == 1
//    Print the sequence and count the steps for any input number.
//
// ============================================================================

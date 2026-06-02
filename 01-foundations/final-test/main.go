package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// ============================================================================
// 🏆 FINAL FOUNDATIONS TEST
// ============================================================================
// Your mission: Build a "Student Grade Analyzer"!
// 
// This challenge combines everything you've learned in Module 1:
// ✅ Variables & Custom Types
// ✅ Control Flow (Loops & Switches)
// ✅ Functions (Variadic & Closures)
// ✅ Error Handling (Errors, Nil, Panic, Recover, Defer)
// ✅ I/O (bufio.Scanner, Printf)
// ✅ Pointers (& and *)
//
// INSTRUCTIONS: Follow the TODOs below to complete the program.
// RUN: go run main.go
// ============================================================================

// TODO 1: Create a custom type called 'Grade' based on float64
type Grade float64 

// TODO 2: Create a function 'getStats' that takes a variadic number of Grades (...Grade).
// It should return three Grades: min, max, and average.
// If the length of the grades slice is 0, it should PANIC with the message: "no grades provided".
// (Hint: set an initial min to a very high number, and max to a very low number before looping)

func getStats(grades ...Grade) (Grade, Grade, Grade) {
	if len(grades) == 0 {
		panic("no grades provided")
	}
	min := Grade(10000000)
	max := Grade(-1000000)
	sum := float64(0)
	for _, g := range grades {
		if g < min {
			min = g
		}
		if g > max {
			max = g
		}
		sum += float64(g)
	}
	avg := Grade(sum / float64(len(grades)))
	return min, max, avg
	
}

// TODO 3: Create a closure factory 'makePassFilter' that takes a passing Grade (e.g., 50.0).
// It should return a function that takes a Grade and returns a boolean (true if it passes, false if not).
func makePassFilter(p Grade) func(Grade) bool {
	return func(g Grade) bool {
		return g >= p
	}
}


// TODO 4: Create a function 'addBonus' that takes a POINTER to a Grade (*Grade).
// It should add 5.0 to the grade, but cap it so it never goes above 100.0.
// (Remember to use * to dereference the pointer when reading or modifying it)
func addBonus(g *Grade){
	*g += 5.0
	if *g > 100.0 {
		*g = 100.0
	}
}


func main() {
	// TODO 5: Set up a defer and recover block.
	// If the program panics (for example, if getStats is called with no grades), 
	// catch the panic here and print a friendly message like: "Analysis aborted: <panic message>"
    defer func(){
		if r := recover(); r != nil {
			fmt.Printf("Analysis aborted: %v\n", r)
		}
	}()


	fmt.Println("=== 🎓 Grade Analyzer ===")
	fmt.Println("Enter grades one by one. Type 'done' when finished.")

	// Create an empty slice to hold the grades
	 var grades []Grade
     
	// TODO 6: Use bufio.Scanner to read input from the user in a loop.
	// - If the user types "done", break the loop.
	// - Use strconv.ParseFloat to convert the input.
	// - If there's an ERROR (e.g. user types "apple"), print a polite error message ("Invalid input, skipping...") and 'continue' the loop.
	// - If it's valid, convert the float64 to your custom 'Grade' type and append it to your grades slice.
	//
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		if input == "done" {
			break
		}
		grade, err := strconv.ParseFloat(input, 64)
		if err != nil {
			fmt.Println("Invalid input, skipping...")
			continue
		}
		grades = append(grades, Grade(grade))
		}



	// TODO 7: Call getStats with your slice of grades.
	// (Hint: use the ... operator to unpack the slice into the variadic function!)
	min, max, avg := getStats(grades...)



	// TODO 8: Use fmt.Printf to print the Min, Max, and Average formatted to 1 decimal place.
    fmt.Printf("Min: %.1f\nMax: %.1f\nAvg: %.1f\n", min, max, avg)


	// TODO 9: Process the final grades!
	// - Create a passing filter for a grade of 50.0 using makePassFilter
	// - Loop through all the grades in your slice
	// - Pass each grade's ADDRESS (&) to addBonus
	// - Check if they passed using your filter
	// - Print the final grade and whether they passed or failed!
	fmt.Println("\n=== Final Results ===")
	isPassing := makePassFilter(50.0)

	for i := range grades {
		// Pass the address to addBonus
		addBonus(&grades[i])

		// Check if passed
		passed := isPassing(grades[i])

		status := "❌ Failed"
		if passed {
			status = "✅ Passed"
		}

		fmt.Printf("Student %d: %.1f - %s\n", i+1, grades[i], status)
	}
}

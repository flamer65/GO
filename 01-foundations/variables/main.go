// ============================================================================
// MODULE 1.2: Variables, Types & Constants
// ============================================================================
//
// KEY CONCEPTS:
//   - Go is STATICALLY TYPED — every variable has a fixed type at compile time
//   - Go has NO implicit type conversion — you must convert explicitly
//   - "Zero values" — Go initializes all variables to a predictable default
//   - Two ways to declare: "var" (explicit) and ":=" (short, type inferred)
//   - Constants are computed at compile time, not runtime
//
// RUN: go run main.go
// ============================================================================

package main

import (
	"fmt"
	"strconv"
)

func main() {
	// -----------------------------------------------------------------------
	// PART 1: Variable Declarations
	// -----------------------------------------------------------------------

	// METHOD 1: "var" keyword — explicit declaration
	var age int         // Declared but not assigned → gets "zero value" (0)
	var name string     // Zero value for string is "" (empty string)
	var isReady bool    // Zero value for bool is false
	var balance float64 // Zero value for float64 is 0.0

	fmt.Println("=== Zero Values (super important in Go!) ===")
	fmt.Printf("int:     %d\n", age)
	fmt.Printf("string:  %q\n", name) // %q adds quotes around strings
	fmt.Printf("bool:    %t\n", isReady)
	fmt.Printf("float64: %f\n", balance)

	// METHOD 2: "var" with initialization
	var city string = "Delhi"
	var temp float64 = 42.5

	// METHOD 3: "var" with type inference (type deduced from value)
	var country = "India" // Go infers this is a string

	// METHOD 4: Short declaration with := (MOST COMMON inside functions)
	// Only works inside functions, not at package level!
	language := "Go"  // Inferred as string
	year := 2009      // Inferred as int
	pi := 3.14159     // Inferred as float64
	isAwesome := true // Inferred as bool

	fmt.Println("\n=== Declared Variables ===")
	fmt.Printf("City: %s, Temp: %.1f°C\n", city, temp)
	fmt.Printf("Country: %s\n", country)
	fmt.Printf("%s was created in %d. Awesome? %t. Pi: %.5f\n", language, year, isAwesome, pi)

	// METHOD 5: Multiple declarations in a block
	var (
		firstName = "Naman"
		lastName  = "Garg"
		score     = 100
	)
	fmt.Printf("\n%s %s scored %d\n", firstName, lastName, score)

	// METHOD 6: Multiple short declarations on one line
	x, y, z := 1, 2.5, "three"
	fmt.Printf("x=%d, y=%.1f, z=%s\n", x, y, z)

	// -----------------------------------------------------------------------
	// PART 2: Basic Types
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Go's Basic Types ===")

	// INTEGERS — signed and unsigned, various sizes
	var i8 int8 = 127     // -128 to 127
	var i16 int16 = 32767 // -32768 to 32767
	var i32 int32 = 2147483647
	var i64 int64 = 9223372036854775807
	var u8 uint8 = 255       // 0 to 255 (same as "byte")
	var platformInt int = 42 // Size depends on platform (32 or 64 bit)

	fmt.Printf("int8: %d, int16: %d, int32: %d\n", i8, i16, i32)
	fmt.Printf("int64: %d, uint8: %d, int: %d\n", i64, u8, platformInt)

	// FLOATING POINT
	var f32 float32 = 3.14
	var f64 float64 = 3.141592653589793 // Prefer float64 (more precise)
	fmt.Printf("float32: %f, float64: %.15f\n", f32, f64)

	// BYTE and RUNE — type aliases with special meaning
	var b byte = 'A' // byte = uint8, used for ASCII characters
	var r rune = '🚀' // rune = int32, used for Unicode code points
	fmt.Printf("byte: %c (value: %d)\n", b, b)
	fmt.Printf("rune: %c (value: %d, Unicode: U+%04X)\n", r, r, r)

	// STRING — immutable sequence of bytes (usually UTF-8)
	greeting := "Hello, 世界" // Go strings are UTF-8 by default!
	fmt.Printf("string: %s (length in bytes: %d)\n", greeting, len(greeting))

	// -----------------------------------------------------------------------
	// PART 3: Type Conversions (Go NEVER converts implicitly!)
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Type Conversions ===")

	intVal := 42
	floatVal := float64(intVal) // int → float64
	fmt.Printf("int %d → float64 %f\n", intVal, floatVal)

	bigFloat := 3.99
	truncated := int(bigFloat) // float64 → int (TRUNCATES, doesn't round!)
	fmt.Printf("float64 %f → int %d (truncated, not rounded!)\n", bigFloat, truncated)

	// This would NOT compile — Go refuses implicit conversion:
	// var myInt int = 42
	// var myFloat float64 = myInt  // ❌ COMPILE ERROR

	// -----------------------------------------------------------------------
	// PART 4: Constants
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Constants ===")

	// Constants are immutable values known at compile time
	const speedOfLight = 299792458 // meters per second
	const appName = "GoLearner"
	const version = 1.0

	fmt.Printf("%s v%.1f\n", appName, version)
	fmt.Printf("Speed of light: %d m/s\n", speedOfLight)

	// IOTA — Go's auto-incrementing constant generator
	// Starts at 0 and increments by 1 for each constant in the block

	// Simple enum-like pattern
	const (
		Sunday    = iota // 0
		Monday           // 1 (iota auto-increments)
		Tuesday          // 2
		Wednesday        // 3
		Thursday         // 4
		Friday           // 5
		Saturday         // 6
	)

	fmt.Printf("\nDays: Sun=%d, Mon=%d, Fri=%d, Sat=%d\n",
		Sunday, Monday, Friday, Saturday)

	// IOTA with expressions — powers of 2 (useful for bit flags)
	const (
		ReadPermission    = 1 << iota // 1 << 0 = 1
		WritePermission               // 1 << 1 = 2
		ExecutePermission             // 1 << 2 = 4
	)

	fmt.Printf("Permissions: Read=%d, Write=%d, Execute=%d\n",
		ReadPermission, WritePermission, ExecutePermission)

	// Combine with bitwise OR
	myPermissions := ReadPermission | WritePermission // 1 | 2 = 3
	fmt.Printf("My permissions (read+write): %d\n", myPermissions)

	// IOTA for byte sizes
	const (
		_  = iota             // Skip 0
		KB = 1 << (10 * iota) // 1 << 10 = 1024
		MB                    // 1 << 20 = 1,048,576
		GB                    // 1 << 30 = 1,073,741,824
	)

	fmt.Printf("\nByte sizes: KB=%d, MB=%d, GB=%d\n", KB, MB, GB)

	// -----------------------------------------------------------------------
	// PART 5: Temperature Converter Exercise (SOLVED EXAMPLE)
	// -----------------------------------------------------------------------

	fmt.Println("\n=== 🌡️ Temperature Converter ===")

	celsius := 100.0
	fahrenheit := celsiusToFahrenheit(Celsius(celsius))
	kelvin := celsiusToKelvin(Celsius(celsius))

	fmt.Printf("%.1f°C = %.1f°F = %.1fK\n", celsius, fahrenheit, kelvin)

	// Convert back
	fmt.Printf("%.1f°F = %.1f°C\n", fahrenheit, fahrenheitToCelsius(fahrenheit))
	fmt.Printf("%.1fK  = %.1f°C\n", kelvin, kelvinToCelsius(kelvin))

	// Fun temperatures
	temperatures := []float64{-40, 0, 37, 100, 1000}
	fmt.Println("\nConversion Table:")
	fmt.Println("  Celsius  │ Fahrenheit │  Kelvin")
	fmt.Println("───────────┼────────────┼──────────")
	for _, c := range temperatures {
		fmt.Printf("  %7.1f°C │  %7.1f°F  │ %7.1fK\n",
			c, celsiusToFahrenheit(Celsius(c)), celsiusToKelvin(Celsius(c)))
	}
	readTemperature()

}

// -----------------------------------------------------------------------
// Temperature conversion functions
// Notice: functions outside main are accessible within the same package
// -----------------------------------------------------------------------
type Celsius float64
type Fahrenheit float64
type Kelvin float64

func celsiusToFahrenheit(c Celsius) Fahrenheit {
	return Fahrenheit(c*9.0/5.0 + 32)
}

func fahrenheitToCelsius(f Fahrenheit) Celsius {
	return Celsius((f - 32) * 5.0 / 9.0)
}

func celsiusToKelvin(c Celsius) Kelvin {
	return Kelvin(c + 273.15)
}

func kelvinToCelsius(k Kelvin) Celsius {
	return Celsius(k - 273.15)
}

func fahrenheitToKelvin(f Fahrenheit) Kelvin {
	return Kelvin((f-32)*5.0/9.0 + 273.15)
}

func kelvinToFahrenheit(k Kelvin) Fahrenheit {
	return Fahrenheit((k-273.15)*5.0/9.0 + 32)
}
func readTemperature() {
	var temp string
	fmt.Print("Enter temperature (e.g. 100C): ")
	fmt.Scan(&temp)
	unit := temp[len(temp)-1:]

	numStr := temp[:len(temp)-1]

	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		fmt.Println("Invalid number")
		return
	}
	switch unit {
	case "C", "c":
		c := Celsius(num)
		fmt.Printf("%.1f°C = %.1f°F = %.1fK\n", c, celsiusToFahrenheit(c), celsiusToKelvin(c))
	case "F", "f":
		f := Fahrenheit(num)
		fmt.Printf("%.1f°F = %.1f°C = %.1fK\n", f, fahrenheitToCelsius(f), fahrenheitToKelvin(f))
	case "K", "k":
		k := Kelvin(num)
		fmt.Printf("%.1fK = %.1f°C = %.1f°F\n", k, kelvinToCelsius(k), kelvinToFahrenheit(k))
	default:
		fmt.Println("Invalid unit")
	}

}

// ============================================================================
// 🏋️ EXERCISES:
//
// 1. Add a function fahrenheitToKelvin() and kelvinToFahrenheit()
//    that compose the existing functions
//
// 2. Make this interactive: read a temperature from the user using
//    fmt.Scan() or bufio.Scanner, detect the unit (C/F/K suffix),
//    and convert to all three
//
// 3. Create a custom type:
//      type Celsius float64
//      type Fahrenheit float64
//    Then add methods on these types for conversion.
//    This previews Module 2's struct methods!
//
// ============================================================================

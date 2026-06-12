// ============================================================================
// MODULE 3.0: Interfaces — The Core Concept
// ============================================================================
//
// KEY CONCEPTS:
//   - An interface is a TYPE that defines a set of METHODS
//   - Any struct that has those methods AUTOMATICALLY satisfies the interface
//   - No "implements" keyword — it's implicit
//   - Interfaces let you write functions that work with MANY types
//   - The empty interface (any) accepts everything
//   - Type assertions let you get the real type back from an interface
//
// This lesson builds slowly from "what is an interface?" to "why is it
// useful?" using examples you can follow step by step.
//
// RUN: go run main.go
// ============================================================================

package main

import (
	"fmt"
	"math"
)

func main() {
	// -----------------------------------------------------------------------
	// PART 1: The Problem — Why Do We Need Interfaces?
	// -----------------------------------------------------------------------

	fmt.Println("=== PART 1: The Problem ===")

	// Let's say we have two shapes:
	c := Circle{Radius: 5}
	r := Rectangle{Width: 10, Height: 3}

	// We can print their areas using separate functions:
	printCircleArea(c)
	printRectangleArea(r)

	// But what if we want ONE function that works with BOTH?
	// We can't do this:
	//   printArea(c)  // ❌ what type does printArea accept? Circle? Rectangle?
	//
	// We need a way to say:
	//   "Accept ANYTHING that has an Area() method"
	//
	// That's what interfaces do! ↓↓↓

	// -----------------------------------------------------------------------
	// PART 2: Defining an Interface
	// -----------------------------------------------------------------------

	fmt.Println("\n=== PART 2: Defining an Interface ===")

	// An interface is defined with the `type` keyword, just like a struct.
	// But instead of fields, it lists METHODS:
	//
	//   type Shape interface {
	//       Area() float64
	//   }
	//
	// This says: "A Shape is ANYTHING that has an Area() method
	// which takes no arguments and returns a float64."
	//
	// See the Shape interface defined below main().

	// Circle has an Area() method → Circle IS a Shape ✅
	// Rectangle has an Area() method → Rectangle IS a Shape ✅
	// No special keyword needed. Just having the method is enough!

	// Now we can use ONE function for both:
	printArea(c) // c is a Circle, but it has Area() → it's a Shape ✅
	printArea(r) // r is a Rectangle, but it has Area() → it's a Shape ✅

	// -----------------------------------------------------------------------
	// PART 3: How Does Go Know Which Area() to Call?
	// -----------------------------------------------------------------------

	fmt.Println("\n=== PART 3: Which Area() Gets Called? ===")

	// When you call s.Area() on an interface variable, Go looks at
	// what ACTUAL type is stored inside, and calls THAT type's method.

	var s Shape // s is an interface variable — can hold any Shape

	s = Circle{Radius: 5}
	fmt.Printf("s holds a Circle  → s.Area() = %.2f\n", s.Area())
	// Go sees: s is actually a Circle → calls Circle's Area()
	// → 3.14159 * 5 * 5 = 78.54

	s = Rectangle{Width: 10, Height: 3}
	fmt.Printf("s holds a Rect    → s.Area() = %.2f\n", s.Area())
	// Go sees: s is actually a Rectangle → calls Rectangle's Area()
	// → 10 * 3 = 30.00

	// Same variable `s`, same method call `s.Area()`,
	// but DIFFERENT code runs depending on what's inside!

	// -----------------------------------------------------------------------
	// PART 4: Interfaces Let You Make Collections of Different Types
	// -----------------------------------------------------------------------

	fmt.Println("\n=== PART 4: Collections of Different Types ===")

	// Without interfaces, you CAN'T put Circle and Rectangle in the same slice:
	//   mixed := []??? { Circle{5}, Rectangle{3,4} }  // ❌ what type is ???

	// With interfaces, you can!
	shapes := []Shape{
		Circle{Radius: 5},
		Rectangle{Width: 10, Height: 3},
		Circle{Radius: 2},
		Rectangle{Width: 7, Height: 7},
		Triangle{Base: 5, Height:7, SideA: 7, SideB: 8, SideC: 9},
	}

	// Now loop over ALL shapes — doesn't matter what they actually are
	fmt.Println("All shapes:")
	totalArea := 0.0
	for i, shape := range shapes {
		area := shape.Area()
		totalArea += area
		fmt.Printf("  [%d] %T → area = %.2f\n", i, shape, area)
		//         %T prints the actual type (Circle or Rectangle)
	}
	fmt.Printf("Total area: %.2f\n", totalArea)

	// -----------------------------------------------------------------------
	// PART 5: An Interface Can Require Multiple Methods
	// -----------------------------------------------------------------------

	fmt.Println("\n=== PART 5: Multiple Methods ===")

	// Our Shape interface only requires Area().
	// Let's make a STRICTER interface that requires MORE methods:
	//
	//   type FullShape interface {
	//       Area() float64
	//       Perimeter() float64
	//       Describe() string
	//   }
	//
	// Now a type must have ALL THREE methods to qualify.

	// Circle has Area(), Perimeter(), AND Describe() → FullShape ✅
	// Rectangle has Area(), Perimeter(), AND Describe() → FullShape ✅
	printFullInfo(Circle{Radius: 5})
	printFullInfo(Rectangle{Width: 10, Height: 3})
    printFullInfo(Triangle{Base: 5, Height:7, SideA: 7, SideB: 8, SideC: 9})
	// What if a type only has Area() but not Perimeter()?
	// → It's a Shape, but NOT a FullShape.
	// The compiler catches this for you!

	// -----------------------------------------------------------------------
	// PART 6: Value Receivers vs Pointer Receivers with Interfaces
	// -----------------------------------------------------------------------

	fmt.Println("\n=== PART 6: Value vs Pointer Receivers ===")

	// RULE:
	//   Value receiver  → both values AND pointers satisfy the interface
	//   Pointer receiver → ONLY pointers satisfy the interface

	// All our methods so far use value receivers:
	//   func (c Circle) Area() float64 { ... }    ← value receiver
	// So both work:
	var s1 Shape = Circle{5}   // ✅ value — works
	var s2 Shape = &Circle{5}  // ✅ pointer — also works
	fmt.Printf("Value:   %.2f\n", s1.Area())
	fmt.Printf("Pointer: %.2f\n", s2.Area())

	// Now look at Counter — it uses a POINTER receiver:
	//   func (c *Counter) Increment() { ... }    ← pointer receiver
	// Because it needs to MODIFY the struct.

	// var ct Incrementer = Counter{0}   // ❌ WON'T COMPILE
	var ct Incrementer = &Counter{0}     // ✅ must use pointer
	ct.Increment()
	ct.Increment()
	ct.Increment()
	fmt.Printf("Counter after 3 increments: %d\n", ct.GetCount())

	// WHY? Because a value receiver gets a COPY — modifications are lost.
	// A pointer receiver modifies the ORIGINAL. If you say the interface
	// needs modification, Go forces you to use a pointer.

	// -----------------------------------------------------------------------
	// PART 7: The Empty Interface — any
	// -----------------------------------------------------------------------

	fmt.Println("\n=== PART 7: The Empty Interface (any) ===")

	// An interface with ZERO methods:
	//   type any = interface{}
	//
	// Zero methods required → EVERY type satisfies it → accepts anything!

	printAnything(42)
	printAnything("hello")
	printAnything(3.14)
	printAnything(true)
	printAnything(Circle{Radius: 5})
	printAnything([]int{1, 2, 3})

	// Useful for maps with mixed value types:
	person := map[string]any{
		"name": "Naman",
		"age":  23,
		"pro":  true,
	}
	fmt.Println("\nMixed map:", person)

	// But you LOSE type safety — you can't call .Area() on an `any`:
	//   var x any = Circle{5}
	//   x.Area()  // ❌ COMPILE ERROR: any has no method Area

	// To use the real type, you need a type assertion (Part 8) ↓

	// -----------------------------------------------------------------------
	// PART 8: Type Assertions — Getting the Real Type Back
	// -----------------------------------------------------------------------

	fmt.Println("\n=== PART 8: Type Assertions ===")

	// If you have an interface variable and want the concrete type back:
	var myShape Shape = Circle{Radius: 10}

	// SAFE way — use the comma-ok pattern (just like maps!)
	circle, ok := myShape.(Circle)
	if ok {
		fmt.Printf("It's a Circle! Radius = %.1f\n", circle.Radius)
	}

	rect, ok := myShape.(Rectangle)
	if ok {
		fmt.Printf("It's a Rectangle! Width = %.1f\n", rect.Width)
	} else {
		fmt.Println("It's NOT a Rectangle")
	}

	// UNSAFE way — panics if wrong type (avoid this usually)
	// circle2 := myShape.(Rectangle)  // 💥 PANIC at runtime

	// -----------------------------------------------------------------------
	// PART 9: Type Switch — Check Multiple Types at Once
	// -----------------------------------------------------------------------

	fmt.Println("\n=== PART 9: Type Switch ===")

	// A type switch is like a regular switch, but checks the TYPE:
	testValues := []any{42, "hello", 3.14, true, Circle{5}, Rectangle{3, 4}}

	for _, v := range testValues {
		whatIsIt(v)
	}

	// -----------------------------------------------------------------------
	// PART 10: Real-World Example — A Notification System
	// -----------------------------------------------------------------------

	fmt.Println("\n=== PART 10: Notification System ===")

	// This is WHY interfaces matter in real code.
	// You write your logic ONCE, and swap implementations easily.

	// Define the interface:  (see Notifier below)
	//   type Notifier interface {
	//       Notify(to, message string)
	//   }

	// Two different implementations:
	email := EmailNotifier{Sender: "app@golearn.dev"}
	sms := SMSNotifier{FromNumber: "+91-9999999999"}
	console := ConsoleNotifier{}

	// Same function works with ALL of them:
	sendAlert(email, "Naman", "Your order shipped!")
	sendAlert(sms, "Naman", "Your OTP is 482913")
	sendAlert(console, "Naman", "Server is healthy")

	// Tomorrow you can add PushNotifier, SlackNotifier, WhatsAppNotifier...
	// WITHOUT changing sendAlert at all. That's the power of interfaces.

	// -----------------------------------------------------------------------
	// PART 11: Interface Satisfaction — A Mental Model
	// -----------------------------------------------------------------------

	fmt.Println("\n=== PART 11: Summary ===")

	fmt.Println(`
  ┌─────────────────────────────────────────────────────────────┐
  │                    HOW INTERFACES WORK                      │
  ├─────────────────────────────────────────────────────────────┤
  │                                                             │
  │  1. Define the CONTRACT (what methods are needed):          │
  │     type Shape interface { Area() float64 }                 │
  │                                                             │
  │  2. ANY struct with those methods AUTOMATICALLY qualifies:  │
  │     Circle has Area()    → Circle is a Shape ✅              │
  │     Rectangle has Area() → Rectangle is a Shape ✅           │
  │     string has no Area() → string is NOT a Shape ❌          │
  │                                                             │
  │  3. Use the interface as a type:                            │
  │     func process(s Shape) { s.Area() }                      │
  │     process(Circle{5})      ← calls Circle's Area           │
  │     process(Rectangle{3,4}) ← calls Rectangle's Area        │
  │                                                             │
  │  4. Go picks the right method AT RUNTIME based on           │
  │     what's actually stored inside the interface variable.   │
  │                                                             │
  │  REMEMBER:                                                  │
  │    • Struct  = what something IS     (data + methods)       │
  │    • Interface = what something CAN DO (just methods)       │
  │    • No "implements" keyword — just have the methods!       │
  │                                                             │
  └─────────────────────────────────────────────────────────────┘
`)
}

// =========================================================================
// TYPE DEFINITIONS
// =========================================================================

// ---- Shape Interface ----
// Any type with Area() is a Shape
type Shape interface {
	Area() float64
}

// ---- FullShape Interface ----
// Requires more methods — stricter contract
type FullShape interface {
	Area() float64
	Perimeter() float64
	Describe() string
}

// ---- Circle ----
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func (c Circle) Describe() string {
	return fmt.Sprintf("Circle with radius %.1f", c.Radius)
}

// ---- Rectangle ----
type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func (r Rectangle) Describe() string {
	return fmt.Sprintf("Rectangle %.1f × %.1f", r.Width, r.Height)
}

// ---- Counter (pointer receiver example) ----
type Incrementer interface {
	Increment()
	GetCount() int
}

type Counter struct {
	count int
}

// Pointer receiver — because we MODIFY count
func (c *Counter) Increment() {
	c.count++
}

func (c *Counter) GetCount() int {
	return c.count
}

// ---- Notifier (real-world example) ----
type Notifier interface {
	Notify(to, message string)
}

type EmailNotifier struct {
	Sender string
}

func (e EmailNotifier) Notify(to, message string) {
	fmt.Printf("  📧 Email from %s → %s: %s\n", e.Sender, to, message)
}

type SMSNotifier struct {
	FromNumber string
}

func (s SMSNotifier) Notify(to, message string) {
	fmt.Printf("  📱 SMS from %s → %s: %s\n", s.FromNumber, to, message)
}

type ConsoleNotifier struct{}

func (c ConsoleNotifier) Notify(to, message string) {
	fmt.Printf("  🖥️  Console → %s: %s\n", to, message)
}

// =========================================================================
// FUNCTIONS
// =========================================================================

// ---- The old way: one function per type ----
func printCircleArea(c Circle) {
	fmt.Printf("  Circle area: %.2f\n", c.Area())
}

func printRectangleArea(r Rectangle) {
	fmt.Printf("  Rectangle area: %.2f\n", r.Area())
}

// ---- The interface way: one function for ALL shapes ----
func printArea(s Shape) {
	fmt.Printf("  Shape area: %.2f\n", s.Area())
}

// ---- Using the stricter FullShape interface ----
func printFullInfo(fs FullShape) {
	fmt.Printf("  %s → area=%.2f, perimeter=%.2f\n",
		fs.Describe(), fs.Area(), fs.Perimeter())
}

// ---- Empty interface: accepts anything ----
func printAnything(v any) {
	fmt.Printf("  any: %-20v (type: %T)\n", v, v)
}

// ---- Type switch ----
func whatIsIt(v any) {
	switch val := v.(type) {
	case int:
		fmt.Printf("  int → %d (doubled: %d)\n", val, val*2)
	case string:
		fmt.Printf("  string → %q (length: %d)\n", val, len(val))
	case float64:
		fmt.Printf("  float64 → %.2f\n", val)
	case bool:
		fmt.Printf("  bool → %t\n", val)
	case Circle:
		fmt.Printf("  Circle → radius=%.1f, area=%.2f\n", val.Radius, val.Area())
	case Rectangle:
		fmt.Printf("  Rectangle → %.1f×%.1f, area=%.2f\n", val.Width, val.Height, val.Area())
	default:
		fmt.Printf("  unknown → %v\n", val)
	}
}

// ---- Notifier function — works with ANY Notifier implementation ----
func sendAlert(n Notifier, to, message string) {
	n.Notify(to, message)
}
type Triangle struct {
	Base float64
	Height float64
	SideA float64
	SideB float64
	SideC float64
}
func (t Triangle) Area() float64{
 return 0.5 * t.Height * t.Base
}
func (t Triangle) Perimeter() float64{
	return t.SideA + t.SideB + t.SideC
}
func (t Triangle) Describe() string{
	return fmt.Sprintf("Triangle with Base:%.1f Height:%.1f SideA:%.1f SideB:%1f SideC:%1f", t.Base, t.Height, t.SideA, t.SideB, t.SideC)
}
type Measurable interface {
	Value() float64
}

type Temperature struct {
	Celsius float64
}
func (t Temperature) Value() float64 { return t.Celsius }

type Distance struct {
	Meters float64
}
func (d Distance) Value() float64 { return d.Meters }

type Weight struct {
	Kilograms float64
}
func (w Weight) Value() float64 { return w.Kilograms }

func average(items []Measurable) float64 {
	if len(items) == 0 {
		return 0
	}
	sum := 0.0
	for _, item := range items {
		sum += item.Value()
	}
	return sum / float64(len(items))
}

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}
type ConsoleLogger struct{}
func (c ConsoleLogger) Info(msg string)  {
	 fmt.Printf("[Info] %s\n",msg)
}
func (c ConsoleLogger) Warn(msg string)  {
	fmt.Printf("[Warn] %s\n", msg)
}
func (c ConsoleLogger) Error(msg string)  {
	fmt.Printf("[Error] %s\n", msg)
}


type MemoryLogger struct{
	logs []string
}
func (m *MemoryLogger) Info(msg string)  { m.logs = append(m.logs, "[INFO] "+msg) }
func (m *MemoryLogger) Warn(msg string)  { m.logs = append(m.logs, "[WARN] "+msg) }
func (m *MemoryLogger) Error(msg string) { m.logs = append(m.logs, "[ERROR] "+msg) }
// ============================================================================
// 🏋️ EXERCISES:
//
// 1. Add a Triangle struct with Base, Height, SideA, SideB, SideC fields.
//    Make it satisfy BOTH Shape AND FullShape interfaces.
//    Add it to the shapes slice in Part 4 and verify it works.
//
// 2. Create an Animal interface with Speak() string.
//    Implement Dog, Cat, and Duck. Write a function chorus(animals []Animal)
//    that calls Speak() on each and prints the result.
//    Notice: you didn't write "Dog implements Animal" anywhere!
//
// 3. Create a Measurable interface with Value() float64.
//    Implement it for Temperature, Distance, and Weight structs.
//    Write a generic average(items []Measurable) float64 function.
//
// 4. [BONUS] Create a Logger interface: Info(msg), Warn(msg), Error(msg).
//    Implement ConsoleLogger (prints to screen) and MemoryLogger (stores
//    in a []string slice). Write code that works with both — then swap
//    the implementation and watch nothing else change.
//
// ============================================================================

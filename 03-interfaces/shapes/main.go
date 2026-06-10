// ============================================================================
// MODULE 3.1: Interfaces — The Power of Implicit Contracts
// ============================================================================
//
// KEY CONCEPTS:
//   - An interface defines a SET OF METHODS — any type that implements
//     those methods automatically satisfies the interface (no "implements")
//   - Interfaces are IMPLICIT — no declaration needed on the implementing type
//   - The empty interface (interface{} or "any") accepts any value
//   - Type assertions let you extract the concrete type from an interface
//   - Type switches are elegant multi-type dispatch
//   - Small interfaces are powerful: io.Reader, io.Writer, fmt.Stringer
//   - "Accept interfaces, return structs" — a key Go design principle
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
	// PART 1: Interface Basics — What and Why
	// -----------------------------------------------------------------------

	fmt.Println("=== Interface Basics ===")

	// An interface is a contract: "any type with these methods"
	// See the Shape interface defined below

	// Circle and Rectangle both implement Shape — no explicit declaration!
	c := Circle{Radius: 5}
	r := Rectangle{Width: 10, Height: 3}
	t := Triangle{Base: 8, Height: 6, SideA: 8, SideB: 6, SideC: 10}

	// We can use them as Shape — polymorphism!
	shapes := []Shape{c, r, t}

	fmt.Println("All shapes:")
	for _, s := range shapes {
		printShapeInfo(s)
	}

	// The power: printShapeInfo works with ANY type that has Area() and Perimeter()
	// You can add new shapes without modifying existing code!

	// -----------------------------------------------------------------------
	// PART 2: Implicit Satisfaction — No "implements" Keyword
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Implicit Satisfaction ===")

	// In Java/C#: class Circle implements Shape { ... }
	// In Go: Circle just HAS the methods — that's it!

	// This means you can satisfy interfaces from OTHER packages
	// without importing them. This is incredibly powerful.

	// Example: fmt.Stringer — any type with String() string
	// See our Circle.String() method below
	fmt.Println("Circle as Stringer:", c) // fmt.Println calls String() automatically!
	fmt.Println("Rectangle as Stringer:", r)
	fmt.Println("Triangle as Stringer:", t)

	// -----------------------------------------------------------------------
	// PART 3: Interface Values — What's Inside?
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Interface Values ===")

	// An interface value has TWO components:
	//   1. A type pointer (what concrete type is stored)
	//   2. A value pointer (the actual data)

	var s Shape
	fmt.Printf("Nil interface: value=%v, type=%T\n", s, s)

	s = Circle{Radius: 3}
	fmt.Printf("Circle:    value=%v, type=%T\n", s, s)

	s = Rectangle{Width: 4, Height: 5}
	fmt.Printf("Rectangle: value=%v, type=%T\n", s, s)

	// GOTCHA: An interface holding a nil pointer is NOT a nil interface!
	var cp *Circle // nil pointer
	var s2 Shape = cp
	fmt.Printf("\nNil pointer in interface: value=%v, type=%T, isNil=%t\n",
		s2, s2, s2 == nil) // isNil=false! The interface HOLDS a type.

	// -----------------------------------------------------------------------
	// PART 4: The Empty Interface — any / interface{}
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Empty Interface (any) ===")

	// interface{} (or "any" since Go 1.18) has zero methods
	// Therefore EVERY type satisfies it — it can hold anything

	var anything any // same as interface{}

	anything = 42
	fmt.Printf("int:    %v (type: %T)\n", anything, anything)

	anything = "hello"
	fmt.Printf("string: %v (type: %T)\n", anything, anything)

	anything = Circle{Radius: 7}
	fmt.Printf("Circle: %v (type: %T)\n", anything, anything)

	anything = []int{1, 2, 3}
	fmt.Printf("slice:  %v (type: %T)\n", anything, anything)

	// Common use: functions that accept anything
	printAnything(42)
	printAnything("Go is awesome")
	printAnything(true)
	printAnything(Circle{Radius: 5})

	// Maps with mixed values
	metadata := map[string]any{
		"name":    "Naman",
		"age":     23,
		"scores":  []int{95, 87, 92},
		"premium": true,
	}
	fmt.Println("\nMixed map:", metadata)

	// -----------------------------------------------------------------------
	// PART 5: Type Assertions
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Type Assertions ===")

	// Type assertion extracts the concrete value from an interface
	// Syntax: value := interfaceVar.(ConcreteType)

	var shape Shape = Circle{Radius: 10}

	// UNSAFE assertion — panics if wrong type
	circle := shape.(Circle)
	fmt.Printf("Asserted Circle: radius=%f\n", circle.Radius)

	// SAFE assertion — returns (value, ok)
	circle2, ok := shape.(Circle)
	fmt.Printf("Safe Circle: radius=%f, ok=%t\n", circle2.Radius, ok)

	rect, ok := shape.(Rectangle)
	fmt.Printf("Safe Rectangle: %v, ok=%t\n", rect, ok) // ok=false

	// Common pattern: check and use
	if c, ok := shape.(Circle); ok {
		fmt.Printf("It's a circle with area %.2f\n", c.Area())
	}

	// -----------------------------------------------------------------------
	// PART 6: Type Switches
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Type Switches ===")

	// Type switch — elegant multi-type dispatch
	testValues := []any{
		42,
		"hello",
		3.14,
		true,
		Circle{Radius: 5},
		Rectangle{Width: 3, Height: 4},
		[]int{1, 2, 3},
		nil,
	}

	for _, v := range testValues {
		describeType(v)
	}

	// -----------------------------------------------------------------------
	// PART 7: Common Standard Library Interfaces
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Standard Library Interfaces ===")

	// fmt.Stringer — already used above
	// type Stringer interface { String() string }

	// sort.Interface — for custom sorting
	team := Team{
		{Name: "Alice", Score: 95},
		{Name: "Charlie", Score: 78},
		{Name: "Bob", Score: 87},
		{Name: "Diana", Score: 92},
	}

	fmt.Println("Before sorting:", team)
	// We implement sort.Interface on Team (see below)
	// sort.Sort(team) — but let's do it manually to show the interface
	bubbleSort(team)
	fmt.Println("After sorting: ", team)

	// io.Reader and io.Writer — the most important interfaces in Go
	// type Reader interface { Read(p []byte) (n int, err error) }
	// type Writer interface { Write(p []byte) (n int, err error) }
	// These will be covered deeply in Module 6 (stdlib)

	// error — the simplest and most used interface
	// type error interface { Error() string }
	// Covered in detail in Module 3.2!

	// -----------------------------------------------------------------------
	// PART 8: Interface Composition
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Interface Composition ===")

	// Interfaces can embed other interfaces — composition!
	// See ReadWriter, Measurable, Drawable interfaces below

	circle3 := NewDrawableCircle(5, "red", true)

	// DrawableCircle satisfies both Shape and Drawable
	// Therefore it satisfies ShapeDrawable (the composed interface)
	demonstrateDrawable(circle3)

	// -----------------------------------------------------------------------
	// PART 9: Practical Patterns
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Practical Patterns ===")

	// Pattern 1: Strategy Pattern — swap behavior via interfaces
	fmt.Println("--- Strategy Pattern ---")

	prices := []float64{29.99, 149.50, 9.99, 89.00, 210.00}

	noDiscount := &NoDiscount{}
	percentage := &PercentageDiscount{Percent: 20}
	bulkDisc := &BulkDiscount{Threshold: 3, DiscountPercent: 15}

	fmt.Printf("  No discount:  $%.2f\n", calculateTotal(prices, noDiscount))
	fmt.Printf("  20%% off:      $%.2f\n", calculateTotal(prices, percentage))
	fmt.Printf("  Bulk (3+):    $%.2f\n", calculateTotal(prices, bulkDisc))

	// Pattern 2: Dependency Injection via interfaces
	fmt.Println("\n--- Notification Service ---")

	emailSender := &EmailNotifier{From: "app@golearn.dev"}
	smsSender := &SMSNotifier{FromNumber: "+1-555-0100"}

	// Same service, different notifiers — swappable!
	service1 := NotificationService{Notifier: emailSender}
	service2 := NotificationService{Notifier: smsSender}

	service1.Alert("Alice", "Your order has shipped!")
	service2.Alert("Bob", "Your OTP is 482913")

	// -----------------------------------------------------------------------
	// PART 10: Practical Exercise — Plugin System (SOLVED)
	// -----------------------------------------------------------------------

	fmt.Println("\n=== 🔌 Plugin System ===")

	// A text processing pipeline using interfaces
	pipeline := Pipeline{
		Plugins: []TextPlugin{
			&UppercasePlugin{},
			&TrimPlugin{},
			&ReplacePlugin{Old: "WORLD", New: "GO"},
			&RepeatPlugin{Times: 2, Separator: " | "},
		},
	}

	input := "  hello, world!  "
	output := pipeline.Process(input)

	fmt.Printf("Input:  %q\n", input)
	fmt.Printf("Output: %q\n", output)

	// Show each step
	fmt.Println("\nStep-by-step:")
	current := input
	for _, plugin := range pipeline.Plugins {
		current = plugin.Transform(current)
		fmt.Printf("  %-20s → %q\n", plugin.Name(), current)
	}
}

// -----------------------------------------------------------------------
// Shape Interface & Implementations
// -----------------------------------------------------------------------

// Shape defines any geometric shape that has area and perimeter
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Circle
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func (c Circle) String() string {
	return fmt.Sprintf("Circle(r=%.1f)", c.Radius)
}

// Rectangle
type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func (r Rectangle) String() string {
	return fmt.Sprintf("Rect(%.1f×%.1f)", r.Width, r.Height)
}

// Triangle
type Triangle struct {
	Base, Height    float64
	SideA, SideB, SideC float64
}

func (t Triangle) Area() float64 {
	return 0.5 * t.Base * t.Height
}

func (t Triangle) Perimeter() float64 {
	return t.SideA + t.SideB + t.SideC
}

func (t Triangle) String() string {
	return fmt.Sprintf("Triangle(b=%.1f, h=%.1f)", t.Base, t.Height)
}

// -----------------------------------------------------------------------
// Helper Functions
// -----------------------------------------------------------------------

func printShapeInfo(s Shape) {
	fmt.Printf("  %-28s area=%-10.2f perimeter=%.2f\n", s, s.Area(), s.Perimeter())
}

func printAnything(v any) {
	fmt.Printf("  printAnything: %v (%T)\n", v, v)
}

func describeType(v any) {
	switch val := v.(type) {
	case int:
		fmt.Printf("  int: %d (doubled: %d)\n", val, val*2)
	case string:
		fmt.Printf("  string: %q (upper: %s)\n", val, strings.ToUpper(val))
	case float64:
		fmt.Printf("  float64: %f\n", val)
	case bool:
		fmt.Printf("  bool: %t\n", val)
	case Shape:
		fmt.Printf("  Shape: %v (area: %.2f)\n", val, val.Area())
	case nil:
		fmt.Println("  nil: <nothing>")
	default:
		fmt.Printf("  unknown: %v (%T)\n", val, val)
	}
}

// -----------------------------------------------------------------------
// sort.Interface example
// -----------------------------------------------------------------------

type Player struct {
	Name  string
	Score int
}

type Team []Player

func (t Team) Len() int           { return len(t) }
func (t Team) Less(i, j int) bool { return t[i].Score > t[j].Score } // descending
func (t Team) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

func (t Team) String() string {
	var parts []string
	for _, p := range t {
		parts = append(parts, fmt.Sprintf("%s(%d)", p.Name, p.Score))
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

// Simple bubble sort that uses the sort.Interface methods
func bubbleSort(data interface {
	Len() int
	Less(i, j int) bool
	Swap(i, j int)
}) {
	n := data.Len()
	for i := 0; i < n; i++ {
		for j := 0; j < n-i-1; j++ {
			if data.Less(j+1, j) {
				data.Swap(j, j+1)
			}
		}
	}
}

// -----------------------------------------------------------------------
// Interface Composition
// -----------------------------------------------------------------------

// Drawable is a separate concern from Shape
type Drawable interface {
	Draw() string
	Color() string
}

// ShapeDrawable composes both interfaces
type ShapeDrawable interface {
	Shape
	Drawable
}

// DrawableCircle satisfies both Shape (via embedding Circle) and Drawable
type DrawableCircle struct {
	Circle
	color  string
	Filled bool
}

// Need to fix: use a field name that doesn't conflict with the method
func NewDrawableCircle(radius float64, color string, filled bool) DrawableCircle {
	return DrawableCircle{
		Circle: Circle{Radius: radius},
		color:  color,
		Filled: filled,
	}
}

func (dc DrawableCircle) Draw() string {
	fill := "outline"
	if dc.Filled {
		fill = "filled"
	}
	return fmt.Sprintf("Drawing %s %s circle (r=%.1f)", dc.color, fill, dc.Radius)
}

func (dc DrawableCircle) Color() string {
	return dc.color
}

func demonstrateDrawable(sd ShapeDrawable) {
	fmt.Printf("  Shape:  area=%.2f, perimeter=%.2f\n", sd.Area(), sd.Perimeter())
	fmt.Printf("  Draw:   %s\n", sd.Draw())
	fmt.Printf("  Color:  %s\n", sd.Color())
}

// -----------------------------------------------------------------------
// Strategy Pattern — Discount Strategies
// -----------------------------------------------------------------------

// DiscountStrategy defines how to calculate a discount
type DiscountStrategy interface {
	Apply(prices []float64) float64
}

type NoDiscount struct{}

func (d *NoDiscount) Apply(prices []float64) float64 {
	total := 0.0
	for _, p := range prices {
		total += p
	}
	return total
}

type PercentageDiscount struct {
	Percent float64
}

func (d *PercentageDiscount) Apply(prices []float64) float64 {
	total := 0.0
	for _, p := range prices {
		total += p
	}
	return total * (1 - d.Percent/100)
}

type BulkDiscount struct {
	Threshold       int
	DiscountPercent float64
}

func (d *BulkDiscount) Apply(prices []float64) float64 {
	total := 0.0
	for _, p := range prices {
		total += p
	}
	if len(prices) >= d.Threshold {
		return total * (1 - d.DiscountPercent/100)
	}
	return total
}

func calculateTotal(prices []float64, strategy DiscountStrategy) float64 {
	return strategy.Apply(prices)
}

// -----------------------------------------------------------------------
// Dependency Injection — Notification Service
// -----------------------------------------------------------------------

// Notifier is the interface — the "port" in hexagonal architecture
type Notifier interface {
	Send(to, message string) error
}

// EmailNotifier is one implementation
type EmailNotifier struct {
	From string
}

func (e *EmailNotifier) Send(to, message string) error {
	fmt.Printf("  📧 Email from %s → %s: %s\n", e.From, to, message)
	return nil
}

// SMSNotifier is another implementation
type SMSNotifier struct {
	FromNumber string
}

func (s *SMSNotifier) Send(to, message string) error {
	fmt.Printf("  📱 SMS from %s → %s: %s\n", s.FromNumber, to, message)
	return nil
}

// NotificationService depends on the INTERFACE, not a concrete type
type NotificationService struct {
	Notifier Notifier // accepts any Notifier implementation
}

func (ns NotificationService) Alert(user, message string) {
	ns.Notifier.Send(user, message)
}

// -----------------------------------------------------------------------
// Plugin System — Text Processing Pipeline
// -----------------------------------------------------------------------

// TextPlugin defines a text transformation
type TextPlugin interface {
	Name() string
	Transform(input string) string
}

type Pipeline struct {
	Plugins []TextPlugin
}

func (p Pipeline) Process(input string) string {
	result := input
	for _, plugin := range p.Plugins {
		result = plugin.Transform(result)
	}
	return result
}

// Plugin implementations

type UppercasePlugin struct{}

func (p *UppercasePlugin) Name() string               { return "Uppercase" }
func (p *UppercasePlugin) Transform(s string) string { return strings.ToUpper(s) }

type TrimPlugin struct{}

func (p *TrimPlugin) Name() string               { return "Trim" }
func (p *TrimPlugin) Transform(s string) string { return strings.TrimSpace(s) }

type ReplacePlugin struct {
	Old, New string
}

func (p *ReplacePlugin) Name() string               { return fmt.Sprintf("Replace(%s→%s)", p.Old, p.New) }
func (p *ReplacePlugin) Transform(s string) string { return strings.ReplaceAll(s, p.Old, p.New) }

type RepeatPlugin struct {
	Times     int
	Separator string
}

func (p *RepeatPlugin) Name() string { return fmt.Sprintf("Repeat(%d)", p.Times) }
func (p *RepeatPlugin) Transform(s string) string {
	parts := make([]string, p.Times)
	for i := range parts {
		parts[i] = s
	}
	return strings.Join(parts, p.Separator)
}

// ============================================================================
// 🏋️ EXERCISES:
//
// 1. Add a new shape: Ellipse (semi-major axis a, semi-minor axis b).
//    Area = π * a * b, Perimeter ≈ π * (3(a+b) - √((3a+b)(a+3b)))
//    It should work seamlessly with printShapeInfo and the shapes slice.
//
// 2. Implement a Logger interface with methods:
//      Info(msg string)
//      Warn(msg string)
//      Error(msg string)
//    Create two implementations: ConsoleLogger (prints to stdout) and
//    BufferLogger (stores messages in a slice for testing).
//    Show how you can swap implementations without changing calling code.
//
// 3. Create a Storage interface:
//      Save(key string, value any) error
//      Load(key string) (any, error)
//      Delete(key string) error
//      Keys() []string
//    Implement MemoryStorage (map-based) and then show how easy it would
//    be to swap in a file-based storage later.
//
// 4. [BONUS] Build a middleware chain for an HTTP-like handler:
//      type Handler interface { Handle(request string) string }
//      type Middleware func(Handler) Handler
//    Create middleware for: Logging, Authentication (check for "auth:"),
//    and Timing. Chain them together so each wraps the next.
//
// ============================================================================

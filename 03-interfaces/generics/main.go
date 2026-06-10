// ============================================================================
// MODULE 3.3: Generics — Type Parameters in Go
// ============================================================================
//
// KEY CONCEPTS:
//   - Generics (added in Go 1.18) let you write functions and types that
//     work with ANY type, while keeping type safety
//   - Type parameters: func Name[T constraint](param T) T
//   - Constraints define what operations are allowed on type parameters
//   - "any" is the loosest constraint (= interface{})
//   - "comparable" allows == and != operations
//   - Custom constraints are just interfaces with type lists
//   - Generic types: type Stack[T any] struct { ... }
//   - Generics eliminate duplicate code without sacrificing type safety
//
// RUN: go run main.go
// ============================================================================

package main

import (
	"cmp"
	"fmt"
	"strings"
)

func main() {
	// -----------------------------------------------------------------------
	// PART 1: The Problem Generics Solve
	// -----------------------------------------------------------------------

	fmt.Println("=== The Problem Generics Solve ===")

	// WITHOUT generics, you'd write the same function for each type:
	fmt.Println("Max int:", maxInt(10, 20))
	fmt.Println("Max float:", maxFloat(3.14, 2.71))
	// maxString, maxTime, maxWhatever... 😩

	// Or use interface{}/any and lose type safety:
	// func maxAny(a, b any) any { ... } // No compile-time checks!

	// WITH generics — one function for ALL comparable types:
	fmt.Println("\nGeneric Max:")
	fmt.Println("  int:    ", Max(10, 20))
	fmt.Println("  float64:", Max(3.14, 2.71))
	fmt.Println("  string: ", Max("apple", "banana"))

	// -----------------------------------------------------------------------
	// PART 2: Generic Functions
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Generic Functions ===")

	// Type parameter syntax: func Name[T constraint](params) returns
	// T is the type parameter, constraint limits what T can be

	// Contains — works with any comparable type
	fmt.Println("Contains 3 in ints?", Contains([]int{1, 2, 3, 4}, 3))
	fmt.Println("Contains 'go' in strings?", Contains([]string{"go", "rust", "python"}, "go"))
	fmt.Println("Contains 9.9?", Contains([]float64{1.1, 2.2, 3.3}, 9.9))

	// Map — transform each element (like JavaScript's .map())
	nums := []int{1, 2, 3, 4, 5}
	doubled := Map(nums, func(n int) int { return n * 2 })
	fmt.Println("\nDoubled:", doubled)

	// Map with type change: int → string
	labels := Map(nums, func(n int) string {
		return fmt.Sprintf("item-%d", n)
	})
	fmt.Println("Labels:", labels)

	// Filter
	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	fmt.Println("Evens:", evens)

	// Reduce
	sum := Reduce(nums, 0, func(acc, n int) int { return acc + n })
	fmt.Println("Sum:", sum)

	product := Reduce(nums, 1, func(acc, n int) int { return acc * n })
	fmt.Println("Product:", product)

	// Chaining: sum of squares of evens
	result := Reduce(
		Map(
			Filter([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				func(n int) bool { return n%2 == 0 }),
			func(n int) int { return n * n }),
		0,
		func(acc, n int) int { return acc + n })
	fmt.Println("Sum of squares of evens (1-10):", result) // 4+16+36+64+100=220

	// -----------------------------------------------------------------------
	// PART 3: Type Constraints
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Type Constraints ===")

	// Built-in constraints:
	//   any        — no requirements (= interface{})
	//   comparable — supports == and != (all basic types, structs, arrays)

	// The cmp.Ordered constraint (from "cmp" package, Go 1.21+):
	//   Covers all types that support <, <=, >, >=
	//   Includes: int, float64, string, and all their variants

	// Custom constraint — see Number interface below
	fmt.Println("Sum ints:", Sum([]int{1, 2, 3, 4, 5}))
	fmt.Println("Sum float64:", Sum([]float64{1.1, 2.2, 3.3}))

	// Constraint with methods — see Stringer constraint below
	people := []Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
	}
	PrintAll(people)

	// -----------------------------------------------------------------------
	// PART 4: Generic Types — Data Structures
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Generic Types ===")

	// --- Stack ---
	fmt.Println("--- Generic Stack ---")
	intStack := NewStack[int]()
	intStack.Push(10)
	intStack.Push(20)
	intStack.Push(30)

	fmt.Printf("Stack: %v (size: %d)\n", intStack.Items(), intStack.Size())

	if val, ok := intStack.Pop(); ok {
		fmt.Println("Popped:", val)
	}
	fmt.Printf("After pop: %v\n", intStack.Items())

	// Same stack with strings
	strStack := NewStack[string]()
	strStack.Push("hello")
	strStack.Push("world")
	fmt.Printf("String stack: %v\n", strStack.Items())

	// --- Pair ---
	fmt.Println("\n--- Generic Pair ---")
	coordinate := NewPair(40.7128, -74.0060) // NYC coordinates
	fmt.Printf("NYC: (%.4f, %.4f)\n", coordinate.First, coordinate.Second)

	nameAge := NewPair("Alice", 30)
	fmt.Printf("Person: %s, age %d\n", nameAge.First, nameAge.Second)

	// --- Optional (Maybe) ---
	fmt.Println("\n--- Generic Optional ---")
	some := Some(42)
	none := None[int]()

	fmt.Printf("Some: %v, hasValue=%t\n", some, some.HasValue)
	fmt.Printf("None: %v, hasValue=%t\n", none, none.HasValue)
	fmt.Println("OrElse(some, 0):", some.OrElse(0))
	fmt.Println("OrElse(none, 0):", none.OrElse(0))

	// -----------------------------------------------------------------------
	// PART 5: Generic Type with Methods
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Generic Set ===")

	// A proper Set implementation using generics
	set := NewSet[string]()
	set.Add("go")
	set.Add("rust")
	set.Add("python")
	set.Add("go") // Duplicate — ignored

	fmt.Println("Set:", set.Items())
	fmt.Println("Contains 'go':", set.Contains("go"))
	fmt.Println("Contains 'java':", set.Contains("java"))
	fmt.Println("Size:", set.Size())

	set.Remove("rust")
	fmt.Println("After removing 'rust':", set.Items())

	// Set operations
	setA := NewSet[int]()
	for _, v := range []int{1, 2, 3, 4, 5} {
		setA.Add(v)
	}
	setB := NewSet[int]()
	for _, v := range []int{3, 4, 5, 6, 7} {
		setB.Add(v)
	}

	fmt.Println("\nSet A:", setA.Items())
	fmt.Println("Set B:", setB.Items())
	fmt.Println("Union:", Union(setA, setB).Items())
	fmt.Println("Intersection:", Intersection(setA, setB).Items())

	// -----------------------------------------------------------------------
	// PART 6: Constraints with Methods (Interface Constraints)
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Interface Constraints ===")

	// Constraints can require methods — just like regular interfaces!
	// This bridges interfaces and generics beautifully

	employees := []Employee{
		{Name: "Alice", Salary: 95000},
		{Name: "Bob", Salary: 87000},
		{Name: "Carol", Salary: 110000},
		{Name: "Dave", Salary: 78000},
	}

	products := []Product{
		{Name: "Laptop", Price: 999.99},
		{Name: "Mouse", Price: 29.99},
		{Name: "Monitor", Price: 449.99},
	}

	// SortByValue works with any type that has a Value() method
	SortByValue(employees)
	fmt.Println("Employees by salary:")
	for _, e := range employees {
		fmt.Printf("  %s: $%.0f\n", e.Name, e.Salary)
	}

	SortByValue(products)
	fmt.Println("\nProducts by price:")
	for _, p := range products {
		fmt.Printf("  %s: $%.2f\n", p.Name, p.Price)
	}

	// -----------------------------------------------------------------------
	// PART 7: Real-World Pattern — Generic Repository
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Generic Repository ===")

	// A generic CRUD repository — works with any entity type
	userRepo := NewRepository[string, User]()

	userRepo.Save("u1", User{ID: "u1", Name: "Alice", Email: "alice@example.com"})
	userRepo.Save("u2", User{ID: "u2", Name: "Bob", Email: "bob@example.com"})
	userRepo.Save("u3", User{ID: "u3", Name: "Carol", Email: "carol@example.com"})

	if user, ok := userRepo.FindByID("u2"); ok {
		fmt.Printf("Found: %+v\n", user)
	}

	fmt.Println("All users:")
	for _, u := range userRepo.FindAll() {
		fmt.Printf("  %s: %s (%s)\n", u.ID, u.Name, u.Email)
	}

	userRepo.Delete("u2")
	fmt.Printf("After delete: %d users\n", userRepo.Count())

	// Same repo for a different type!
	orderRepo := NewRepository[int, OrderRecord]()
	orderRepo.Save(1, OrderRecord{ID: 1, Product: "Laptop", Amount: 999.99})
	orderRepo.Save(2, OrderRecord{ID: 2, Product: "Mouse", Amount: 29.99})
	fmt.Printf("Orders: %d\n", orderRepo.Count())

	// -----------------------------------------------------------------------
	// PART 8: Practical Exercise — Generic Pipeline (SOLVED)
	// -----------------------------------------------------------------------

	fmt.Println("\n=== 🔗 Generic Pipeline ===")

	// A type-safe data processing pipeline using generics
	words := []string{"  Hello  ", "  WORLD  ", "go", "  Generics  ", "ROCK"}

	// Build a pipeline of transformations
	processed := Pipeline(words,
		func(s string) string { return strings.TrimSpace(s) },
		func(s string) string { return strings.ToLower(s) },
		func(s string) string { return strings.Title(s) },
	)
	fmt.Println("Processed words:", processed)

	// Numeric pipeline
	numbers := []float64{1, 4, 9, 16, 25, 36}
	results := Pipeline(numbers,
		func(n float64) float64 { return n + 1 },     // Add 1
		func(n float64) float64 { return n * 2 },     // Double
		func(n float64) float64 { return n / 10.0 },  // Scale down
	)
	fmt.Println("Processed numbers:", results)

	// GroupBy — generic grouping
	fmt.Println("\nGroupBy length:")
	wordList := []string{"go", "is", "an", "amazing", "language", "for", "building", "systems"}
	grouped := GroupBy(wordList, func(s string) int { return len(s) })
	for length, words := range grouped {
		fmt.Printf("  len=%d: %v\n", length, words)
	}
}

// -----------------------------------------------------------------------
// Non-Generic Functions (the old way)
// -----------------------------------------------------------------------

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// -----------------------------------------------------------------------
// Generic Functions
// -----------------------------------------------------------------------

// Max returns the larger of two ordered values
func Max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Contains checks if a slice contains a value
func Contains[T comparable](slice []T, target T) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

// Map transforms each element using fn
func Map[T any, U any](slice []T, fn func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

// Filter returns elements where fn returns true
func Filter[T any](slice []T, fn func(T) bool) []T {
	var result []T
	for _, v := range slice {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// Reduce combines all elements into a single value
func Reduce[T any, U any](slice []T, initial U, fn func(U, T) U) U {
	result := initial
	for _, v := range slice {
		result = fn(result, v)
	}
	return result
}

// -----------------------------------------------------------------------
// Type Constraints
// -----------------------------------------------------------------------

// Number is a custom constraint for numeric types
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~float32 | ~float64
}

// Sum adds all numbers in a slice
func Sum[T Number](nums []T) T {
	var total T
	for _, n := range nums {
		total += n
	}
	return total
}

// Stringer constraint — requires String() method
type Stringer interface {
	String() string
}

type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	return fmt.Sprintf("%s (age %d)", p.Name, p.Age)
}

// PrintAll prints any slice of Stringers
func PrintAll[T Stringer](items []T) {
	for _, item := range items {
		fmt.Printf("  → %s\n", item.String())
	}
}

// -----------------------------------------------------------------------
// Generic Types
// -----------------------------------------------------------------------

// Stack — generic LIFO data structure
type Stack[T any] struct {
	items []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (s *Stack[T]) Push(val T) {
	s.items = append(s.items, val)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	val := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return val, true
}

func (s *Stack[T]) Peek() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

func (s *Stack[T]) Size() int     { return len(s.items) }
func (s *Stack[T]) IsEmpty() bool { return len(s.items) == 0 }
func (s *Stack[T]) Items() []T    { return s.items }

// Pair — holds two values of potentially different types
type Pair[T any, U any] struct {
	First  T
	Second U
}

func NewPair[T any, U any](first T, second U) Pair[T, U] {
	return Pair[T, U]{First: first, Second: second}
}

// Optional — a value that may or may not be present
type Optional[T any] struct {
	Value    T
	HasValue bool
}

func Some[T any](val T) Optional[T] {
	return Optional[T]{Value: val, HasValue: true}
}

func None[T any]() Optional[T] {
	return Optional[T]{}
}

func (o Optional[T]) OrElse(defaultVal T) T {
	if o.HasValue {
		return o.Value
	}
	return defaultVal
}

// -----------------------------------------------------------------------
// Generic Set
// -----------------------------------------------------------------------

type Set[T comparable] struct {
	items map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{items: make(map[T]struct{})}
}

func (s *Set[T]) Add(val T)              { s.items[val] = struct{}{} }
func (s *Set[T]) Remove(val T)           { delete(s.items, val) }
func (s *Set[T]) Contains(val T) bool    { _, ok := s.items[val]; return ok }
func (s *Set[T]) Size() int              { return len(s.items) }

func (s *Set[T]) Items() []T {
	result := make([]T, 0, len(s.items))
	for k := range s.items {
		result = append(result, k)
	}
	return result
}

func Union[T comparable](a, b *Set[T]) *Set[T] {
	result := NewSet[T]()
	for _, v := range a.Items() {
		result.Add(v)
	}
	for _, v := range b.Items() {
		result.Add(v)
	}
	return result
}

func Intersection[T comparable](a, b *Set[T]) *Set[T] {
	result := NewSet[T]()
	for _, v := range a.Items() {
		if b.Contains(v) {
			result.Add(v)
		}
	}
	return result
}

// -----------------------------------------------------------------------
// Interface Constraints with Methods
// -----------------------------------------------------------------------

// Valuable — anything with a numeric value (for sorting)
type Valuable interface {
	Value() float64
}

type Employee struct {
	Name   string
	Salary float64
}

func (e Employee) Value() float64 { return e.Salary }

type Product struct {
	Name  string
	Price float64
}

func (p Product) Value() float64 { return p.Price }

// SortByValue sorts any slice of Valuables (simple insertion sort)
func SortByValue[T Valuable](items []T) {
	for i := 1; i < len(items); i++ {
		for j := i; j > 0 && items[j].Value() < items[j-1].Value(); j-- {
			items[j], items[j-1] = items[j-1], items[j]
		}
	}
}

// -----------------------------------------------------------------------
// Generic Repository
// -----------------------------------------------------------------------

type User struct {
	ID    string
	Name  string
	Email string
}

type OrderRecord struct {
	ID      int
	Product string
	Amount  float64
}

type Repository[K comparable, V any] struct {
	data map[K]V
}

func NewRepository[K comparable, V any]() *Repository[K, V] {
	return &Repository[K, V]{data: make(map[K]V)}
}

func (r *Repository[K, V]) Save(id K, entity V) {
	r.data[id] = entity
}

func (r *Repository[K, V]) FindByID(id K) (V, bool) {
	val, ok := r.data[id]
	return val, ok
}

func (r *Repository[K, V]) FindAll() []V {
	result := make([]V, 0, len(r.data))
	for _, v := range r.data {
		result = append(result, v)
	}
	return result
}

func (r *Repository[K, V]) Delete(id K) {
	delete(r.data, id)
}

func (r *Repository[K, V]) Count() int {
	return len(r.data)
}

// -----------------------------------------------------------------------
// Generic Pipeline & GroupBy
// -----------------------------------------------------------------------

// Pipeline applies a chain of transformations to each element
func Pipeline[T any](data []T, transforms ...func(T) T) []T {
	result := make([]T, len(data))
	copy(result, data)
	for _, transform := range transforms {
		for i, v := range result {
			result[i] = transform(v)
		}
	}
	return result
}

// GroupBy groups elements by a key function
func GroupBy[T any, K comparable](items []T, keyFn func(T) K) map[K][]T {
	groups := make(map[K][]T)
	for _, item := range items {
		key := keyFn(item)
		groups[key] = append(groups[key], item)
	}
	return groups
}

// ============================================================================
// 🏋️ EXERCISES:
//
// 1. Write a generic MinMax[T cmp.Ordered](slice []T) (T, T) function
//    that returns both the minimum and maximum values of a slice in a
//    single pass. Handle the empty slice case gracefully.
//
// 2. Implement a generic LinkedList[T any]:
//      type Node[T any] struct { Value T; Next *Node[T] }
//      type LinkedList[T any] struct { Head *Node[T]; size int }
//    Methods: Prepend, Append, Contains (needs comparable), ToSlice, Len
//
// 3. Write a generic Cache[K comparable, V any] with:
//      Get(key K) (V, bool)
//      Set(key K, value V, ttl time.Duration)
//      Delete(key K)
//    Entries should automatically expire after their TTL.
//    Hint: store the expiration time alongside each value.
//
// 4. [BONUS] Implement a generic MapReduce[T, U any] function:
//      MapReduce(data []T, mapFn func(T) []U, reduceFn func(U, U) U) U
//    It should: (1) apply mapFn to each element, (2) flatten the results,
//    (3) reduce them to a single value. Test with a word-count example:
//    map each sentence to words, reduce by counting.
//
// ============================================================================

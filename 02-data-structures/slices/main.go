// ============================================================================
// MODULE 2.1: Slices & Arrays
// ============================================================================
//
// KEY CONCEPTS:
//   - Arrays are FIXED-SIZE value types — their length is part of the type
//   - Slices are dynamic, flexible views into arrays — this is what you'll
//     use 99% of the time in Go
//   - A slice has three components: pointer, length, capacity
//   - append() may allocate a new backing array when capacity is exceeded
//   - Slices are reference types — multiple slices can share the same memory
//   - Use copy() to make independent copies
//
// RUN: go run main.go
// ============================================================================

package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	// -----------------------------------------------------------------------
	// PART 1: Arrays (fixed-size, rarely used directly)
	// -----------------------------------------------------------------------

	fmt.Println("=== Arrays ===")

	// Arrays have a FIXED size — the size is part of the type!
	var scores [5]int // All elements initialized to zero value (0)
	fmt.Println("Zero-value array:", scores)

	// Initialize with values
	weekdays := [5]string{"Mon", "Tue", "Wed", "Thu", "Fri"}
	fmt.Println("Weekdays:", weekdays)

	// Let the compiler count the elements with [...]
	primes := [...]int{2, 3, 5, 7, 11, 13}
	fmt.Printf("Primes: %v (length: %d)\n", primes, len(primes))

	// Access and modify elements
	scores[0] = 95
	scores[1] = 87
	scores[4] = 100
	fmt.Println("Scores:", scores)

	// CRITICAL: Arrays are VALUE TYPES — assignment copies the entire array!
	original := [3]int{1, 2, 3}
	copied := original // This is a FULL COPY, not a reference
	copied[0] = 999
	fmt.Printf("Original: %v (unchanged!)\n", original)
	fmt.Printf("Copied:   %v (only this changed)\n", copied)

	// [3]int and [5]int are DIFFERENT TYPES — you can't assign one to the other!
	// var a [3]int = [5]int{1,2,3,4,5}  // ❌ COMPILE ERROR

	// -----------------------------------------------------------------------
	// PART 2: Slice Basics
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Slice Basics ===")

	// METHOD 1: Slice literal (most common)
	fruits := []string{"apple", "banana", "cherry", "date"}
	fmt.Println("Fruits:", fruits)

	// METHOD 2: make(type, length, capacity)
	// Use make() when you know the size upfront
	buffer := make([]int, 3, 10) // length=3, capacity=10
	fmt.Printf("Buffer: %v (len=%d, cap=%d)\n", buffer, len(buffer), cap(buffer))

	// METHOD 3: make with just length (capacity = length)
	zeros := make([]float64, 5) // length=5, capacity=5
	fmt.Println("Zeros:", zeros)

	// METHOD 4: nil slice (zero value of a slice)
	var empty []int
	fmt.Printf("Nil slice: %v (len=%d, cap=%d, isNil=%t)\n",
		empty, len(empty), cap(empty), empty == nil)

	// A nil slice is perfectly safe to append to!
	empty = append(empty, 42)
	fmt.Println("After append:", empty)

	// -----------------------------------------------------------------------
	// PART 3: Slicing — creating sub-slices
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Slicing ===")

	numbers := []int{10, 20, 30, 40, 50, 60, 70, 80}
	fmt.Println("Original:", numbers)

	// Slice syntax: slice[low:high] — includes low, excludes high
	fmt.Println("numbers[2:5]  =", numbers[2:5]) // [30 40 50]
	fmt.Println("numbers[:3]   =", numbers[:3])  // [10 20 30]   (from start)
	fmt.Println("numbers[5:]   =", numbers[5:])  // [60 70 80]   (to end)
	fmt.Println("numbers[:]    =", numbers[:])   // [10 20 ... ] (everything)

	// ⚠️ WARNING: Sub-slices SHARE the same backing array!
	slice1 := numbers[1:4] // [20 30 40]
	slice1[0] = 999        // This modifies the original too!
	fmt.Println("After modifying sub-slice:")
	fmt.Println("  numbers:", numbers) // [10 999 30 40 50 60 70 80]
	fmt.Println("  slice1: ", slice1)  // [999 30 40]

	// Fix: use copy() to make an independent slice
	numbers[1] = 20 // Reset
	independent := make([]int, 3)
	copy(independent, numbers[1:4])
	independent[0] = 999
	fmt.Println("With copy():")
	fmt.Println("  numbers:    ", numbers)     // unchanged
	fmt.Println("  independent:", independent) // [999 30 40]

	// -----------------------------------------------------------------------
	// PART 4: Appending & Growing
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Appending & Growing ===")

	s := make([]int, 0, 4) // Start with capacity 4
	fmt.Printf("Start: len=%d cap=%d %v\n", len(s), cap(s), s)

	// Watch how capacity grows
	for i := 1; i <= 10; i++ {
		oldCap := cap(s)

		s = append(s, i)
		fmt.Println(s)
		if cap(s) != oldCap {
			fmt.Printf("  Grew! len=%d cap=%d (was %d) → new backing array allocated\n",
				len(s), cap(s), oldCap)
		}
	}
	fmt.Println("Final:", s)

	// Append multiple elements at once
	colors := []string{"red", "green"}
	colors = append(colors, "blue", "yellow", "purple")
	fmt.Println("Colors:", colors)

	// Append one slice to another (note the ... spread operator)
	moreColors := []string{"orange", "pink"}
	colors = append(colors, moreColors...)
	fmt.Println("All colors:", colors)

	// -----------------------------------------------------------------------
	// PART 5: Slice Internals
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Slice Internals ===")

	// A slice is a 3-word struct (called the "slice header"):
	//   1. Pointer — to the first element in the backing array
	//   2. Length  — number of elements in the slice
	//   3. Capacity — number of elements from start to end of backing array

	data := []int{10, 20, 30, 40, 50}
	sub := data[1:3] // [20, 30]

	fmt.Printf("data: ptr=%p len=%d cap=%d %v\n", data, len(data), cap(data), data)
	fmt.Printf("sub:  ptr=%p len=%d cap=%d %v\n", sub, len(sub), cap(sub), sub)
	// Notice: sub's pointer is offset, capacity is reduced from the start

	// The 3-index slice: slice[low:high:max] — controls capacity!
	limited := data[1:3:3] // len=2, cap=2 (NOT 4!)
	fmt.Printf("limited (3-index): ptr=%p len=%d cap=%d %v\n",
		limited, len(limited), cap(limited), limited)
	// Now appending to 'limited' will allocate a new array, not corrupt 'data'

	// -----------------------------------------------------------------------
	// PART 6: Common Slice Operations
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Common Slice Operations ===")

	// --- Copy ---
	src := []int{1, 2, 3, 4, 5}
	dst := make([]int, len(src))
	n := copy(dst, src) // Returns number of elements copied
	fmt.Printf("Copied %d elements: %v\n", n, dst)

	// --- Delete element at index ---
	items := []string{"a", "b", "c", "d", "e"}
	i := 2 // Delete "c"
	items = append(items[:i], items[i+1:]...)
	fmt.Println("After deleting index 2:", items) // [a b d e]

	// --- Insert element at index ---
	nums := []int{1, 2, 4, 5}
	insertAt := 2
	// Make room, then shift elements right
	nums = append(nums[:insertAt+1], nums[insertAt:]...)
	fmt.Println(nums)
	nums[insertAt] = 3
	fmt.Println("After inserting 3 at index 2:", nums) // [1 2 3 4 5]

	// --- Filter (keep only elements matching a condition) ---
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	evens := filter(values, func(n int) bool { return n%2 == 0 })
	fmt.Println("Evens:", evens)

	// --- Reverse ---
	toReverse := []int{1, 2, 3, 4, 5}
	reverse(toReverse)
	fmt.Println("Reversed:", toReverse)

	// --- Contains ---
	fmt.Println("Contains 3?", contains([]int{1, 2, 3, 4}, 3))
	fmt.Println("Contains 9?", contains([]int{1, 2, 3, 4}, 9))

	// -----------------------------------------------------------------------
	// PART 7: Multi-dimensional Slices
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Multi-dimensional Slices ===")

	// 2D slice (slice of slices) — rows can have different lengths!
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	fmt.Println("Matrix:")
	for row, cols := range matrix {
		fmt.Printf("  Row %d: %v\n", row, cols)
	}
	fmt.Printf("  matrix[1][2] = %d\n", matrix[1][2]) // 6

	// Jagged array — rows with different lengths
	triangle := make([][]int, 5)
	for i := range triangle {
		triangle[i] = make([]int, i+1)
		for j := range triangle[i] {
			triangle[i][j] = i + j
		}
	}

	fmt.Println("\nPascal-like triangle:")
	for _, row := range triangle {
		fmt.Printf("  %v\n", row)
	}

	// -----------------------------------------------------------------------
	// PART 8: Practical Exercise — Word Frequency Counter (SOLVED)
	// -----------------------------------------------------------------------

	fmt.Println("\n=== 📊 Word Frequency Counter ===")

	text := "the quick brown fox jumps over the lazy dog the fox the dog"
	words := strings.Fields(text) // Split by whitespace

	// Count frequencies using a map (previewing Module 2.2!)
	freq := make(map[string]int)
	for _, word := range words {
		freq[strings.ToLower(word)]++
	}

	// Collect unique words and sort them by frequency (descending)
	type wordCount struct {
		word  string
		count int
	}

	var counts []wordCount
	for word, count := range freq {
		counts = append(counts, wordCount{word, count})
	}
	fmt.Println(counts)
	sort.Slice(counts, func(i, j int) bool {
		return counts[i].count > counts[j].count
	})

	fmt.Printf("Text: %q\n", text)
	fmt.Printf("Total words: %d, Unique: %d\n\n", len(words), len(counts))
	fmt.Println("  Word       │ Count │ Bar")
	fmt.Println("─────────────┼───────┼──────────────")
	for _, wc := range counts {
		bar := strings.Repeat("█", wc.count)
		fmt.Printf("  %-10s │   %d   │ %s\n", wc.word, wc.count, bar)
	}

	fmt.Println()
	fmt.Printf("------testing----------------------\n")
	array := []int{1, 2, 3, 4, 5}
	r := chunks(array, 2)
	for _, row := range r {
		fmt.Printf("  %v\n", row)
	}
}

// -----------------------------------------------------------------------
// Helper Functions
// -----------------------------------------------------------------------

// filter returns a new slice containing only elements where fn returns true
func filter(s []int, fn func(int) bool) []int {
	var result []int // nil slice — append works fine
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// reverse reverses a slice in-place
func reverse(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// contains checks if a slice contains a value
func contains(s []int, target int) bool {
	for _, v := range s {
		if v == target {
			return true
		}
	}
	return false
}
func removeDuplicate(s []int) []int {
	// 1. Create a map to act as our "Set" (tracks what we've seen)
	seen := make(map[int]bool)
	// 2. Create a slice to hold the final unique numbers
	var result []int

	for _, v := range s {
		// If we haven't seen the number yet...
		if !seen[v] {
			// Mark it as seen in the map
			seen[v] = true
			// Append it to our result slice
			result = append(result, v)
		}
	}
	return result
}
func rotate(s []int, k int) []int {
	k = k % len(s)
	result := make([]int, 0, len(s))
	result = append(result, s[len(s)-k:]...)
	result = append(result, s[:len(s)-k]...)
	return result
}
func chunks(s []int, size int) [][]int {
	numsChunks := (len(s) + size - 1) / size
	chunk := make([][]int, numsChunks)

	for i := range chunk {
		start := i * size
		end := start + size
		if end > len(s) {
			end = len(s)
		}
		chunk[i] = s[start:end]
	}
	return chunk
}
type Stack struct {
	items []int
}
func NewStack() *Stack{
	return &Stack{
		items: []int {},
	}
}
func (c *Stack) Push(v int) {
	c.items = append(c.items, v)
}
func (c *Stack) Pop() (int,bool) {
	if len(c.items) == 0 {
		return 0, false
	}
	n := len(c.items)
	value := c.items[n -1]
	c.items = c.items[:n-1]
	return value, true
}
func (c *Stack) Peek() (int, bool) {
	if len(c.items) == 0{
		return 0, false
	}
	return c.items[len(c.items) -1], true
}
func (c *Stack) IsEmpty() bool{
	return len(c.items) == 0
}
func (c *Stack) Size() int{
	return len(c.items)
}
// ============================================================================
// 🏋️ EXERCISES:
//
// 1. Write a function removeDuplicates(s []int) []int that returns a new
//    slice with all duplicate values removed, preserving order.
//    Example: [1, 2, 3, 2, 1, 4] → [1, 2, 3, 4]
//
// 2. Write a function rotate(s []int, k int) []int that rotates a slice
//    k positions to the right.
//    Example: rotate([1,2,3,4,5], 2) → [4,5,1,2,3]
//    Hint: use the reverse trick — reverse all, reverse first k, reverse rest
//
// 3. Write chunk(s []int, size int) [][]int that splits a slice into
//    chunks of the given size. The last chunk may be smaller.
//    Example: chunk([1,2,3,4,5], 2) → [[1,2], [3,4], [5]]
//
// 4. [BONUS] Implement a Stack using slices with these methods:
//      Push(val int)
//      Pop() (int, bool)    — returns value and ok
//      Peek() (int, bool)
//      IsEmpty() bool
//      Size() int
//    Hint: define a type Stack struct { items []int } and add methods
//
// ============================================================================

// ============================================================================
// MODULE 2.2: Maps
// ============================================================================
//
// KEY CONCEPTS:
//   - Maps are Go's hash table / dictionary type: map[KeyType]ValueType
//   - Zero value of a map is nil — you MUST initialize before writing
//   - The comma-ok idiom: value, ok := m[key] — THE way to check existence
//   - Maps are reference types — passing a map shares the underlying data
//   - Iteration order is intentionally RANDOMIZED (don't depend on it!)
//   - Maps are NOT safe for concurrent access (use sync.Mutex or sync.Map)
//   - Keys must be comparable (==), so slices and maps can't be keys
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
	// PART 1: Map Basics
	// -----------------------------------------------------------------------

	fmt.Println("=== Map Basics ===")

	// METHOD 1: make() — creates an empty, initialized map
	population := make(map[string]int)
	population["Tokyo"] = 37400000
	population["Delhi"] = 32900000
	population["Shanghai"] = 28500000
	population["São Paulo"] = 22400000
	fmt.Println("Population:", population)

	// METHOD 2: Map literal — most common for known data
	capitals := map[string]string{
		"India":  "New Delhi",
		"Japan":  "Tokyo",
		"France": "Paris",
		"Brazil": "Brasília",
	}
	fmt.Println("Capitals:", capitals)

	// METHOD 3: Empty literal (NOT nil — safe to write to!)
	emptyMap := map[string]int{}
	emptyMap["test"] = 42
	fmt.Println("Empty literal map:", emptyMap)

	// Zero value is nil — reading returns zero value, but writing PANICS!
	var nilMap map[string]int
	fmt.Printf("Nil map: %v (isNil=%t)\n", nilMap, nilMap == nil)
	fmt.Println("Reading nil map:", nilMap["anything"]) // Returns 0, no panic
	// nilMap["key"] = 1  // ❌ RUNTIME PANIC: assignment to entry in nil map

	// -----------------------------------------------------------------------
	// PART 2: CRUD Operations
	// -----------------------------------------------------------------------

	fmt.Println("\n=== CRUD Operations ===")

	inventory := map[string]int{
		"laptop":  15,
		"mouse":   50,
		"monitor": 30,
	}

	// CREATE — add new entries
	inventory["keyboard"] = 45
	inventory["webcam"] = 20
	fmt.Println("After adding:", inventory)

	// READ — access by key (returns zero value if missing)
	fmt.Printf("Laptops in stock: %d\n", inventory["laptop"])
	fmt.Printf("Tablets in stock: %d (key doesn't exist — zero value!)\n",
		inventory["tablet"])

	// UPDATE — just assign again
	inventory["laptop"] = 12
	fmt.Printf("Updated laptops: %d\n", inventory["laptop"])

	// DELETE — built-in delete() function
	delete(inventory, "webcam")
	fmt.Println("After deleting webcam:", inventory)

	// delete() is safe on missing keys — no panic
	delete(inventory, "nonexistent") // Does nothing, no error

	// -----------------------------------------------------------------------
	// PART 3: The Comma-OK Pattern
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Comma-OK Pattern ===")

	// Problem: map returns zero value for missing keys
	// How do you tell "key exists with value 0" from "key doesn't exist"?
	scores := map[string]int{
		"Alice": 95,
		"Bob":   0, // Bob scored 0, but he took the test!
	}

	// ❌ BAD: Can't distinguish missing from zero
	fmt.Printf("Charlie's score: %d (but did Charlie take the test?)\n",
		scores["Charlie"])

	// ✅ GOOD: The comma-ok pattern
	if score, ok := scores["Bob"]; ok {
		fmt.Printf("Bob took the test and scored: %d\n", score)
	}

	if score, ok := scores["Charlie"]; ok {
		fmt.Printf("Charlie scored: %d\n", score)
	} else {
		fmt.Printf("Charlie hasn't taken the test (ok=%t, score=%d)\n",
			ok, score)
	}

	// Common pattern: use in conditionals
	if _, exists := scores["Alice"]; exists {
		fmt.Println("Alice is in the system ✓")
	}

	// -----------------------------------------------------------------------
	// PART 4: Iterating Maps
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Iterating Maps ===")

	languages := map[string]int{
		"Go":         2009,
		"Python":     1991,
		"JavaScript": 1995,
		"Rust":       2010,
		"Java":       1995,
	}

	// range returns (key, value) — ORDER IS RANDOM!
	fmt.Println("Random order (run multiple times to see it change):")
	for lang, year := range languages {
		fmt.Printf("  %s → %d\n", lang, year)
	}

	// For sorted output: extract keys, sort, then iterate
	keys := make([]string, 0, len(languages))
	for k := range languages {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Println("\nSorted by name:")
	for _, k := range keys {
		fmt.Printf("  %s → %d\n", k, languages[k])
	}

	// Count entries
	fmt.Printf("\nTotal languages: %d\n", len(languages))

	// -----------------------------------------------------------------------
	// PART 5: Maps as Sets
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Maps as Sets ===")

	// Go has no built-in set type. Use maps instead!

	// Pattern 1: map[T]bool — readable, but uses more memory
	visited := map[string]bool{
		"home":    true,
		"about":   true,
		"contact": true,
	}

	// Check membership
	page := "about"
	if visited[page] {
		fmt.Printf("Already visited: %s\n", page)
	}
	
	if !visited["pricing"] {
		fmt.Println("Haven't visited: pricing")
	}

	// Pattern 2: map[T]struct{} — zero memory for values (more efficient)
	uniqueWords := map[string]struct{}{}
	sentence := "the cat sat on the mat the cat"
	for _, word := range strings.Fields(sentence) {
		uniqueWords[word] = struct{}{} // empty struct takes 0 bytes
	}

	fmt.Printf("\nUnique words in %q:\n", sentence)
	for word := range uniqueWords {
		fmt.Printf("  • %s\n", word)
	}
	fmt.Printf("Count: %d unique words\n", len(uniqueWords))

	// Set operations
	setA := map[string]struct{}{"a": {}, "b": {}, "c": {}, "d": {}}
	setB := map[string]struct{}{"c": {}, "d": {}, "e": {}, "f": {}}

	// Intersection
	fmt.Print("\nIntersection: ")
	for k := range setA {
		if _, ok := setB[k]; ok {
			fmt.Printf("%s ", k)
		}
	}
	fmt.Println()

	// Union
	union := make(map[string]struct{})
	for k := range setA {
		union[k] = struct{}{}
	}
	for k := range setB {
		union[k] = struct{}{}
	}
	fmt.Print("Union: ")
	for k := range union {
		fmt.Printf("%s ", k)
	}
	fmt.Println()

	// -----------------------------------------------------------------------
	// PART 6: Nested Maps & Complex Values
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Nested Maps & Complex Values ===")

	// Map of slices — one key maps to multiple values
	courses := map[string][]string{
		"Alice": {"Math", "Physics", "CS"},
		"Bob":   {"History", "English"},
		"Carol": {"CS", "Math", "Art", "Music"},
	}

	for student, classes := range courses {
		fmt.Printf("  %s: %s\n", student, strings.Join(classes, ", "))
	}

	// Add a course for Bob
	courses["Bob"] = append(courses["Bob"], "Math")
	fmt.Println("\nBob's updated courses:", courses["Bob"])

	// Map of maps — nested structure
	config := map[string]map[string]string{
		"database": {
			"host":     "localhost",
			"port":     "5432",
			"name":     "myapp",
			"user":     "admin",
			"password": "secret",
		},
		"cache": {
			"host": "localhost",
			"port": "6379",
			"ttl":  "3600",
		},
	}

	fmt.Println("\nConfig:")
	for section, values := range config {
		fmt.Printf("  [%s]\n", section)
		for k, v := range values {
			if k == "password" {
				v = "****"
			}
			fmt.Printf("    %s = %s\n", k, v)
		}
	}

	// -----------------------------------------------------------------------
	// PART 7: Map Internals & Gotchas
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Map Gotchas ===")

	// Gotcha 1: Maps are reference types
	original := map[string]int{"a": 1, "b": 2}
	alias := original     // NOT a copy — both point to the same data!
	alias["a"] = 999
	fmt.Println("Original after modifying alias:", original) // a=999!

	// To make a true copy, you must copy manually
	trueCopy := make(map[string]int, len(original))
	for k, v := range original {
		trueCopy[k] = v
	}
	trueCopy["a"] = 1 // Doesn't affect original
	fmt.Println("Original after modifying trueCopy:", original)

	// Gotcha 2: Can't take address of map values
	// m := map[string]int{"x": 1}
	// ptr := &m["x"]  // ❌ COMPILE ERROR: cannot take address of m["x"]
	// Why? Because the map might reallocate its internal storage!

	// Gotcha 3: Maps are NOT safe for concurrent access
	// If you read/write from multiple goroutines → data race!
	// Solutions: sync.Mutex, sync.RWMutex, or sync.Map (Module 5)

	// Gotcha 4: Key types must be comparable (support ==)
	// ✅ string, int, float, bool, pointer, struct (if all fields comparable)
	// ❌ slices, maps, functions — cannot be map keys

	fmt.Println("\n(See comments for concurrent access and key type rules)")

	// -----------------------------------------------------------------------
	// PART 8: Practical Exercise — Contact Book (SOLVED)
	// -----------------------------------------------------------------------

	fmt.Println("\n=== 📒 Contact Book ===")

	book := NewContactBook()

	// Add contacts
	book.Add("Alice", "555-0101", "alice@example.com")
	book.Add("Bob", "555-0102", "bob@example.com")
	book.Add("Carol", "555-0103", "carol@example.com")
	book.Add("Dave", "555-0104", "dave@example.com")
	book.Add("Alice", "555-0199", "alice.new@example.com") // Update Alice

	// Search
	fmt.Println("--- Search ---")
	if contact, found := book.Search("Alice"); found {
		fmt.Printf("Found: %s — 📞 %s — 📧 %s\n",
			contact.Name, contact.Phone, contact.Email)
	}

	if _, found := book.Search("Eve"); !found {
		fmt.Println("Eve: not found")
	}

	// Delete
	book.Delete("Dave")
	fmt.Println("\n--- After deleting Dave ---")

	// List all
	book.ListAll()
}

// -----------------------------------------------------------------------
// Contact Book Implementation
// -----------------------------------------------------------------------

// Contact holds a person's contact information
type Contact struct {
	Name  string
	Phone string
	Email string
}

// ContactBook manages a collection of contacts
type ContactBook struct {
	contacts map[string]Contact
}

// NewContactBook creates and initializes a new ContactBook
func NewContactBook() *ContactBook {
	return &ContactBook{
		contacts: make(map[string]Contact),
	}
}

// Add adds or updates a contact
func (cb *ContactBook) Add(name, phone, email string) {
	cb.contacts[strings.ToLower(name)] = Contact{
		Name:  name,
		Phone: phone,
		Email: email,
	}
	fmt.Printf("  ✅ Saved: %s\n", name)
}

// Search finds a contact by name (case-insensitive)
func (cb *ContactBook) Search(name string) (Contact, bool) {
	contact, ok := cb.contacts[strings.ToLower(name)]
	return contact, ok
}

// Delete removes a contact by name
func (cb *ContactBook) Delete(name string) {
	key := strings.ToLower(name)
	if _, ok := cb.contacts[key]; ok {
		delete(cb.contacts, key)
		fmt.Printf("  🗑️  Deleted: %s\n", name)
	} else {
		fmt.Printf("  ⚠️  Not found: %s\n", name)
	}
}

// ListAll prints all contacts in sorted order
func (cb *ContactBook) ListAll() {
	if len(cb.contacts) == 0 {
		fmt.Println("  📒 Contact book is empty")
		return
	}

	// Sort by name for consistent output
	names := make([]string, 0, len(cb.contacts))
	for k := range cb.contacts {
		names = append(names, k)
	}
	sort.Strings(names)

	fmt.Println("  ┌────────────────────────────────────────────────┐")
	fmt.Println("  │              📒 All Contacts                   │")
	fmt.Println("  ├──────────┬─────────────┬───────────────────────┤")
	fmt.Println("  │ Name     │ Phone       │ Email                 │")
	fmt.Println("  ├──────────┼─────────────┼───────────────────────┤")
	for _, key := range names {
		c := cb.contacts[key]
		fmt.Printf("  │ %-8s │ %-11s │ %-21s │\n", c.Name, c.Phone, c.Email)
	}
	fmt.Println("  └──────────┴─────────────┴───────────────────────┘")
}
func wordCount(text string) map[string]int {
	counts := make(map[string]int)
	for _, word := range strings.Fields(strings.ToLower(text)) {
		counts[word]++
	}
	return counts
}

func invertMap(m map[string]int) map[int][]string {
	invertmap := make(map[int][]string)
	for k, value := range m {
		invertmap[value] = append(invertmap[value], k) 
	}
	return invertmap
}

type Person struct{
	Name string
	Age int
	City string
}
func groupBy(people []Person, fn func(Person) string) map[string][]Person {
	groups := make(map[string][]Person)
	for _, person := range people {
		key := fn(person)                              // call fn on ONE person to get the group name
		groups[key] = append(groups[key], person)      // same pattern as invertMap!
	}
	return groups
}
// entry stores a value and when it was last used
type entry struct {
	Value    int
	LastUsed int // higher number = more recently used
}

type LRUCache struct {
	cache    map[string]entry
	capacity int
	counter  int // goes up by 1 on every Get or Put
}

func NewLRU(capacity int) *LRUCache {
	return &LRUCache{
		cache:    make(map[string]entry),
		capacity: capacity,
		counter:  0,
	}
}

// Get retrieves a value and marks it as recently used
func (c *LRUCache) Get(key string) (int, bool) {
	e, ok := c.cache[key]
	if !ok {
		return 0, false
	}
	// Mark as recently used by updating the counter
	c.counter++
	e.LastUsed = c.counter
	c.cache[key] = e // maps hold copies of structs, so we must write it back!
	return e.Value, true
}

// Put adds or updates a value. If full, evicts the least recently used entry.
func (c *LRUCache) Put(key string, value int) {
	c.counter++

	// If key already exists, just update it
	if _, ok := c.cache[key]; ok {
		c.cache[key] = entry{Value: value, LastUsed: c.counter}
		return
	}

	// If cache is full, find and evict the least recently used entry
	if len(c.cache) >= c.capacity {
		lruKey := ""
		lruTime := c.counter // start with the highest possible value
		for k, e := range c.cache {
			if e.LastUsed < lruTime {
				lruTime = e.LastUsed
				lruKey = k
			}
		}
		delete(c.cache, lruKey)
		fmt.Printf("  ♻️  Evicted: %s\n", lruKey)
	}

	// Add the new entry
	c.cache[key] = entry{Value: value, LastUsed: c.counter}
}

// ============================================================================
// 🏋️ EXERCISES:
//
// 1. Write a function wordCount(text string) map[string]int that counts
//    word frequencies in text, case-insensitive.
//    Example: "The cat and the Cat" → {"the": 2, "cat": 2, "and": 1}
//
// 2. Write a function invertMap(m map[string]int) map[int][]string
//    that inverts a map. Since multiple keys may map to the same value,
//    the inverted map's values are slices.
//    Example: {"a":1, "b":2, "c":1} → {1:["a","c"], 2:["b"]}
//
// 3. Define a Person struct {Name string, Age int, City string}.
//    Write groupBy(people []Person, fn func(Person) string) map[string][]Person
//    that groups people by whatever criteria fn returns.
//    Test with grouping by City and by age range ("young"/"adult"/"senior").
//
// 4. [BONUS] Implement a simple LRU (Least Recently Used) cache:
//      type LRUCache struct { ... }
//      func NewLRU(capacity int) *LRUCache
//      func (c *LRUCache) Get(key string) (int, bool)
//      func (c *LRUCache) Put(key string, value int)
//    When capacity is exceeded, evict the least recently used entry.
//    Hint: track access order with a counter or timestamp in each entry.
//
// ============================================================================

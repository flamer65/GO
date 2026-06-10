// ============================================================================
// MODULE 3.2: Error Handling — The Go Way
// ============================================================================
//
// KEY CONCEPTS:
//   - Errors are VALUES, not exceptions — they implement the error interface
//   - error interface: just Error() string — the simplest interface in Go
//   - Always check errors! if err != nil { ... } is the Go rhythm
//   - fmt.Errorf with %w WRAPS errors — preserving the chain
//   - errors.Is checks if an error matches (even through wrapping)
//   - errors.As extracts a specific error type from the chain
//   - Custom error types carry structured data (status codes, context)
//   - Sentinel errors (var ErrNotFound = ...) are named, reusable errors
//   - panic/recover is NOT for normal error handling — use errors!
//
// RUN: go run main.go
// ============================================================================

package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

func main() {
	// -----------------------------------------------------------------------
	// PART 1: The error Interface
	// -----------------------------------------------------------------------

	fmt.Println("=== The error Interface ===")

	// The error interface is Go's simplest:
	//   type error interface {
	//       Error() string
	//   }
	// Any type with an Error() string method IS an error

	// Creating basic errors
	err1 := errors.New("something went wrong")
	err2 := fmt.Errorf("failed to process item %d: %s", 42, "timeout")

	fmt.Println("errors.New:", err1)
	fmt.Println("fmt.Errorf:", err2)

	// Errors are just values — you can store, compare, return them
	var savedError error = err1
	fmt.Printf("Type: %T, Value: %v\n", savedError, savedError)

	// nil means "no error" — this is THE success signal
	var noError error = nil
	fmt.Printf("No error: %v (isNil: %t)\n", noError, noError == nil)

	// -----------------------------------------------------------------------
	// PART 2: The if err != nil Pattern
	// -----------------------------------------------------------------------

	fmt.Println("\n=== if err != nil ===")

	// This is the heartbeat of Go. You'll write it thousands of times.
	// It's explicit, it's simple, and it works.

	// Example: a function that can fail
	result, err := divide(10, 3)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("10 / 3 = %.4f\n", result)
	}

	// Error case
	result, err = divide(10, 0)
	if err != nil {
		fmt.Println("Error:", err) // Error: division by zero
	}

	// Real-world pattern: early return on error
	fmt.Println("\nProcessing users:")
	users := []string{"alice", "bob", "", "dave"}
	for _, name := range users {
		if err := processUser(name); err != nil {
			fmt.Printf("  ❌ %s: %v\n", name, err)
			continue // or return, depending on context
		}
		fmt.Printf("  ✅ %s: processed\n", name)
	}

	// -----------------------------------------------------------------------
	// PART 3: Sentinel Errors
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Sentinel Errors ===")

	// Sentinel errors are package-level variables — named, reusable errors
	// Convention: var ErrXxx = errors.New("description")
	// Examples from stdlib: io.EOF, sql.ErrNoRows, os.ErrNotExist

	// See our sentinel errors defined below
	db := NewMockDB()
	db.Set("user:1", "Alice")
	db.Set("user:2", "Bob")

	// Using sentinel errors for control flow
	for _, key := range []string{"user:1", "user:3", ""} {
		val, err := db.Get(key)
		switch {
		case err == nil:
			fmt.Printf("  Found: %s = %s\n", key, val)
		case errors.Is(err, ErrNotFound):
			fmt.Printf("  Not found: %s\n", key)
		case errors.Is(err, ErrEmptyKey):
			fmt.Printf("  Bad request: empty key\n")
		default:
			fmt.Printf("  Unexpected: %v\n", err)
		}
	}

	// -----------------------------------------------------------------------
	// PART 4: Custom Error Types
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Custom Error Types ===")

	// Custom error types carry structured data beyond just a message
	// They implement the error interface by having Error() string

	// See ValidationError and HTTPError defined below

	// Validation errors
	if err := validateAge(-5); err != nil {
		fmt.Println("Validation:", err)

		// Type assert to get structured data
		var ve *ValidationError
		if errors.As(err, &ve) {
			fmt.Printf("  Field: %s, Value: %v, Rule: %s\n",
				ve.Field, ve.Value, ve.Rule)
		}
	}

	// HTTP-like errors
	if err := fetchResource("/api/secret"); err != nil {
		fmt.Println("\nHTTP Error:", err)

		var he *HTTPError
		if errors.As(err, &he) {
			fmt.Printf("  Status: %d, Retryable: %t\n",
				he.StatusCode, he.StatusCode >= 500)
		}
	}

	// -----------------------------------------------------------------------
	// PART 5: Error Wrapping with %w
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Error Wrapping ===")

	// fmt.Errorf with %w WRAPS an error — creating a chain
	// This preserves context while maintaining the original error

	err = loadUserProfile("user:999")
	if err != nil {
		fmt.Println("Top-level error:", err)
		// The error message shows the full chain:
		// "load user profile: get user: key not found: user:999"

		// errors.Is walks the chain to find a match
		fmt.Println("Is ErrNotFound?", errors.Is(err, ErrNotFound))

		// Unwrap manually to see the chain
		fmt.Println("\nError chain:")
		printErrorChain(err)
	}

	// -----------------------------------------------------------------------
	// PART 6: errors.Is and errors.As
	// -----------------------------------------------------------------------

	fmt.Println("\n=== errors.Is and errors.As ===")

	// errors.Is — checks if any error in the chain matches a target
	// (like == but works through wrapping)

	wrapped := fmt.Errorf("layer 3: %w",
		fmt.Errorf("layer 2: %w",
			fmt.Errorf("layer 1: %w", ErrNotFound)))

	fmt.Println("Error:", wrapped)
	fmt.Println("Is ErrNotFound?", errors.Is(wrapped, ErrNotFound)) // true!
	fmt.Println("Is ErrEmptyKey?", errors.Is(wrapped, ErrEmptyKey)) // false

	// errors.As — extracts a specific error TYPE from the chain
	complexErr := fmt.Errorf("request failed: %w",
		&HTTPError{StatusCode: 503, Message: "Service Unavailable",
			Endpoint: "/api/data"})

	var httpErr *HTTPError
	if errors.As(complexErr, &httpErr) {
		fmt.Printf("\nExtracted HTTPError: %d %s at %s\n",
			httpErr.StatusCode, httpErr.Message, httpErr.Endpoint)
	}

	// -----------------------------------------------------------------------
	// PART 7: Error Handling Patterns
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Error Handling Patterns ===")

	// Pattern 1: Collect multiple errors
	fmt.Println("--- Multi-Error ---")
	errs := validateForm(map[string]string{
		"name":  "",
		"email": "not-an-email",
		"age":   "-5",
	})
	if errs != nil {
		fmt.Println("Form errors:")
		for _, e := range errs {
			fmt.Printf("  • %v\n", e)
		}
	}

	// Pattern 2: Error with retry
	fmt.Println("\n--- Retry Pattern ---")
	attempts := 0
	err = retry(3, 100*time.Millisecond, func() error {
		attempts++
		if attempts < 3 {
			return fmt.Errorf("connection refused (attempt %d)", attempts)
		}
		return nil // Success on 3rd try
	})
	if err != nil {
		fmt.Println("Final error:", err)
	} else {
		fmt.Printf("Succeeded after %d attempts\n", attempts)
	}

	// Pattern 3: Must — panic on error (only for initialization!)
	fmt.Println("\n--- Must Pattern ---")
	// Use sparingly — only in init() or main setup
	value := must(divide(100, 4))
	fmt.Printf("must(divide(100, 4)) = %.2f\n", value)

	// This would panic:
	// must(divide(100, 0))  // 💥 panic: division by zero

	// -----------------------------------------------------------------------
	// PART 8: Practical Exercise — Result Type & Error Pipeline (SOLVED)
	// -----------------------------------------------------------------------

	fmt.Println("\n=== 🔧 Order Processing Pipeline ===")

	orders := []Order{
		{ID: "ORD-001", Customer: "Alice", Items: []Item{
			{Name: "Laptop", Price: 999.99, Qty: 1},
			{Name: "Mouse", Price: 29.99, Qty: 2},
		}},
		{ID: "ORD-002", Customer: "", Items: []Item{ // Empty customer
			{Name: "Keyboard", Price: 79.99, Qty: 1},
		}},
		{ID: "ORD-003", Customer: "Carol", Items: []Item{}}, // No items
		{ID: "ORD-004", Customer: "Dave", Items: []Item{
			{Name: "Monitor", Price: 499.99, Qty: 1},
			{Name: "Cable", Price: -5.00, Qty: 1}, // Negative price
		}},
		{ID: "ORD-005", Customer: "Eve", Items: []Item{
			{Name: "Headphones", Price: 199.99, Qty: 3},
		}},
	}

	fmt.Println()
	var processed, failed int
	for _, order := range orders {
		if err := processOrder(order); err != nil {
			fmt.Printf("  ❌ %s: %v\n", order.ID, err)
			failed++
		} else {
			processed++
		}
	}
	fmt.Printf("\n  Summary: %d processed, %d failed\n", processed, failed)
}

// -----------------------------------------------------------------------
// Basic Functions
// -----------------------------------------------------------------------

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

func processUser(name string) error {
	if name == "" {
		return errors.New("empty username")
	}
	return nil
}

// -----------------------------------------------------------------------
// Sentinel Errors
// -----------------------------------------------------------------------

var (
	ErrNotFound = errors.New("not found")
	ErrEmptyKey = errors.New("empty key")
	ErrExpired  = errors.New("expired")
)

// MockDB simulates a simple key-value store
type MockDB struct {
	data map[string]string
}

func NewMockDB() *MockDB {
	return &MockDB{data: make(map[string]string)}
}

func (db *MockDB) Set(key, value string) {
	db.data[key] = value
}

func (db *MockDB) Get(key string) (string, error) {
	if key == "" {
		return "", ErrEmptyKey
	}
	val, ok := db.data[key]
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrNotFound, key)
	}
	return val, nil
}

// -----------------------------------------------------------------------
// Custom Error Types
// -----------------------------------------------------------------------

// ValidationError carries field-level validation details
type ValidationError struct {
	Field   string
	Value   any
	Rule    string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %s %s (got: %v)", e.Field, e.Rule, e.Value)
}

func validateAge(age int) error {
	if age < 0 {
		return &ValidationError{
			Field:   "age",
			Value:   age,
			Rule:    "must be non-negative",
			Message: "Age cannot be negative",
		}
	}
	if age > 150 {
		return &ValidationError{
			Field:   "age",
			Value:   age,
			Rule:    "must be ≤ 150",
			Message: "Unrealistic age",
		}
	}
	return nil
}

// HTTPError represents an HTTP error response
type HTTPError struct {
	StatusCode int
	Message    string
	Endpoint   string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s (endpoint: %s)", e.StatusCode, e.Message, e.Endpoint)
}

func fetchResource(endpoint string) error {
	if strings.HasPrefix(endpoint, "/api/secret") {
		return &HTTPError{
			StatusCode: 403,
			Message:    "Forbidden",
			Endpoint:   endpoint,
		}
	}
	return nil
}

// -----------------------------------------------------------------------
// Error Wrapping
// -----------------------------------------------------------------------

func loadUserProfile(userID string) error {
	userData, err := getUserFromDB(userID)
	if err != nil {
		return fmt.Errorf("load user profile: %w", err)
	}
	_ = userData
	return nil
}

func getUserFromDB(userID string) (string, error) {
	db := NewMockDB()
	val, err := db.Get(userID)
	if err != nil {
		return "", fmt.Errorf("get user: %w", err)
	}
	return val, nil
}

func printErrorChain(err error) {
	for i := 0; err != nil; i++ {
		fmt.Printf("  [%d] %v\n", i, err)
		err = errors.Unwrap(err)
	}
}

// -----------------------------------------------------------------------
// Error Handling Patterns
// -----------------------------------------------------------------------

// Multi-error collection
func validateForm(fields map[string]string) []error {
	var errs []error

	if name := fields["name"]; name == "" {
		errs = append(errs, &ValidationError{
			Field: "name", Value: name, Rule: "required",
		})
	}

	if email := fields["email"]; !strings.Contains(email, "@") {
		errs = append(errs, &ValidationError{
			Field: "email", Value: email, Rule: "must contain @",
		})
	}

	if age := fields["age"]; strings.HasPrefix(age, "-") {
		errs = append(errs, &ValidationError{
			Field: "age", Value: age, Rule: "must be non-negative",
		})
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

// Retry with backoff
func retry(maxAttempts int, delay time.Duration, fn func() error) error {
	var lastErr error
	for i := 1; i <= maxAttempts; i++ {
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		fmt.Printf("  Attempt %d failed: %v\n", i, lastErr)
		if i < maxAttempts {
			time.Sleep(delay)
			delay *= 2 // Exponential backoff
		}
	}
	return fmt.Errorf("failed after %d attempts: %w", maxAttempts, lastErr)
}

// Must — panics on error (use only for initialization)
func must(val float64, err error) float64 {
	if err != nil {
		panic(err)
	}
	return val
}

// -----------------------------------------------------------------------
// Order Processing Pipeline
// -----------------------------------------------------------------------

type Item struct {
	Name  string
	Price float64
	Qty   int
}

type Order struct {
	ID       string
	Customer string
	Items    []Item
}

func processOrder(order Order) error {
	// Step 1: Validate customer
	if order.Customer == "" {
		return fmt.Errorf("validate: customer name is required")
	}

	// Step 2: Validate items
	if len(order.Items) == 0 {
		return fmt.Errorf("validate: order must have at least one item")
	}

	// Step 3: Validate each item
	total := 0.0
	for _, item := range order.Items {
		if item.Price < 0 {
			return fmt.Errorf("validate: item %q has negative price ($%.2f)",
				item.Name, item.Price)
		}
		if item.Qty <= 0 {
			return fmt.Errorf("validate: item %q has invalid quantity (%d)",
				item.Name, item.Qty)
		}
		total += item.Price * float64(item.Qty)
	}

	// Step 4: Process payment (simulated)
	if err := chargePayment(order.Customer, total); err != nil {
		return fmt.Errorf("payment: %w", err)
	}

	fmt.Printf("  ✅ %s: %s — %d items — $%.2f\n",
		order.ID, order.Customer, len(order.Items), total)
	return nil
}

func chargePayment(customer string, amount float64) error {
	if amount > 10000 {
		return fmt.Errorf("amount $%.2f exceeds limit", amount)
	}
	// Simulate successful payment
	_ = customer
	_ = math.Round(amount*100) / 100
	return nil
}

// ============================================================================
// 🏋️ EXERCISES:
//
// 1. Create a FileError custom type with fields: Op (operation like "read",
//    "write"), Path (file path), and Err (underlying error).
//    Implement Error() and Unwrap(). Show how errors.Is and errors.As
//    work with it through multiple layers of wrapping.
//
// 2. Implement a Result[T] generic type (or use struct with value + error)
//    that wraps either a success value or an error. Add methods:
//      IsOk() bool
//      IsErr() bool
//      Unwrap() T          — panics if error
//      UnwrapOr(default T) — returns default if error
//      Map(fn func(T) T) Result[T]
//
// 3. Build an error aggregator that collects errors from multiple
//    concurrent operations (simulate with goroutines if you've done
//    Module 5, or just sequential calls):
//      type MultiError struct { Errors []error }
//    Implement Error() that formats all errors, and add HasErrors() bool.
//
// 4. [BONUS] Implement a circuit breaker pattern:
//      type CircuitBreaker struct { ... }
//    It wraps a function and tracks failures. After N consecutive failures,
//    it "opens" the circuit and returns errors immediately without calling
//    the function. After a timeout, it allows one "probe" call through.
//    States: Closed (normal) → Open (failing) → HalfOpen (testing)
//
// ============================================================================

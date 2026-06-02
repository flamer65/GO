// ============================================================================
// MODULE 2.4: Structs & Methods
// ============================================================================
//
// KEY CONCEPTS:
//   - Structs are Go's way to define custom data types (like classes, minus
//     inheritance)
//   - Methods are functions with a "receiver" argument attached to a type
//   - Value receivers get a COPY; pointer receivers can MODIFY the original
//   - Embedding provides composition — Go's alternative to inheritance
//   - Struct tags control JSON encoding and other serialization
//   - Constructor functions (NewXxx pattern) initialize structs properly
//
// RUN: go run main.go
// ============================================================================

package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"math"
)

func main() {
	// -----------------------------------------------------------------------
	// PART 1: Struct Basics
	// -----------------------------------------------------------------------

	fmt.Println("=== Struct Basics ===")

	// Defining a struct type (usually done outside functions)
	// See the Person struct defined below main()

	// Creating struct values
	// METHOD 1: Named fields (preferred — order doesn't matter)
	alice := Person{
		FirstName: "Alice",
		LastName:  "Johnson",
		Age:       30,
		Email:     "alice@example.com",
	}
	fmt.Println("Alice:", alice)
	fmt.Printf("Full name: %s %s, Age: %d\n",
		alice.FirstName, alice.LastName, alice.Age)

	// METHOD 2: Positional (fragile — breaks if you add fields!)
	bob := Person{"Bob", "Smith", 25, "bob@example.com"}
	fmt.Println("Bob:", bob)

	// METHOD 3: Zero value — all fields get their zero values
	var unknown Person
	fmt.Printf("Zero person: %+v\n", unknown) // %+v shows field names

	// METHOD 4: Partial initialization — unset fields get zero values
	carol := Person{FirstName: "Carol", Age: 28}
	fmt.Printf("Partial: %+v\n", carol)

	// Field access and modification
	alice.Age = 31
	fmt.Printf("Alice is now %d\n", alice.Age)

	// -----------------------------------------------------------------------
	// PART 2: Struct Initialization Patterns
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Struct Initialization ===")

	// Using new() — returns a pointer to a zero-valued struct
	dave := new(Person)
	dave.FirstName = "Dave"
	dave.LastName = "Wilson"
	dave.Age = 35
	fmt.Printf("new(): %+v (type: %T)\n", *dave, dave)

	// Address-of operator — most common way to get a pointer
	eve := &Person{
		FirstName: "Eve",
		LastName:  "Davis",
		Age:       22,
		Email:     "eve@example.com",
	} 
	fmt.Printf("&literal: %+v (type: %T)\n", *eve, eve)

	// Anonymous structs — useful for one-off data, tests, and JSON
	config := struct {
		Host string
		Port int
		TLS  bool
	}{
		Host: "localhost",
		Port: 8080,
		TLS:  false,
	}
	fmt.Printf("Anonymous: %+v\n", config)

	// -----------------------------------------------------------------------
	// PART 3: Pointers to Structs
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Pointers to Structs ===")

	// Structs are VALUE TYPES — assignment copies the entire struct
	original := Person{FirstName: "Original", Age: 50}
	copied := original
	copied.FirstName = "Copied"
	fmt.Printf("Original: %s (unchanged!)\n", original.FirstName)
	fmt.Printf("Copied:   %s\n", copied.FirstName)

	// Use pointers to share and modify
	ptr := &original
	ptr.Age = 51 // Go auto-dereferences — same as (*ptr).Age = 51
	fmt.Printf("Via pointer: %s is now %d\n", original.FirstName, original.Age)

	// Why pointers matter: passing structs to functions
	p := Person{FirstName: "Test", Age: 20}
	cantModify(p)
	fmt.Printf("After cantModify: Age=%d (unchanged)\n", p.Age)
	canModify(&p)
	fmt.Printf("After canModify:  Age=%d (modified!)\n", p.Age)

	// -----------------------------------------------------------------------
	// PART 4: Methods
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Methods ===")

	// Methods are functions with a receiver
	// See method definitions below main()

	rect := Rectangle{Width: 10, Height: 5}

	// Value receiver methods — work on a copy
	fmt.Printf("Rectangle: %s\n", rect.String())
	fmt.Printf("Area: %.1f\n", rect.Area())
	fmt.Printf("Perimeter: %.1f\n", rect.Perimeter())
	fmt.Printf("Is Square? %t\n", rect.IsSquare())

	// Pointer receiver methods — modify the original
	fmt.Printf("Before scale: %s\n", rect.String())
	rect.Scale(2.0) // Go auto-takes address: (&rect).Scale(2.0)
	fmt.Printf("After Scale(2): %s\n", rect.String())

	// Rule of thumb for receivers:
	// Use POINTER receiver when:
	//   1. The method modifies the receiver
	//   2. The struct is large (avoid copying)
	//   3. Consistency — if one method needs pointer, use pointer for all
	// Use VALUE receiver when:
	//   1. The struct is small and read-only
	//   2. You want the method to work on a copy (safety)

	// -----------------------------------------------------------------------
	// PART 5: Struct Embedding (Composition)
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Struct Embedding ===")

	// Go doesn't have inheritance. Instead, use EMBEDDING for composition.
	// An embedded struct's fields and methods are "promoted" — accessible directly.

	// See Employee struct defined below
	emp := Employee{
		Person: Person{
			FirstName: "Alice",
			LastName:  "Johnson",
			Age:       30,
			Email:     "alice@company.com",
		},
		Company:  "Gopher Corp",
		Position: "Senior Engineer",
		Salary:   120000,
	}

	// Promoted fields — access Person fields directly!
	fmt.Printf("Name: %s %s\n", emp.FirstName, emp.LastName) // No need for emp.Person.FirstName
	fmt.Printf("Company: %s, Position: %s\n", emp.Company, emp.Position)
	fmt.Printf("Salary: $%d\n", emp.Salary)

	// Promoted methods work too
	fmt.Printf("Full name: %s\n", emp.FullName()) // Person's method

	// Employee can have its own methods that override or extend
	fmt.Printf("Badge: %s\n", emp.Badge())

	// You can still access the embedded struct explicitly
	fmt.Printf("Embedded Person: %+v\n", emp.Person)

	// Multiple embedding
	manager := Manager{
		Employee: emp,
		Reports:  []string{"Bob", "Carol", "Dave"},
	}
	// Fields promoted through multiple levels
	fmt.Printf("\nManager: %s %s\n", manager.FirstName, manager.LastName)
	fmt.Printf("Reports: %v\n", manager.Reports)
	fmt.Printf("Team size: %d\n", manager.TeamSize())

	// -----------------------------------------------------------------------
	// PART 6: Struct Tags (JSON)
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Struct Tags (JSON) ===")

	// Struct tags are string annotations that control encoding/decoding
	// The most common use: JSON field mapping

	// See Product struct defined below
	product := Product{
		ID:          "prod-001",
		Name:        "Mechanical Keyboard",
		Price:       149.99,
		InStock:     true,
		Tags:        []string{"electronics", "peripherals", "gaming"},
		Description: "", // Will be omitted from JSON (omitempty)
	}

	// Marshal — Go struct → JSON
	jsonBytes, err := json.MarshalIndent(product, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Go → JSON:")
	fmt.Println(string(jsonBytes))

	// Unmarshal — JSON → Go struct
	jsonInput := `{
		"id": "prod-002",
		"name": "Wireless Mouse",
		"price": 59.99,
		"in_stock": false,
		"tags": ["electronics", "peripherals"],
		"description": "Ergonomic wireless mouse"
	}`

	var decoded Product
	err = json.Unmarshal([]byte(jsonInput), &decoded)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("\nJSON → Go: %+v\n", decoded)

	// -----------------------------------------------------------------------
	// PART 7: Constructor Patterns
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Constructor Patterns ===")

	// Pattern 1: Simple constructor (NewXxx)
	acc := NewBankAccount("Alice", 1000.00)
	fmt.Printf("Account: %s, Balance: $%.2f\n", acc.Owner, acc.Balance)

	// Pattern 2: Constructor with validation
	_, err = NewValidatedAccount("", -100)
	if err != nil {
		fmt.Println("Validation error:", err)
	}

	validAcc, _ := NewValidatedAccount("Bob", 500)
	fmt.Printf("Valid account: %+v\n", *validAcc)

	// Pattern 3: Functional Options (advanced but powerful!)
	server := NewServer(
		WithPort(9090),
		WithTimeout(30*time.Second),
		WithMaxConnections(100),
		WithTLS(true),
	)
	fmt.Printf("Server: %+v\n", *server)

	// Default server (no options)
	defaultServer := NewServer()
	fmt.Printf("Default: %+v\n", *defaultServer)

	// -----------------------------------------------------------------------
	// PART 8: Practical Exercise — Bank Account System (SOLVED)
	// -----------------------------------------------------------------------

	fmt.Println("\n=== 🏦 Bank Account System ===")

	// Create accounts
	checking := NewBankAccount("Naman", 5000.00)
	savings := NewBankAccount("Naman (Savings)", 10000.00)

	// Transactions
	checking.Deposit(1500.00)
	checking.Deposit(750.00)
	checking.Withdraw(200.00)
	checking.Withdraw(50.00)

	savings.Deposit(2000.00)

	// Transfer between accounts
	fmt.Println()
	checking.Transfer(savings, 1000.00)

	// Try invalid operations
	fmt.Println()
	checking.Withdraw(999999.00)
	checking.Deposit(-100)

	// Print statements
	fmt.Println()
	checking.Statement()
	fmt.Println()
	savings.Statement()
}

// -----------------------------------------------------------------------
// Struct Definitions
// -----------------------------------------------------------------------

// Person represents a basic person
type Person struct {
	FirstName string
	LastName  string
	Age       int
	Email     string
}

// FullName returns the person's full name
func (p Person) FullName() string {
	return p.FirstName + " " + p.LastName
}

// cantModify receives a copy — changes don't affect the original
func cantModify(p Person) {
	p.Age = 999 // Modifies the copy, not the original
}

// canModify receives a pointer — changes affect the original
func canModify(p *Person) {
	p.Age = 21 // Modifies the original
}

// -----------------------------------------------------------------------
// Rectangle with Methods
// -----------------------------------------------------------------------

// Rectangle represents a rectangle with width and height
type Rectangle struct {
	Width, Height float64
}

// Value receiver — doesn't modify the struct
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func (r Rectangle) IsSquare() bool {
	return r.Width == r.Height
}

func (r Rectangle) String() string {
	return fmt.Sprintf("Rectangle(%.1f × %.1f)", r.Width, r.Height)
}

// Pointer receiver — modifies the struct
func (r *Rectangle) Scale(factor float64) {
	r.Width *= factor
	r.Height *= factor
}

// -----------------------------------------------------------------------
// Embedding (Composition)
// -----------------------------------------------------------------------

// Employee embeds Person — gains all Person fields and methods
type Employee struct {
	Person            // Embedded (anonymous) field
	Company  string
	Position string
	Salary   int
}

// Badge creates an employee badge string
func (e Employee) Badge() string {
	return fmt.Sprintf("[%s] %s — %s", e.Company, e.FullName(), e.Position)
}

// Manager embeds Employee
type Manager struct {
	Employee
	Reports []string
}

// TeamSize returns the number of direct reports
func (m Manager) TeamSize() int {
	return len(m.Reports)
}

// -----------------------------------------------------------------------
// Struct Tags (JSON)
// -----------------------------------------------------------------------

// Product demonstrates struct tags for JSON encoding
type Product struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Price       float64  `json:"price"`
	InStock     bool     `json:"in_stock"`
	Tags        []string `json:"tags,omitempty"`
	Description string   `json:"description,omitempty"` // Omit if empty
	Internal    string   `json:"-"`                     // Always omit from JSON
}

// -----------------------------------------------------------------------
// Constructor Patterns
// -----------------------------------------------------------------------

// Simple constructor
func NewBankAccount(owner string, initialBalance float64) *BankAccount {
	return &BankAccount{
		Owner:   owner,
		Balance: initialBalance,
		History: []Transaction{
			{
				Type:      "OPEN",
				Amount:    initialBalance,
				Balance:   initialBalance,
				Timestamp: time.Now(),
				Note:      "Account opened",
			},
		},
	}
}

// Constructor with validation
type ValidatedAccount struct {
	Owner   string
	Balance float64
}

func NewValidatedAccount(owner string, balance float64) (*ValidatedAccount, error) {
	if owner == "" {
		return nil, fmt.Errorf("owner name cannot be empty")
	}
	if balance < 0 {
		return nil, fmt.Errorf("initial balance cannot be negative: %.2f", balance)
	}
	return &ValidatedAccount{Owner: owner, Balance: balance}, nil
}

// Functional Options Pattern
type Server struct {
	Port           int
	Timeout        time.Duration
	MaxConnections int
	TLS            bool
}

// Option is a function that configures a Server
type Option func(*Server)

func WithPort(port int) Option {
	return func(s *Server) { s.Port = port }
}

func WithTimeout(d time.Duration) Option {
	return func(s *Server) { s.Timeout = d }
}

func WithMaxConnections(max int) Option {
	return func(s *Server) { s.MaxConnections = max }
}

func WithTLS(tls bool) Option {
	return func(s *Server) { s.TLS = tls }
}

func NewServer(opts ...Option) *Server {
	// Start with sensible defaults
	s := &Server{
		Port:           8080,
		Timeout:        10 * time.Second,
		MaxConnections: 50,
		TLS:            false,
	}
	// Apply each option
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// -----------------------------------------------------------------------
// Bank Account System
// -----------------------------------------------------------------------

// Transaction records a single account operation
type Transaction struct {
	Type      string // DEPOSIT, WITHDRAW, TRANSFER_IN, TRANSFER_OUT
	Amount    float64
	Balance   float64
	Timestamp time.Time
	Note      string
}

// BankAccount represents a bank account with transaction history
type BankAccount struct {
	Owner   string
	Balance float64
	History []Transaction
}

// Deposit adds money to the account
func (a *BankAccount) Deposit(amount float64) error {
	if amount <= 0 {
		fmt.Printf("  ❌ Invalid deposit amount: $%.2f\n", amount)
		return fmt.Errorf("deposit amount must be positive")
	}
	a.Balance += amount
	a.History = append(a.History, Transaction{
		Type:      "DEPOSIT",
		Amount:    amount,
		Balance:   a.Balance,
		Timestamp: time.Now(),
		Note:      "Cash deposit",
	})
	fmt.Printf("  💰 Deposited $%.2f → Balance: $%.2f\n", amount, a.Balance)
	return nil
}

// Withdraw removes money from the account
func (a *BankAccount) Withdraw(amount float64) error {
	if amount <= 0 {
		fmt.Printf("  ❌ Invalid withdrawal amount: $%.2f\n", amount)
		return fmt.Errorf("withdrawal amount must be positive")
	}
	if amount > a.Balance {
		fmt.Printf("  ❌ Insufficient funds: need $%.2f, have $%.2f\n",
			amount, a.Balance)
		return fmt.Errorf("insufficient funds")
	}
	a.Balance -= amount
	a.History = append(a.History, Transaction{
		Type:      "WITHDRAW",
		Amount:    amount,
		Balance:   a.Balance,
		Timestamp: time.Now(),
		Note:      "Cash withdrawal",
	})
	fmt.Printf("  💸 Withdrew $%.2f → Balance: $%.2f\n", amount, a.Balance)
	return nil
}

// Transfer moves money from this account to another
func (a *BankAccount) Transfer(to *BankAccount, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("transfer amount must be positive")
	}
	if amount > a.Balance {
		fmt.Printf("  ❌ Cannot transfer $%.2f: insufficient funds ($%.2f)\n",
			amount, a.Balance)
		return fmt.Errorf("insufficient funds")
	}

	a.Balance -= amount
	a.History = append(a.History, Transaction{
		Type:      "TRANSFER_OUT",
		Amount:    amount,
		Balance:   a.Balance,
		Timestamp: time.Now(),
		Note:      fmt.Sprintf("Transfer to %s", to.Owner),
	})

	to.Balance += amount
	to.History = append(to.History, Transaction{
		Type:      "TRANSFER_IN",
		Amount:    amount,
		Balance:   to.Balance,
		Timestamp: time.Now(),
		Note:      fmt.Sprintf("Transfer from %s", a.Owner),
	})

	fmt.Printf("  🔄 Transferred $%.2f: %s → %s\n", amount, a.Owner, to.Owner)
	return nil
}

// Statement prints a formatted account statement
func (a *BankAccount) Statement() {
	width := 83
	line := strings.Repeat("─", width)

	fmt.Println("  ┌" + line + "┐")
	title := fmt.Sprintf("  │  🏦 Account Statement: %-37s│", a.Owner)
	fmt.Println(title)
	fmt.Println("  ├" + line + "┤")
	fmt.Println("  │  Type           │   Amount   │   Balance  │ Note               │ Time              |")
	fmt.Println("  ├─────────────────┼────────────┼────────────┼───────────────────--┤───────────────────|")

	for _, t := range a.History {
		sign := "+"
		if t.Type == "WITHDRAW" || t.Type == "TRANSFER_OUT" {
			sign = "-"
		}
		// Truncate note if needed
		note := t.Note
		if len(note) > 17 {
			note = note[:14] + "..."
		}
		fmt.Printf("  │  %-14s │ %s%9.2f │ %10.2f │ %-17s │ %-16s |\n",
			t.Type, sign, t.Amount, t.Balance, note, t.Timestamp.Format("2006-01-02 15:04"))
	}

	fmt.Println("  ├" + line + "┤")
	fmt.Printf("  │  Current Balance: $%-43.2f│\n", a.Balance)
	fmt.Println("  └" + line + "┘")
}

type BaseShape struct {
	Color  string
	Filled bool
}

// Describe belongs to BaseShape — Circle and Rect inherit it automatically!
func (b BaseShape) Describe() string {
	filled := "not filled"
	if b.Filled {
		filled = "filled"
	}
	return fmt.Sprintf("color=%s, %s", b.Color, filled)
}

type Circle struct {
	BaseShape        // embedded — Circle inherits Describe()!
	Radius    float64
}

func (c Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }

type Rect struct {
	BaseShape              // ← semicolon NOT comma for embedding!
	Width, Height float64  // Width and Height are regular fields sharing float64
}

func (r Rect) Area() float64      { return r.Width * r.Height }
func (r Rect) Perimeter() float64 { return 2 * (r.Width + r.Height) }

type Config struct {
	AppName  string `json:"app_name"`
	Port     int    `json:"port"`
	Debug    bool   `json:"debug"`           // ← was missing quotes around "debug"
	DBConfig struct {
		Host string `json:"host"`
		Port int    `json:"port"`
		Name string `json:"name"`
	} `json:"database"`
}

// LoadConfig parses a JSON string and returns a Config
func LoadConfig(data string) (Config, error) {   // ← "json" is a package, not a type!
	var c Config
	err := json.Unmarshal([]byte(data), &c)       // ← Unmarshal reads JSON INTO the struct
	if err != nil {
		return Config{}, err
	}
	return c, nil
}

// SaveConfig serializes a Config to pretty-printed JSON string
func SaveConfig(c Config) (string, error) {
	jsonBytes, err := json.MarshalIndent(c, "", "  ")  // ← Marshal converts struct TO JSON
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
// ============================================================================
// 🏋️ EXERCISES:
//
// 1. Add a Statement() method to BankAccount that prints a nicely formatted
//    transaction history with date, type, amount, and running balance.
//    (Already solved above! Extend it: add date formatting, filter by type)
//
// 2. Create a Shape hierarchy using embedding:
//      type BaseShape struct { Color string; Filled bool }
//      type Circle struct { BaseShape; Radius float64 }
//      type Rect struct { BaseShape; Width, Height float64 }
//    Add Area() and Perimeter() methods to both. Add a Describe() method
//    to BaseShape that both shapes inherit.
//
// 3. Create a Config struct with JSON tags:
//      type Config struct {
//          AppName  string `json:"app_name"`
//          Port     int    `json:"port"`
//          Debug    bool   `json:"debug"`
//          DBConfig struct { ... } `json:"database"`
//      }
//    Write LoadConfig() that parses a JSON string and returns a Config.
//    Write SaveConfig() that serializes a Config to pretty-printed JSON.
//
// 4. [BONUS] Implement the Builder pattern for a complex struct:
//      type QueryBuilder struct { ... }
//      qb := NewQueryBuilder("users").
//          Select("name", "email").
//          Where("age > ?", 18).
//          OrderBy("name").
//          Limit(10).
//          Build()
//    The Build() method should return the SQL string.
//
// ============================================================================

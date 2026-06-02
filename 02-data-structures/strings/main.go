// ============================================================================
// MODULE 2.3: Strings & Runes
// ============================================================================
//
// KEY CONCEPTS:
//   - Strings are IMMUTABLE sequences of bytes (usually UTF-8 encoded)
//   - len() returns BYTE count, not character count!
//   - Use []rune for character-level operations
//   - range on a string iterates RUNES (Unicode code points), not bytes
//   - strings.Builder is the efficient way to build strings (not +=)
//   - strconv handles string ↔ number conversions
//   - The strings package has everything you need for string manipulation
//
// RUN: go run main.go
// ============================================================================

package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	// -----------------------------------------------------------------------
	// PART 1: String Basics
	// -----------------------------------------------------------------------

	fmt.Println("=== String Basics ===")

	// Strings in Go are immutable sequences of bytes
	greeting := "Hello, World!"
	fmt.Println(greeting)
	fmt.Printf("Type: %T, Length (bytes): %d\n", greeting, len(greeting))

	// Individual bytes can be accessed with indexing
	fmt.Printf("First byte: %c (value: %d)\n", greeting[0], greeting[0])
	// But you CANNOT modify them:
	// greeting[0] = 'h'  // ❌ COMPILE ERROR: strings are immutable

	// Multi-line strings with backticks (raw strings)
	rawStr := `This is a raw string.
It preserves newlines and "quotes" literally.
No escape sequences like \n are processed.
Great for regex, JSON templates, SQL queries!`
	fmt.Println("\nRaw string:")
	fmt.Println(rawStr)

	// Escape sequences in regular strings
	escaped := "Tab:\there\nNewline above\t\"quotes\" and \\backslash"
	fmt.Println("\nEscaped string:")
	fmt.Println(escaped)

	// String comparison
	fmt.Println("\n\"Go\" == \"Go\":", "Go" == "Go")
	fmt.Println("\"Go\" < \"Python\":", "Go" < "Python") // Lexicographic

	// -----------------------------------------------------------------------
	// PART 2: Strings vs Bytes vs Runes
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Strings vs Bytes vs Runes ===")

	// THIS IS CRUCIAL: strings are bytes, but characters can be multi-byte!
	ascii := "Hello"
	emoji := "Hello 🌍🚀"
	chinese := "你好世界" // "Hello World" in Chinese

	fmt.Printf("%-12s → len()=%2d bytes, runes=%2d\n",
		`"Hello"`, len(ascii), len([]rune(ascii)))
	fmt.Printf("%-12s → len()=%2d bytes, runes=%2d\n",
		`"Hello 🌍🚀"`, len(emoji), len([]rune(emoji)))
	fmt.Printf("%-12s → len()=%2d bytes, runes=%2d\n",
		`"你好世界"`, len(chinese), len([]rune(chinese)))

	// See the difference: byte iteration vs rune iteration
	s := "Go🚀"
	fmt.Printf("\nString: %q\n", s)

	fmt.Println("Byte-by-byte (indexing):")
	for i := 0; i < len(s); i++ {
		fmt.Printf("  [%d] byte=0x%02X (%c)\n", i, s[i], s[i])
	}

	fmt.Println("Rune-by-rune (range):")
	for i, r := range s {
		fmt.Printf("  [%d] rune=U+%04X (%c) — %d bytes\n",
			i, r, r, len(string(r)))
	}

	// Converting between types
	bytes := []byte("Hello")
	fmt.Printf("\n[]byte: %v → string: %s\n", bytes, string(bytes))

	runes := []rune("Hello 🌍")
	fmt.Printf("[]rune: %v\n", runes)
	fmt.Printf("Rune count: %d, last rune: %c\n", len(runes), runes[len(runes)-1])

	// -----------------------------------------------------------------------
	// PART 3: The strings Package
	// -----------------------------------------------------------------------

	fmt.Println("\n=== The strings Package ===")

	text := "  Go is an Open Source programming language  "

	// Searching
	fmt.Println("Contains 'Open':", strings.Contains(text, "Open"))
	fmt.Println("Contains 'open':", strings.Contains(text, "open")) // Case-sensitive!
	fmt.Println("HasPrefix '  Go':", strings.HasPrefix(text, "  Go"))
	fmt.Println("HasSuffix 'age  ':", strings.HasSuffix(text, "age  "))
	fmt.Println("Index of 'Open':", strings.Index(text, "Open"))
	fmt.Println("Count of 'a':", strings.Count(text, "a"))

	// Transforming
	fmt.Println("\nToUpper:", strings.ToUpper("hello"))
	fmt.Println("ToLower:", strings.ToLower("HELLO"))
	fmt.Println("Title:", strings.ToTitle("hello world"))
	fmt.Println("TrimSpace:", fmt.Sprintf("%q", strings.TrimSpace(text)))
	fmt.Println("Trim dots:", strings.Trim("...hello...", "."))
	fmt.Println("Replace:", strings.Replace("foo bar foo baz foo", "foo", "qux", 2))
	fmt.Println("ReplaceAll:", strings.ReplaceAll("aabbcc", "b", "X"))

	// Splitting and Joining
	csv := "apple,banana,cherry,date"
	parts := strings.Split(csv, ",")
	fmt.Println("\nSplit:", parts)
	fmt.Println("Join:", strings.Join(parts, " | "))

	// Fields splits by ANY whitespace (better than Split for text)
	messy := "  hello   world   go   "
	fields := strings.Fields(messy)
	fmt.Println("Fields:", fields) // ["hello", "world", "go"]

	// Repeat
	fmt.Println("Repeat:", strings.Repeat("Go! ", 3))

	// Map — transform each rune
	shouted := strings.Map(func(r rune) rune {
		if r == ' ' {
			return '!'
		}
		return unicode.ToUpper(r)
	}, "hello world")
	fmt.Println("Map:", shouted) // HELLO!WORLD

	// -----------------------------------------------------------------------
	// PART 4: String Building
	// -----------------------------------------------------------------------

	fmt.Println("\n=== String Building ===")

	// ❌ BAD: += creates a new string every time — O(n²)!
	// result := ""
	// for i := 0; i < 1000; i++ {
	//     result += "x"  // Creates 1000 intermediate strings!
	// }

	// ✅ GOOD: strings.Builder — efficient, O(n)
	var builder strings.Builder
	for i := 0; i < 10; i++ {
		builder.WriteString("Go")
		if i < 9 {
			builder.WriteString(", ")
		}
	}
	fmt.Println("Builder:", builder.String())

	// Builder with various write methods
	var b strings.Builder
	b.WriteString("Name: ")
	b.WriteString("Naman")
	b.WriteByte('\n')
	b.WriteString("Score: ")
	b.WriteRune('💯')
	fmt.Println("\nMulti-write builder:")
	fmt.Print(b.String())

	// fmt.Sprintf for formatted strings
	name := "Naman"
	age := 23
	formatted := fmt.Sprintf("Hello, %s! You are %d years old.", name, age)
	fmt.Println("\n\nSprintf:", formatted)

	// -----------------------------------------------------------------------
	// PART 5: String Conversions (strconv)
	// -----------------------------------------------------------------------

	fmt.Println("\n=== String Conversions (strconv) ===")

	// Int ↔ String
	num := 42
	numStr := strconv.Itoa(num) // Int to ASCII
	fmt.Printf("Itoa: %d → %q (type: %T)\n", num, numStr, numStr)

	parsed, err := strconv.Atoi("123") // ASCII to Int
	if err == nil {
		fmt.Printf("Atoi: %q → %d (type: %T)\n", "123", parsed, parsed)
	}

	_, err = strconv.Atoi("not a number")
	if err != nil {
		fmt.Println("Atoi error:", err)
	}

	// Float ↔ String
	f := 3.14159
	fStr := strconv.FormatFloat(f, 'f', 2, 64) // format, precision, bitSize
	fmt.Printf("FormatFloat: %f → %q\n", f, fStr)

	fParsed, _ := strconv.ParseFloat("2.71828", 64)
	fmt.Printf("ParseFloat: %q → %f\n", "2.71828", fParsed)

	// Bool ↔ String
	bStr := strconv.FormatBool(true)
	fmt.Printf("FormatBool: %t → %q\n", true, bStr)

	bParsed, _ := strconv.ParseBool("true")
	fmt.Printf("ParseBool: %q → %t\n", "true", bParsed)

	// Quick tip: fmt.Sprintf works too, but strconv is faster
	quick := fmt.Sprintf("%d", 42) // Slower but more flexible
	_ = quick

	// -----------------------------------------------------------------------
	// PART 6: Unicode & Rune Handling
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Unicode & Rune Handling ===")

	// The unicode package has character classification functions
	testRunes := []rune{'A', 'a', '1', ' ', '!', '你', '🚀'}

	fmt.Println("  Char  │ Letter │ Digit  │ Upper  │ Space  │ Punct")
	fmt.Println("────────┼────────┼────────┼────────┼────────┼───────")
	for _, r := range testRunes {
		fmt.Printf("  %-4c  │ %-6t │ %-6t │ %-6t │ %-6t │ %t\n",
			r,
			unicode.IsLetter(r),
			unicode.IsDigit(r),
			unicode.IsUpper(r),
			unicode.IsSpace(r),
			unicode.IsPunct(r))
	}

	// Rune conversions
	fmt.Println("\nToUpper('a'):", string(unicode.ToUpper('a')))
	fmt.Println("ToLower('Z'):", string(unicode.ToLower('Z')))

	// Count characters by category
	sample := "Hello, World! 你好 123 🎉"
	var letters, digits, spaces, other int
	for _, r := range sample {
		switch {
		case unicode.IsLetter(r):
			letters++
		case unicode.IsDigit(r):
			digits++
		case unicode.IsSpace(r):
			spaces++
		default:
			other++
		}
	}
	fmt.Printf("\nIn %q:\n", sample)
	fmt.Printf("  Letters: %d, Digits: %d, Spaces: %d, Other: %d\n",
		letters, digits, spaces, other)

	// -----------------------------------------------------------------------
	// PART 7: Regular Expressions
	// -----------------------------------------------------------------------

	fmt.Println("\n=== Regular Expressions ===")

	// regexp.Compile returns a compiled regex and an error
	// regexp.MustCompile panics on invalid regex (use for constants)
	emailPattern := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)

	// MatchString — does it match?
	fmt.Println("Match 'user@example.com':", emailPattern.MatchString("user@example.com"))
	fmt.Println("Match 'not-an-email':", emailPattern.MatchString("not-an-email"))

	// FindString — find first match
	mixedText := "Contact us at support@golearn.dev or admin@golearn.dev"
	first := emailPattern.FindString(mixedText)
	fmt.Println("First email:", first)

	// FindAllString — find all matches
	allEmails := emailPattern.FindAllString(mixedText, -1) // -1 = find all
	fmt.Println("All emails:", allEmails)

	// ReplaceAllString — replace matches
	redacted := emailPattern.ReplaceAllString(mixedText, "[REDACTED]")
	fmt.Println("Redacted:", redacted)

	// Capture groups with FindStringSubmatch
	urlPattern := regexp.MustCompile(`(https?)://([^/]+)(/\S*)?`)
	url := "Visit https://golang.org/doc for docs"
	matches := urlPattern.FindStringSubmatch(url)
	if matches != nil {
		fmt.Println("\nURL breakdown:")
		fmt.Println("  Full match:", matches[0])
		fmt.Println("  Protocol:", matches[1])
		fmt.Println("  Domain:", matches[2])
		fmt.Println("  Path:", matches[3])
	}

	// Split by regex
	wordSplitter := regexp.MustCompile(`\W+`)
	wordList := wordSplitter.Split("Hello, World! How's it going?", -1)
	fmt.Println("\nWords:", wordList)

	// -----------------------------------------------------------------------
	// PART 8: Practical Exercise — Text Analyzer (SOLVED)
	// -----------------------------------------------------------------------

	fmt.Println("\n=== 📝 Text Analyzer ===")

	passage := `Go is an open source programming language that makes it easy to
build simple, reliable, and efficient software. Go was designed at Google
by Robert Griesemer, Rob Pike, and Ken Thompson. Go is sometimes called
Golang because of its domain name, golang.org. Go has built-in concurrency
and a robust standard library. Go is great for building web services,
command-line tools, and distributed systems.`

	stats := analyzeText(passage)

	fmt.Println("─────────────────────────────────────")
	fmt.Println("           Text Statistics           ")
	fmt.Println("─────────────────────────────────────")
	fmt.Printf("  Characters (with spaces):  %d\n", stats.chars)
	fmt.Printf("  Characters (no spaces):    %d\n", stats.charsNoSpace)
	fmt.Printf("  Words:                     %d\n", stats.words)
	fmt.Printf("  Sentences:                 %d\n", stats.sentences)
	fmt.Printf("  Average word length:       %.1f\n", stats.avgWordLen)
	fmt.Printf("  Unique words:              %d\n", stats.uniqueWords)
	fmt.Println("─────────────────────────────────────")

	fmt.Println("\n  Top 5 most common words:")
	for i, wc := range stats.topWords {
		bar := strings.Repeat("█", wc.count)
		fmt.Printf("    %d. %-12s %2d  %s\n", i+1, wc.word, wc.count, bar)
	}
}

// -----------------------------------------------------------------------
// Text Analyzer Implementation
// -----------------------------------------------------------------------

type wordFreq struct {
	word  string
	count int
}

type textStats struct {
	chars        int
	charsNoSpace int
	words        int
	sentences    int
	avgWordLen   float64
	uniqueWords  int
	topWords     []wordFreq
}

func analyzeText(text string) textStats {
	var stats textStats

	// Character counts
	stats.chars = len([]rune(text))
	for _, r := range text {
		if !unicode.IsSpace(r) {
			stats.charsNoSpace++
		}
	}

	// Word counts
	wordList := strings.Fields(text)
	stats.words = len(wordList)

	// Sentence count (naive: count '.', '!', '?')
	for _, r := range text {
		if r == '.' || r == '!' || r == '?' {
			stats.sentences++
		}
	}

	// Word frequencies
	freq := make(map[string]int)
	totalLen := 0
	for _, word := range wordList {
		// Clean punctuation from word edges
		cleaned := strings.Trim(strings.ToLower(word), ".,!?;:'\"()[]")
		if cleaned != "" {
			freq[cleaned]++
			totalLen += len([]rune(cleaned))
		}
	}

	stats.uniqueWords = len(freq)
	if stats.words > 0 {
		stats.avgWordLen = float64(totalLen) / float64(stats.words)
	}

	// Sort by frequency
	var sorted []wordFreq
	for w, c := range freq {
		sorted = append(sorted, wordFreq{w, c})
	}
	// Sort descending by count, then alphabetically
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].count > sorted[i].count ||
				(sorted[j].count == sorted[i].count && sorted[j].word < sorted[i].word) {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	if len(sorted) > 5 {
		stats.topWords = sorted[:5]
	} else {
		stats.topWords = sorted
	}

	return stats
}

// ============================================================================
// 🏋️ EXERCISES:
//
// 1. Write isPalindrome(s string) bool that checks if a string reads the
//    same forwards and backwards, ignoring case, spaces, and punctuation.
//    Examples: "racecar" → true, "A man a plan a canal Panama" → true
//    Hint: convert to []rune, filter only letters, compare in reverse
//
// 2. Write caesarCipher(s string, shift int) string that shifts each letter
//    by 'shift' positions in the alphabet. Handle wraparound (z+1=a),
//    preserve case, and leave non-letters unchanged.
//    Example: caesarCipher("Hello", 3) → "Khoor"
//
// 3. Write truncate(s string, maxLen int) string that truncates a string
//    to maxLen CHARACTERS (not bytes!) and appends "..." if truncated.
//    Example: truncate("Hello, 世界!", 7) → "Hello, …"
//    Hint: work with []rune, not bytes
//
// 4. [BONUS] Write a simple Markdown stripper that converts Markdown to
//    plain text: remove # headers, **bold**, *italic*, [text](url) links
//    (keep text, remove URL), and ``` code blocks.
//    Hint: use regexp.ReplaceAllString for each pattern
//
// ============================================================================

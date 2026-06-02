# 🚀 Learn Go — Hands-On Workspace

Welcome to your Go learning workspace! This is a structured, progressive curriculum
where every concept comes with working code you can run, modify, and experiment with.

## Getting Started

```bash
# You have Go 1.22.5 installed. Verify with:
go version

# Navigate to any lesson and run it:
cd 01-foundations/hello
go run main.go
go run main.go YourName
```

## Directory Structure

```
Go/
├── 01-foundations/       ← START HERE
│   ├── hello/            Module 1.1: Hello World & CLI args
│   ├── variables/        Module 1.2: Types, zero values, constants, iota
│   ├── control-flow/     Module 1.3: if/for/switch, FizzBuzz, guessing game
│   └── functions/        Module 1.4: Multiple returns, closures, defer
├── 02-data-structures/
│   ├── slices/           Module 2.1: Arrays, slices, append, copy, internals
│   ├── maps/             Module 2.2: Maps, comma-ok, sets, nested maps
│   ├── strings/          Module 2.3: Strings, runes, Builder, strconv, regex
│   └── structs/          Module 2.4: Structs, methods, embedding, JSON tags
├── 03-interfaces/
├── 04-packages/
├── 05-concurrency/       ⭐ Go's killer feature
├── 06-stdlib/
├── 07-testing/
├── 08-advanced/
└── 09-projects/          🚀 Capstone projects
```

## How to Work Through Each Lesson

1. **Read** the comments at the top — they explain the key concepts
2. **Run** the code: `go run main.go`
3. **Modify** things and re-run — experiment!
4. **Do the exercises** at the bottom of each file (marked with 🏋️)
5. **Move on** to the next lesson when you're comfortable

## Essential Go Commands

| Command              | What it does                 |
|----------------------|------------------------------|
| `go run main.go`     | Compile and run in one step  |
| `go build`           | Compile to binary            |
| `go fmt ./...`       | Auto-format all code         |
| `go vet ./...`       | Find common mistakes         |
| `go test ./...`      | Run all tests                |
| `go mod tidy`        | Clean up dependencies        |
| `go doc fmt.Println` | View docs for any function   |

## Tips

- 💡 Go enforces formatting — run `go fmt` and never argue about style again
- 💡 Unused imports and variables cause **compile errors** (not warnings!)
- 💡 Read the error messages — Go's compiler errors are excellent
- 💡 When stuck, try `go doc <package>.<Function>` for instant docs

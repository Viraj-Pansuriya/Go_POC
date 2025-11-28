# Go Mono Repo - Learning Journey 🚀

Welcome to your Go learning monorepo! This repository is structured to help you transition from Java/C++ to Go.

## 📁 Repository Structure

```
go-mono-repo/
├── notes/                    # Quick reference guides (cheatsheets)
│   ├── 01-basics.md          # Variables, types, control flow
│   ├── 02-functions.md       # Functions, methods, interfaces
│   ├── 03-concurrency.md     # Goroutines, channels
│   ├── 04-error-handling.md  # Error handling patterns
│   └── 05-packages.md        # Modules, packages, imports
│
├── docs/                     # Project-specific documentation
│   └── projects/             # One .md file per project you create
│
├── projects/                 # Your Go projects live here
│   └── hello-world/          # Sample starter project
│
├── go.work                   # Go workspace file (monorepo magic)
└── README.md                 # You are here!
```

## 🎯 Quick Start

```bash
# Navigate to any project
cd projects/hello-world

# Run the project
go run .

# Build the project
go build -o bin/hello-world

# Run tests
go test ./...
```

## 🔄 Java/C++ → Go Mental Model

| Concept | Java | C++ | Go |
|---------|------|-----|-----|
| Entry Point | `public static void main(String[] args)` | `int main()` | `func main()` |
| Package Manager | Maven/Gradle | CMake/Conan | `go mod` (built-in) |
| Classes | `class Foo {}` | `class Foo {};` | `type Foo struct {}` (no classes!) |
| Inheritance | `extends` | `: public Base` | **Composition over inheritance** |
| Interfaces | `implements` | virtual functions | Implicit (duck typing) |
| Memory | GC | Manual/RAII | GC (like Java) |
| Threads | `Thread`, `ExecutorService` | `std::thread` | `goroutines` (lightweight) |
| Error Handling | Exceptions | Exceptions | Return values (`error` type) |

## 🧭 Learning Path

1. **Week 1**: Basics (syntax, types, control flow)
2. **Week 2**: Functions, methods, and structs
3. **Week 3**: Interfaces and composition
4. **Week 4**: Concurrency (goroutines, channels)
5. **Week 5**: Error handling and testing
6. **Week 6**: Building real projects

## 💡 Go Philosophy (Different from Java/C++)

1. **Simplicity over features** - No generics abuse, no class hierarchies
2. **Composition over inheritance** - Embed structs, don't extend
3. **Explicit over implicit** - Errors are returned, not thrown
4. **Concurrency is first-class** - goroutines are dirt cheap
5. **One way to do things** - `gofmt` standardizes formatting

---

*Happy coding! Check the `notes/` folder for quick references.*



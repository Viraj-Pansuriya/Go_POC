# 📚 Quick Reference Notes

Fast lookup guides for Go concepts. Coming from Java/C++, these notes highlight the differences and mental model shifts.

## 📋 Index

| Note | Topic | Key Concepts |
|------|-------|--------------|
| [01-basics.md](./01-basics.md) | Variables, Types, Control Flow | `:=` syntax, no semicolons, `for` is the only loop |
| [02-functions.md](./02-functions.md) | Functions, Methods, Interfaces | Multiple returns, pointer receivers, duck typing |
| [03-concurrency.md](./03-concurrency.md) | Goroutines, Channels | `go` keyword, channels, `select`, `sync` package |
| [04-error-handling.md](./04-error-handling.md) | Error Patterns | No exceptions, `error` return values, `panic`/`recover` |
| [05-packages.md](./05-packages.md) | Modules, Imports | `go mod`, visibility rules, project structure |

## 🔥 Most Common Gotchas (Java/C++ → Go)

### 1. No Classes, No Inheritance
```go
// Not this (Java/C++)
class Dog extends Animal { }

// But this (Go)
type Dog struct {
    Animal  // Embed for composition
}
```

### 2. Visibility by Naming
```go
func PublicFunc() {}   // Exported (like public)
func privateFunc() {}  // Unexported (like private)
```

### 3. Error Handling is Explicit
```go
// Not this
try { } catch (Exception e) { }

// But this
result, err := something()
if err != nil {
    return err
}
```

### 4. Interfaces are Implicit
```go
// No "implements" keyword!
// If your struct has the methods, it implements the interface
type Writer interface {
    Write([]byte) (int, error)
}

// MyWriter implements Writer automatically if it has Write()
```

### 5. Goroutines are Not Threads
```go
// Not this (Java)
new Thread(() -> doWork()).start();

// But this (Go) - lightweight!
go doWork()  // Can spawn millions
```

---

## 🚀 One-Liner Cheat Sheet

```go
// Variable
x := 42

// Function
func add(a, b int) int { return a + b }

// Struct
type User struct { Name string }

// Method
func (u User) Greet() string { return "Hi, " + u.Name }

// Goroutine
go doSomething()

// Channel
ch := make(chan int)
ch <- 42     // send
x := <-ch    // receive

// Error check
if err != nil { return err }

// Slice operations
s = append(s, item)
s = s[1:3]

// Map
m := map[string]int{"a": 1}
val, ok := m["key"]

// JSON
json.Marshal(obj)
json.Unmarshal(data, &obj)
```

---

*Add your own notes as you learn! 📝*



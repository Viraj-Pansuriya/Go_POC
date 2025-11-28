# Go Basics - Quick Reference

## 📦 Package Declaration

Every Go file starts with a package declaration.

```go
package main  // Executable program (has main function)
package utils // Library package (no main function)
```

**Java equivalent:** `package com.example.utils;`  
**C++ equivalent:** Namespaces (`namespace utils {}`)

---

## 🔤 Variables

### Declaration Styles

```go
// Method 1: var keyword (like Java final type inference)
var name string = "Viraj"
var age int = 25

// Method 2: Short declaration (most common, inside functions only)
name := "Viraj"    // Type inferred
age := 25

// Method 3: Multiple variables
var x, y int = 1, 2
a, b := "hello", 42

// Zero values (like default values in Java)
var i int      // 0
var s string   // "" (empty string)
var b bool     // false
var p *int     // nil (like null in Java)
```

### Comparison Table

| Feature | Java | C++ | Go |
|---------|------|-----|-----|
| Declaration | `int x = 5;` | `int x = 5;` | `x := 5` or `var x int = 5` |
| Type inference | `var x = 5;` (Java 10+) | `auto x = 5;` | `x := 5` |
| Constants | `final int X = 5;` | `const int X = 5;` | `const X = 5` |
| No value | `null` | `nullptr` | `nil` |

---

## 📊 Basic Types

```go
// Numeric
int, int8, int16, int32, int64
uint, uint8, uint16, uint32, uint64
float32, float64
complex64, complex128

// Other
bool
string
byte    // alias for uint8
rune    // alias for int32 (Unicode code point)
```

**Tradeoff:** Go has explicit sized types (`int32`, `int64`). In Java, `int` is always 32-bit. Go's `int` is platform-dependent (32 or 64 bit).

---

## 🔁 Control Flow

### If Statement

```go
// No parentheses required!
if x > 10 {
    fmt.Println("big")
} else if x > 5 {
    fmt.Println("medium")
} else {
    fmt.Println("small")
}

// With initialization (unique to Go!)
if err := doSomething(); err != nil {
    // err is only scoped here
    return err
}
```

**Java equivalent:**
```java
if (x > 10) { ... }  // Parentheses required
```

### For Loop (The ONLY loop in Go!)

```go
// Classic for loop
for i := 0; i < 10; i++ {
    fmt.Println(i)
}

// While-style loop
for x < 100 {
    x *= 2
}

// Infinite loop
for {
    // break to exit
}

// Range loop (like for-each)
nums := []int{1, 2, 3}
for index, value := range nums {
    fmt.Println(index, value)
}

// Ignore index with _
for _, value := range nums {
    fmt.Println(value)
}
```

**Tradeoff:** No `while` or `do-while` keywords. Go uses `for` for everything - simpler but different mindset.

### Switch Statement

```go
// No break needed! (automatic)
switch day {
case "Monday":
    fmt.Println("Start of week")
case "Friday":
    fmt.Println("TGIF!")
default:
    fmt.Println("Regular day")
}

// Switch without expression (cleaner if-else chain)
switch {
case x < 0:
    fmt.Println("negative")
case x > 0:
    fmt.Println("positive")
default:
    fmt.Println("zero")
}
```

**Java/C++ difference:** In Java/C++, switch falls through by default. In Go, it breaks by default. Use `fallthrough` keyword if needed.

---

## 📝 Strings

```go
// String declaration
s := "Hello, World"

// String length
len(s)  // 12 (bytes, not runes!)

// String concatenation
greeting := "Hello, " + "World"

// Multi-line strings (raw strings)
query := `
    SELECT *
    FROM users
    WHERE id = 1
`

// String formatting (like printf)
name := "Viraj"
msg := fmt.Sprintf("Hello, %s!", name)
```

**Tradeoff:** Strings in Go are immutable (like Java). Concatenation creates new strings. Use `strings.Builder` for efficient building.

---

## 🧮 Arrays and Slices

```go
// Array (fixed size - rarely used directly)
var arr [5]int = [5]int{1, 2, 3, 4, 5}

// Slice (dynamic size - used 99% of the time)
slice := []int{1, 2, 3, 4, 5}

// Create slice with make
slice := make([]int, 5)      // length 5
slice := make([]int, 5, 10)  // length 5, capacity 10

// Append to slice
slice = append(slice, 6, 7, 8)

// Slice a slice
subSlice := slice[1:4]  // elements 1, 2, 3
```

**Java equivalent:** `ArrayList<Integer>`  
**C++ equivalent:** `std::vector<int>`

**Tradeoff:** Slices are references (like ArrayList). Modifying a slice can affect the original array. This is a common gotcha!

---

## 🗺️ Maps

```go
// Create a map
ages := map[string]int{
    "Alice": 30,
    "Bob":   25,
}

// Create empty map with make
ages := make(map[string]int)

// Set value
ages["Charlie"] = 35

// Get value (with existence check)
age, exists := ages["Alice"]
if exists {
    fmt.Println(age)
}

// Delete key
delete(ages, "Bob")

// Iterate
for name, age := range ages {
    fmt.Println(name, age)
}
```

**Java equivalent:** `HashMap<String, Integer>`  
**C++ equivalent:** `std::unordered_map<std::string, int>`

---

## 🎯 Quick Tips for Java/C++ Developers

1. **No semicolons** - Go inserts them automatically
2. **No parentheses** in if/for/switch conditions
3. **Braces are mandatory** - even for single-line blocks
4. **Unused variables = compile error** - use `_` to ignore
5. **Unused imports = compile error** - keep it clean!
6. **Public = Capitalized** - `func Public()` vs `func private()`

---

*Next: [02-functions.md](./02-functions.md)*



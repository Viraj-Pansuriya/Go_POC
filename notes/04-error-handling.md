# Error Handling in Go - Quick Reference

## 🎯 The Big Difference

| Aspect | Java/C++ | Go |
|--------|----------|-----|
| Mechanism | Exceptions (throw/try/catch) | Return values |
| Control flow | Jumps up call stack | Explicit, local |
| Philosophy | Exceptional cases | Errors are values |
| Checked errors | Java: Yes, C++: No | Yes (compiler checks return) |

**Go Philosophy:** Errors are not exceptional. Handle them explicitly.

---

## 📌 Basic Error Handling

```go
import "errors"

// Function that returns an error
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// Calling and handling error
result, err := divide(10, 0)
if err != nil {
    log.Fatal(err)  // or handle appropriately
}
fmt.Println(result)
```

### Java/C++ Equivalent

```java
// Java
public int divide(int a, int b) throws ArithmeticException {
    if (b == 0) throw new ArithmeticException("division by zero");
    return a / b;
}

try {
    int result = divide(10, 0);
} catch (ArithmeticException e) {
    System.err.println(e.getMessage());
}
```

```cpp
// C++
int divide(int a, int b) {
    if (b == 0) throw std::runtime_error("division by zero");
    return a / b;
}

try {
    int result = divide(10, 0);
} catch (const std::exception& e) {
    std::cerr << e.what() << std::endl;
}
```

---

## 🛠️ Creating Errors

### Simple Errors

```go
import "errors"

// Method 1: errors.New
err := errors.New("something went wrong")

// Method 2: fmt.Errorf (with formatting)
err := fmt.Errorf("failed to process user %d", userID)
```

### Custom Error Types

```go
// Define custom error type
type ValidationError struct {
    Field   string
    Message string
}

// Implement error interface (just needs Error() string)
func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Usage
func validateAge(age int) error {
    if age < 0 {
        return &ValidationError{
            Field:   "age",
            Message: "must be non-negative",
        }
    }
    return nil
}
```

**Java equivalent:** Custom exception class extending `Exception`

---

## 🔗 Error Wrapping (Go 1.13+)

```go
import (
    "errors"
    "fmt"
)

// Wrap error with context
func readConfig(path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return fmt.Errorf("reading config %s: %w", path, err)
    }
    // ... process data
    return nil
}

// Unwrap and check error type
err := readConfig("/bad/path")

// Check if it's a specific error
if errors.Is(err, os.ErrNotExist) {
    fmt.Println("File not found")
}

// Extract wrapped error type
var pathErr *os.PathError
if errors.As(err, &pathErr) {
    fmt.Println("Path:", pathErr.Path)
}
```

### Error Chain

```go
// Original error
originalErr := errors.New("database connection failed")

// Wrapped once
err1 := fmt.Errorf("user service: %w", originalErr)

// Wrapped again
err2 := fmt.Errorf("API handler: %w", err1)

// err2.Error() = "API handler: user service: database connection failed"

// errors.Is checks the entire chain
errors.Is(err2, originalErr)  // true
```

---

## 🎭 Sentinel Errors (Predefined Errors)

```go
import "errors"

// Define sentinel errors (usually at package level)
var (
    ErrNotFound     = errors.New("not found")
    ErrUnauthorized = errors.New("unauthorized")
    ErrInvalidInput = errors.New("invalid input")
)

func GetUser(id int) (*User, error) {
    user := db.Find(id)
    if user == nil {
        return nil, ErrNotFound
    }
    return user, nil
}

// Caller checks for specific errors
user, err := GetUser(123)
if errors.Is(err, ErrNotFound) {
    // Handle not found case
}
```

**Java equivalent:** Specific exception types like `NotFoundException`

---

## 🏗️ Error Handling Patterns

### Pattern 1: Early Return

```go
func processRequest(r *Request) error {
    if r == nil {
        return errors.New("request is nil")
    }
    
    if err := validateRequest(r); err != nil {
        return fmt.Errorf("validation: %w", err)
    }
    
    if err := saveRequest(r); err != nil {
        return fmt.Errorf("saving: %w", err)
    }
    
    return nil
}
```

### Pattern 2: Defer with Error

```go
func writeFile(path string, data []byte) (err error) {
    f, err := os.Create(path)
    if err != nil {
        return err
    }
    
    defer func() {
        closeErr := f.Close()
        if err == nil {
            err = closeErr  // Capture close error if no prior error
        }
    }()
    
    _, err = f.Write(data)
    return err
}
```

### Pattern 3: Error Aggregation

```go
type MultiError struct {
    Errors []error
}

func (m *MultiError) Error() string {
    msgs := make([]string, len(m.Errors))
    for i, err := range m.Errors {
        msgs[i] = err.Error()
    }
    return strings.Join(msgs, "; ")
}

func validateAll(items []Item) error {
    var errs []error
    for _, item := range items {
        if err := validate(item); err != nil {
            errs = append(errs, err)
        }
    }
    if len(errs) > 0 {
        return &MultiError{Errors: errs}
    }
    return nil
}
```

---

## 🔥 Panic and Recover (Use Sparingly!)

**Panic is like an exception - but only for truly exceptional cases.**

```go
// Panic - crashes the program (unless recovered)
func mustParse(s string) int {
    n, err := strconv.Atoi(s)
    if err != nil {
        panic(fmt.Sprintf("failed to parse %q: %v", s, err))
    }
    return n
}

// Recover - catches panic (only works in defer)
func safeCall(fn func()) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic recovered: %v", r)
        }
    }()
    
    fn()
    return nil
}
```

### When to Panic

| Panic ✓ | Don't Panic ✗ |
|---------|---------------|
| Programmer error (bug) | User input error |
| Impossible state | File not found |
| Init/setup failure | Network error |
| Assertion failure | Validation error |

```go
// OK to panic - programmer error
func (s *Server) MustStart() {
    if s.config == nil {
        panic("server config is nil")  // Bug in code
    }
}

// NOT OK to panic - handle gracefully
func (s *Server) HandleRequest(r *Request) error {
    if r.UserID == "" {
        return ErrInvalidRequest  // User error, return error
    }
}
```

---

## 📋 Common Patterns Comparison

### Java try-with-resources vs Go defer

```java
// Java
try (FileInputStream fis = new FileInputStream("file.txt")) {
    // use fis
} catch (IOException e) {
    // handle
}
```

```go
// Go
f, err := os.Open("file.txt")
if err != nil {
    return err
}
defer f.Close()
// use f
```

### Java Optional vs Go (value, ok)

```java
// Java
Optional<User> user = findUser(id);
user.ifPresent(u -> process(u));
```

```go
// Go - maps
value, ok := myMap[key]
if ok {
    process(value)
}

// Go - type assertion
str, ok := value.(string)
if ok {
    process(str)
}
```

---

## 🎯 Best Practices

1. **Always check errors** - don't use `_` to ignore

```go
// BAD
result, _ := mightFail()

// GOOD
result, err := mightFail()
if err != nil {
    return err
}
```

2. **Add context when wrapping**

```go
// BAD
return err

// GOOD
return fmt.Errorf("processing user %d: %w", userID, err)
```

3. **Handle errors once**

```go
// BAD - logging and returning
if err != nil {
    log.Println(err)  // Logged here
    return err        // And will be logged again up the stack
}

// GOOD - either log OR return
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}
```

4. **Use sentinel errors for expected errors**

```go
// Define at package level
var ErrNotFound = errors.New("not found")

// Caller can check specifically
if errors.Is(err, ErrNotFound) {
    // Handle missing item
}
```

---

*Next: [05-packages.md](./05-packages.md)*



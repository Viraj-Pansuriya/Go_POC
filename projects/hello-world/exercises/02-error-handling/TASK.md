# Exercise 02: Error Handling in Go

## 🎯 Goal
Master Go's error handling - NO try/catch, errors are **values** you return and check!

---

## 📚 Key Concept: Java vs Go Error Handling

### Java (Exceptions)
```java
try {
    User user = userService.findById(id);
    user.updateEmail(newEmail);
} catch (UserNotFoundException e) {
    log.error("User not found", e);
} catch (ValidationException e) {
    log.error("Invalid email", e);
} finally {
    // cleanup
}
```

### Go (Errors as Values)
```go
user, err := userService.FindById(id)
if err != nil {
    log.Printf("Failed to find user: %v", err)
    return err
}

err = user.UpdateEmail(newEmail)
if err != nil {
    log.Printf("Failed to update email: %v", err)
    return err
}
```

**Key Insight**: In Go, errors are just values returned from functions. You MUST check them explicitly!

---

## 🏋️ Your Task: Build a User Registration System

### Requirements:

### 1. Create custom error types

```go
// errors.go
package registration

import "fmt"

// Sentinel errors (predefined errors to compare against)
var (
    ErrUserNotFound     = errors.New("user not found")
    ErrEmailExists      = errors.New("email already exists")
    ErrInvalidEmail     = errors.New("invalid email format")
)

// Custom error type with more context
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed on field '%s': %s", e.Field, e.Message)
}
```

### 2. Create a UserService with these methods:

```go
// user_service.go
type User struct {
    ID    int
    Name  string
    Email string
    Age   int
}

type UserService struct {
    users map[int]*User  // in-memory storage
}

func NewUserService() *UserService { ... }

// Register a new user - returns error if validation fails
func (s *UserService) Register(name, email string, age int) (*User, error) { ... }

// FindByID - returns ErrUserNotFound if not found
func (s *UserService) FindByID(id int) (*User, error) { ... }

// UpdateEmail - returns validation error if email invalid
func (s *UserService) UpdateEmail(id int, newEmail string) error { ... }

// Delete - returns error if user doesn't exist
func (s *UserService) Delete(id int) error { ... }
```

### 3. Validation Rules:
- Name: cannot be empty, min 2 characters
- Email: must contain "@" (simple check is fine)
- Age: must be >= 18

### 4. In main(), demonstrate:
- Register valid users
- Try to register with invalid data (show error handling)
- Find existing and non-existing users
- Update email with valid and invalid emails
- Use `errors.Is()` to check for specific errors
- Use `errors.As()` to extract ValidationError details

---

## 📁 Files to Create

```
02-error-handling/
├── main.go
└── registration/
    ├── errors.go       # Custom errors
    └── user_service.go # UserService
```

---

## ✅ Expected Output (example)
```
=== Registering Users ===
✓ Registered: {ID:1 Name:John Doe Email:john@example.com Age:25}
✗ Registration failed: validation failed on field 'name': name must be at least 2 characters
✗ Registration failed: validation failed on field 'email': invalid email format
✗ Registration failed: validation failed on field 'age': must be 18 or older

=== Finding Users ===
✓ Found user: John Doe
✗ Find failed: user not found

=== Updating Email ===
✓ Email updated successfully
✗ Update failed: invalid email format

=== Error Type Checking ===
Is ErrUserNotFound: true
ValidationError details - Field: email, Message: invalid email format
```

---

## 💡 Hints

### 1. Creating errors
```go
import "errors"

// Simple error
err := errors.New("something went wrong")

// Formatted error
err := fmt.Errorf("user %d not found", id)

// Wrapping errors (adds context)
err := fmt.Errorf("failed to update user: %w", originalErr)
```

### 2. Checking error types
```go
import "errors"

// Check if error IS a specific error
if errors.Is(err, ErrUserNotFound) {
    // handle not found
}

// Extract custom error type
var validationErr *ValidationError
if errors.As(err, &validationErr) {
    fmt.Println("Field:", validationErr.Field)
}
```

### 3. The error interface
```go
// Any type with Error() method implements error interface
type error interface {
    Error() string
}
```

---

## 🎓 What You'll Learn

1. **Errors as values** - Return and check, no throwing
2. **Sentinel errors** - Predefined errors for comparison
3. **Custom error types** - Structs that implement error interface
4. **Error wrapping** - Add context with `%w`
5. **errors.Is()** - Check error identity (even through wrapping)
6. **errors.As()** - Extract specific error types

---

## ⚠️ Common Mistakes to Avoid

```go
// ❌ BAD: Ignoring errors
user, _ := userService.FindByID(id)

// ✅ GOOD: Always check
user, err := userService.FindByID(id)
if err != nil {
    return nil, err
}

// ❌ BAD: Comparing errors with ==
if err == ErrUserNotFound { }  // Doesn't work with wrapped errors!

// ✅ GOOD: Use errors.Is
if errors.Is(err, ErrUserNotFound) { }

// ❌ BAD: Type assertion for errors
if validationErr, ok := err.(*ValidationError); ok { }  // Doesn't work with wrapped!

// ✅ GOOD: Use errors.As
var validationErr *ValidationError
if errors.As(err, &validationErr) { }
```

---

## ⏱️ Estimated Time: 25-30 minutes

When done, let me know and I'll review your code! 🚀


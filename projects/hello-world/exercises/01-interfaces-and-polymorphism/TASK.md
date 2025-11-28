# Exercise 01: Interfaces and Polymorphism in Go

## 🎯 Goal
Understand how Go's **implicit interfaces** work - this is VERY different from Java!

---

## 📚 Key Concept: Java vs Go Interfaces

### Java (Explicit)
```java
public interface PaymentProcessor {
    void process(double amount);
}

public class CreditCard implements PaymentProcessor {  // EXPLICIT
    @Override
    public void process(double amount) { ... }
}
```

### Go (Implicit)
```go
type PaymentProcessor interface {
    Process(amount float64) error
}

type CreditCard struct { }

// NO "implements" keyword! Just define the method.
func (c CreditCard) Process(amount float64) error {
    // implementation
    return nil
}
// CreditCard automatically satisfies PaymentProcessor!
```

**Key Insight**: In Go, if a type has all the methods an interface requires, it automatically implements that interface. No explicit declaration needed!

---

## 🏋️ Your Task

Create a **notification system** with the following requirements:

### 1. Create an interface `Notifier`
```go
type Notifier interface {
    Send(message string) error
    GetType() string
}
```

### 2. Implement 3 notification types:
- `EmailNotifier` - with fields: `to`, `from` (both strings)
- `SMSNotifier` - with fields: `phoneNumber` (string)  
- `SlackNotifier` - with fields: `channel`, `webhookURL` (both strings)

### 3. Create a function `NotifyAll`
```go
func NotifyAll(notifiers []Notifier, message string) error
```
This should send the message through ALL notifiers and return the first error encountered (or nil if all succeed).

### 4. In main(), demonstrate:
- Create instances of all 3 notifier types
- Put them in a slice of `Notifier`
- Call `NotifyAll` to send "System Alert: Server is down!"

---

## 📁 Files to Create

Create a single file: `main.go`

---

## ✅ Expected Output (example)
```
Sending Email to user@example.com from alerts@company.com: System Alert: Server is down!
Sending SMS to +1234567890: System Alert: Server is down!
Sending Slack message to #alerts: System Alert: Server is down!
All notifications sent successfully!
```

---

## 💡 Hints

1. For now, just `fmt.Println` instead of actually sending notifications
2. Use pointer receivers `func (e *EmailNotifier)` for methods that might modify state
3. Return `nil` for error if everything is successful
4. Package should be `package main`

---

## 🎓 What You'll Learn

1. **Implicit interface implementation** - The Go way
2. **Interface slices** - Polymorphism without inheritance
3. **Error handling pattern** - Returning errors vs throwing exceptions
4. **Method receivers** - Value vs pointer receivers

---

## ⏱️ Estimated Time: 15-20 minutes

When done, let me know and I'll review your code! 🚀


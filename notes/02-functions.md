# Functions, Methods & Interfaces - Quick Reference

## 📌 Function Basics

```go
// Basic function
func greet(name string) string {
    return "Hello, " + name
}

// Multiple parameters of same type
func add(a, b int) int {
    return a + b
}

// Multiple return values (unique to Go!)
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// Usage of multiple returns
result, err := divide(10, 2)
if err != nil {
    log.Fatal(err)
}

// Named return values
func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    return  // naked return
}
```

### Comparison with Java/C++

| Feature | Java | C++ | Go |
|---------|------|-----|-----|
| Multiple returns | ❌ (use objects) | ❌ (use tuple/struct) | ✅ Native |
| Named returns | ❌ | ❌ | ✅ |
| Default params | ❌ | ✅ | ❌ |
| Method overloading | ✅ | ✅ | ❌ |

**Tradeoff:** Go has no method overloading. You must use different function names or variadic functions.

---

## 🔄 Variadic Functions

```go
// Accept any number of arguments
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

// Usage
sum(1, 2, 3)       // 6
sum(1, 2, 3, 4, 5) // 15

// Pass slice as variadic
nums := []int{1, 2, 3}
sum(nums...)  // Spread operator
```

**Java equivalent:** `public int sum(int... nums)`  
**C++ equivalent:** Variadic templates or `std::initializer_list`

---

## 📦 Structs (Go's "Classes")

```go
// Define a struct (like a class without methods)
type Person struct {
    Name string   // Public (capitalized)
    Age  int
    city string   // private (lowercase)
}

// Create instances
p1 := Person{Name: "Viraj", Age: 25, city: "Mumbai"}
p2 := Person{"Viraj", 25, "Mumbai"}  // Order matters
p3 := Person{}                        // Zero values

// Pointer to struct
p4 := &Person{Name: "Viraj"}

// Access fields (same syntax for pointer and value!)
p1.Name     // "Viraj"
p4.Name     // "Viraj" (no -> like C++!)
```

### Java/C++ Comparison

```java
// Java
public class Person {
    private String name;
    private int age;
    
    public Person(String name, int age) {
        this.name = name;
        this.age = age;
    }
}
```

```go
// Go - Much simpler!
type Person struct {
    Name string
    Age  int
}
```

---

## 🔧 Methods (Functions on Structs)

```go
type Rectangle struct {
    Width  float64
    Height float64
}

// Value receiver (copy of struct)
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// Pointer receiver (can modify struct)
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor   // Modifies original
    r.Height *= factor
}

// Usage
rect := Rectangle{Width: 10, Height: 5}
area := rect.Area()     // 50

rect.Scale(2)           // rect is now 20 x 10
```

### When to use pointer receiver?

| Use Pointer Receiver | Use Value Receiver |
|---------------------|-------------------|
| Need to modify struct | Read-only operations |
| Struct is large | Struct is small |
| Consistency (if any method needs pointer) | Immutability desired |

**Java equivalent:** All methods can modify `this`  
**C++ equivalent:** `const` methods vs non-const methods

---

## 🎭 Interfaces (Duck Typing!)

**This is where Go shines differently from Java!**

```go
// Define interface
type Speaker interface {
    Speak() string
}

// Structs implement interfaces IMPLICITLY
type Dog struct {
    Name string
}

func (d Dog) Speak() string {
    return "Woof!"
}

type Cat struct {
    Name string
}

func (c Cat) Speak() string {
    return "Meow!"
}

// Both Dog and Cat implement Speaker automatically!
// No "implements" keyword needed

func MakeSound(s Speaker) {
    fmt.Println(s.Speak())
}

// Usage
dog := Dog{Name: "Buddy"}
cat := Cat{Name: "Whiskers"}

MakeSound(dog)  // "Woof!"
MakeSound(cat)  // "Meow!"
```

### Java vs Go Interfaces

```java
// Java - Explicit implementation
public class Dog implements Speaker {
    @Override
    public String speak() {
        return "Woof!";
    }
}
```

```go
// Go - Implicit implementation
type Dog struct{}

func (d Dog) Speak() string {
    return "Woof!"
}
// Dog now implements Speaker, no declaration needed!
```

**Tradeoff:**
- ✅ More flexible - types can implement interfaces retroactively
- ✅ No import cycles for interfaces
- ❌ Less explicit - harder to see what interfaces a type implements
- ❌ Easy to accidentally implement (or break) interfaces

---

## 🏗️ Composition over Inheritance

**Go has NO inheritance. Use embedding instead!**

```go
// Base "class" (struct)
type Animal struct {
    Name string
}

func (a Animal) Eat() {
    fmt.Println(a.Name, "is eating")
}

// "Child" class - EMBEDDING, not inheritance
type Dog struct {
    Animal  // Embedded struct
    Breed string
}

func (d Dog) Bark() {
    fmt.Println("Woof!")
}

// Usage
dog := Dog{
    Animal: Animal{Name: "Buddy"},
    Breed:  "Labrador",
}

dog.Eat()   // Promoted method from Animal
dog.Bark()  // Dog's own method
dog.Name    // Promoted field from Animal
```

### Java Inheritance vs Go Composition

```java
// Java
class Dog extends Animal {
    private String breed;
}
```

```go
// Go
type Dog struct {
    Animal       // Embedded
    Breed string
}
```

**Key Differences:**
- No `super()` calls
- No method overriding (but you can shadow)
- Multiple embedding allowed (like multiple inheritance, but safer)

---

## 📦 Empty Interface (like Object in Java)

```go
// interface{} accepts ANY type
func PrintAnything(v interface{}) {
    fmt.Println(v)
}

// Go 1.18+ uses 'any' alias
func PrintAnything(v any) {
    fmt.Println(v)
}

// Type assertion
func Process(v interface{}) {
    // Check if v is a string
    if s, ok := v.(string); ok {
        fmt.Println("String:", s)
    }
    
    // Type switch
    switch val := v.(type) {
    case int:
        fmt.Println("Integer:", val)
    case string:
        fmt.Println("String:", val)
    default:
        fmt.Println("Unknown type")
    }
}
```

**Java equivalent:** `Object` type  
**C++ equivalent:** `std::any` or templates

---

## 🎯 Quick Tips

1. **Start with value receivers** - only use pointers when needed
2. **Small interfaces** - Go idiom: 1-2 methods per interface
3. **Accept interfaces, return structs** - common pattern
4. **No constructors** - use factory functions like `NewPerson()`

```go
// Factory function pattern
func NewPerson(name string, age int) *Person {
    return &Person{
        Name: name,
        Age:  age,
    }
}
```

---

*Next: [03-concurrency.md](./03-concurrency.md)*



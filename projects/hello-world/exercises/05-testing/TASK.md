# Exercise 05: Testing in Go

## 🎯 Goal
Master Go's built-in testing - **table-driven tests**, **mocking**, and **test coverage**!

---

## 📚 Key Concept: Java vs Go Testing

### Java (JUnit + Mockito)
```java
@Test
public void testAdd() {
    Calculator calc = new Calculator();
    assertEquals(5, calc.add(2, 3));
}

@Mock
private UserRepository userRepo;

@Test
public void testGetUser() {
    when(userRepo.findById(1)).thenReturn(new User("John"));
    // ...
}
```

### Go (Built-in testing package)
```go
func TestAdd(t *testing.T) {
    calc := Calculator{}
    result := calc.Add(2, 3)
    if result != 5 {
        t.Errorf("Expected 5, got %d", result)
    }
}

// Mocking = just use interfaces!
type MockUserRepo struct{}
func (m *MockUserRepo) FindByID(id int) *User {
    return &User{Name: "John"}
}
```

---

## 📖 Go Testing Basics

### File Naming
```
calculator.go       # Source file
calculator_test.go  # Test file (must end with _test.go)
```

### Test Function Naming
```go
func TestXxx(t *testing.T)     # Test (starts with Test)
func BenchmarkXxx(b *testing.B) # Benchmark
func ExampleXxx()               # Example (shown in docs)
```

### Running Tests
```bash
go test              # Run tests in current package
go test ./...        # Run all tests recursively
go test -v           # Verbose output
go test -cover       # Show coverage %
go test -run TestAdd # Run specific test
```

---

## 🏋️ Your Task: Test a Calculator Service

### 1. Create Calculator with these methods

```go
// calculator/calculator.go
type Calculator struct{}

func (c *Calculator) Add(a, b int) int
func (c *Calculator) Subtract(a, b int) int
func (c *Calculator) Multiply(a, b int) int
func (c *Calculator) Divide(a, b int) (int, error)  // Error if b == 0
```

### 2. Write Table-Driven Tests

```go
// calculator/calculator_test.go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive numbers", 2, 3, 5},
        {"negative numbers", -2, -3, -5},
        {"zero", 0, 0, 0},
        {"mixed", -5, 10, 5},
    }

    calc := &Calculator{}
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := calc.Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d", 
                    tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

### 3. Test Error Cases

```go
func TestDivide(t *testing.T) {
    calc := &Calculator{}
    
    // Test division by zero
    _, err := calc.Divide(10, 0)
    if err == nil {
        t.Error("Expected error for division by zero")
    }
    
    // Test normal division
    result, err := calc.Divide(10, 2)
    if err != nil {
        t.Errorf("Unexpected error: %v", err)
    }
    if result != 5 {
        t.Errorf("Expected 5, got %d", result)
    }
}
```

### 4. Create a Service with Dependency (for mocking)

```go
// service/math_service.go
type MathRepository interface {
    GetOperands() (int, int, error)
}

type MathService struct {
    repo MathRepository
}

func NewMathService(repo MathRepository) *MathService {
    return &MathService{repo: repo}
}

func (s *MathService) AddFromRepo() (int, error) {
    a, b, err := s.repo.GetOperands()
    if err != nil {
        return 0, err
    }
    return a + b, nil
}
```

### 5. Mock the Repository in Tests

```go
// service/math_service_test.go
type MockMathRepo struct {
    a, b int
    err  error
}

func (m *MockMathRepo) GetOperands() (int, int, error) {
    return m.a, m.b, m.err
}

func TestAddFromRepo(t *testing.T) {
    // Test success case
    mockRepo := &MockMathRepo{a: 5, b: 3, err: nil}
    service := NewMathService(mockRepo)
    
    result, err := service.AddFromRepo()
    if err != nil {
        t.Errorf("Unexpected error: %v", err)
    }
    if result != 8 {
        t.Errorf("Expected 8, got %d", result)
    }
    
    // Test error case
    mockRepo = &MockMathRepo{err: errors.New("db error")}
    service = NewMathService(mockRepo)
    
    _, err = service.AddFromRepo()
    if err == nil {
        t.Error("Expected error, got nil")
    }
}
```

---

## 📁 Files to Create

```
05-testing/
├── calculator/
│   ├── calculator.go
│   └── calculator_test.go
└── service/
    ├── math_service.go
    └── math_service_test.go
```

---

## ✅ Expected Output

```bash
$ go test ./... -v

=== RUN   TestAdd
=== RUN   TestAdd/positive_numbers
=== RUN   TestAdd/negative_numbers
=== RUN   TestAdd/zero
=== RUN   TestAdd/mixed
--- PASS: TestAdd (0.00s)
    --- PASS: TestAdd/positive_numbers (0.00s)
    --- PASS: TestAdd/negative_numbers (0.00s)
    --- PASS: TestAdd/zero (0.00s)
    --- PASS: TestAdd/mixed (0.00s)
=== RUN   TestDivide
--- PASS: TestDivide (0.00s)
=== RUN   TestAddFromRepo
--- PASS: TestAddFromRepo (0.00s)
PASS
coverage: 100.0% of statements
ok      .../calculator    0.005s
ok      .../service       0.003s
```

---

## 💡 Hints

### 1. Table-Driven Tests Pattern
```go
tests := []struct {
    name     string  // Descriptive name
    input    int     // Input values
    expected int     // Expected output
    wantErr  bool    // Expect error?
}{
    {"case 1", 10, 20, false},
    {"error case", -1, 0, true},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // Test logic here
    })
}
```

### 2. Common Test Assertions
```go
// Check equality
if got != want {
    t.Errorf("got %v, want %v", got, want)
}

// Check error exists
if err == nil {
    t.Error("expected error, got nil")
}

// Check no error
if err != nil {
    t.Fatalf("unexpected error: %v", err)
}

// Fail immediately (t.Fatal vs t.Error)
t.Error("continues running other tests")
t.Fatal("stops this test immediately")
```

### 3. Test Setup/Teardown
```go
func TestMain(m *testing.M) {
    // Setup before all tests
    setup()
    
    // Run tests
    code := m.Run()
    
    // Teardown after all tests
    teardown()
    
    os.Exit(code)
}
```

### 4. Subtests for Organization
```go
func TestCalculator(t *testing.T) {
    t.Run("Addition", func(t *testing.T) {
        // Add tests
    })
    
    t.Run("Division", func(t *testing.T) {
        t.Run("normal", func(t *testing.T) {
            // Normal division
        })
        t.Run("by zero", func(t *testing.T) {
            // Division by zero
        })
    })
}
```

---

## 🎓 What You'll Learn

1. **Test file naming** - `*_test.go` convention
2. **Table-driven tests** - Go's idiomatic test pattern
3. **t.Run()** - Subtests for organization
4. **Mocking with interfaces** - No framework needed!
5. **Test coverage** - `go test -cover`
6. **t.Error vs t.Fatal** - When to stop vs continue

---

## 🆚 Java vs Go Testing Comparison

| Aspect | Java | Go |
|--------|------|-----|
| Framework | JUnit, TestNG | Built-in `testing` |
| Mocking | Mockito, EasyMock | Just use interfaces! |
| Assertions | AssertJ, Hamcrest | Manual `if` + `t.Error` |
| Annotations | `@Test`, `@Before` | Function naming convention |
| Coverage | JaCoCo | Built-in `go test -cover` |
| Running | Maven/Gradle | `go test` |

---

## 🚀 Bonus Challenge

Add benchmarks:

```go
func BenchmarkAdd(b *testing.B) {
    calc := &Calculator{}
    for i := 0; i < b.N; i++ {
        calc.Add(100, 200)
    }
}
```

Run with: `go test -bench=.`

---

## ⏱️ Estimated Time: 20-25 minutes

This is simpler than context! Go build it! 🧪



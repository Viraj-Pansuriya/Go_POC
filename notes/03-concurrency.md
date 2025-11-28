# Concurrency in Go - Quick Reference

## 🚀 Why Go Concurrency is Different

| Feature | Java | C++ | Go |
|---------|------|-----|-----|
| Unit | `Thread` (OS thread) | `std::thread` (OS thread) | `goroutine` (lightweight) |
| Memory | ~1MB per thread | ~1MB per thread | ~2KB per goroutine |
| Creation cost | Expensive | Expensive | Dirt cheap |
| Communication | Shared memory + locks | Shared memory + locks | Channels (CSP) |
| Max practical count | ~thousands | ~thousands | ~millions |

**Go's Philosophy:** "Don't communicate by sharing memory; share memory by communicating."

---

## 🏃 Goroutines

A goroutine is a lightweight thread managed by Go runtime.

```go
// Start a goroutine
go doSomething()

// With anonymous function
go func() {
    fmt.Println("Running in background")
}()

// With parameters (capture by value!)
go func(msg string) {
    fmt.Println(msg)
}("Hello")
```

### Java Comparison

```java
// Java - verbose
new Thread(() -> {
    System.out.println("Running in background");
}).start();

// Or with ExecutorService
executor.submit(() -> doSomething());
```

```go
// Go - simple
go doSomething()
```

---

## 📡 Channels (The Heart of Go Concurrency)

Channels are typed conduits for sending and receiving values.

```go
// Create a channel
ch := make(chan int)       // Unbuffered channel
ch := make(chan int, 10)   // Buffered channel (capacity 10)

// Send to channel
ch <- 42

// Receive from channel
value := <-ch

// Receive and check if closed
value, ok := <-ch
if !ok {
    // Channel is closed
}
```

### Basic Pattern: Worker

```go
func worker(id int, jobs <-chan int, results chan<- int) {
    for job := range jobs {
        result := job * 2  // Do work
        results <- result
    }
}

func main() {
    jobs := make(chan int, 100)
    results := make(chan int, 100)
    
    // Start 3 workers
    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }
    
    // Send 5 jobs
    for j := 1; j <= 5; j++ {
        jobs <- j
    }
    close(jobs)
    
    // Collect results
    for r := 1; r <= 5; r++ {
        fmt.Println(<-results)
    }
}
```

### Channel Directions (Type Safety)

```go
// Send-only channel
func sender(ch chan<- int) {
    ch <- 42
    // <-ch  // Compile error!
}

// Receive-only channel
func receiver(ch <-chan int) {
    value := <-ch
    // ch <- 42  // Compile error!
}
```

---

## 🔄 Select Statement (Multiplexing)

`select` lets you wait on multiple channel operations.

```go
select {
case msg := <-ch1:
    fmt.Println("Received from ch1:", msg)
case msg := <-ch2:
    fmt.Println("Received from ch2:", msg)
case ch3 <- 42:
    fmt.Println("Sent to ch3")
default:
    fmt.Println("No channel ready")
}
```

### Timeout Pattern

```go
select {
case result := <-ch:
    fmt.Println("Got result:", result)
case <-time.After(3 * time.Second):
    fmt.Println("Timeout!")
}
```

### Non-blocking Channel Operations

```go
select {
case msg := <-ch:
    fmt.Println("Received:", msg)
default:
    fmt.Println("No message available")
}
```

---

## 🔐 sync Package (When You Need Locks)

### WaitGroup (Wait for goroutines to finish)

```go
var wg sync.WaitGroup

for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        fmt.Println("Worker", id)
    }(i)
}

wg.Wait()  // Block until all done
fmt.Println("All workers finished")
```

**Java equivalent:** `CountDownLatch` or `CompletableFuture.allOf()`

### Mutex (Mutual Exclusion)

```go
type SafeCounter struct {
    mu    sync.Mutex
    count int
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}

func (c *SafeCounter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.count
}
```

**Java equivalent:** `synchronized` block or `ReentrantLock`

### RWMutex (Read-Write Lock)

```go
type SafeCache struct {
    mu    sync.RWMutex
    cache map[string]string
}

func (c *SafeCache) Get(key string) string {
    c.mu.RLock()         // Multiple readers allowed
    defer c.mu.RUnlock()
    return c.cache[key]
}

func (c *SafeCache) Set(key, value string) {
    c.mu.Lock()          // Exclusive write
    defer c.mu.Unlock()
    c.cache[key] = value
}
```

### Once (Run exactly once)

```go
var once sync.Once
var config *Config

func GetConfig() *Config {
    once.Do(func() {
        config = loadConfig()  // Only runs once
    })
    return config
}
```

**Java equivalent:** Double-checked locking or `Suppliers.memoize()`

---

## 🎯 Common Patterns

### Fan-Out, Fan-In

```go
func fanOut(input <-chan int, workers int) []<-chan int {
    channels := make([]<-chan int, workers)
    for i := 0; i < workers; i++ {
        channels[i] = worker(input)
    }
    return channels
}

func fanIn(channels ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    
    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for v := range c {
                out <- v
            }
        }(ch)
    }
    
    go func() {
        wg.Wait()
        close(out)
    }()
    
    return out
}
```

### Pipeline Pattern

```go
func generator(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

// Usage
nums := generator(1, 2, 3, 4)
squares := square(nums)
for s := range squares {
    fmt.Println(s)  // 1, 4, 9, 16
}
```

---

## ⚠️ Common Gotchas

### 1. Goroutine Leak

```go
// BAD - goroutine leaks if nobody reads
func bad() chan int {
    ch := make(chan int)
    go func() {
        ch <- 42  // Blocks forever if not read
    }()
    return ch
}

// GOOD - use buffered channel or ensure reader
func good() chan int {
    ch := make(chan int, 1)  // Buffered
    go func() {
        ch <- 42
    }()
    return ch
}
```

### 2. Loop Variable Capture

```go
// BAD - all goroutines see same value
for i := 0; i < 5; i++ {
    go func() {
        fmt.Println(i)  // Might print 5,5,5,5,5
    }()
}

// GOOD - pass as parameter
for i := 0; i < 5; i++ {
    go func(n int) {
        fmt.Println(n)  // Prints 0,1,2,3,4 (order varies)
    }(i)
}

// GOOD (Go 1.22+) - loop variable is per-iteration
for i := 0; i < 5; i++ {
    go func() {
        fmt.Println(i)  // Works correctly in Go 1.22+
    }()
}
```

### 3. Race Condition

```go
// BAD - race condition
counter := 0
for i := 0; i < 1000; i++ {
    go func() {
        counter++  // DATA RACE!
    }()
}

// GOOD - use atomic
var counter int64
for i := 0; i < 1000; i++ {
    go func() {
        atomic.AddInt64(&counter, 1)
    }()
}
```

---

## 🔍 Race Detector

Go has a built-in race detector!

```bash
# Run with race detection
go run -race main.go

# Test with race detection
go test -race ./...

# Build with race detection
go build -race
```

---

## 📊 Context (Cancellation & Timeout)

```go
import "context"

func doWork(ctx context.Context) error {
    select {
    case <-time.After(5 * time.Second):
        return nil  // Work completed
    case <-ctx.Done():
        return ctx.Err()  // Cancelled or timeout
    }
}

func main() {
    // With timeout
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    
    if err := doWork(ctx); err != nil {
        fmt.Println("Error:", err)
    }
}
```

---

*Next: [04-error-handling.md](./04-error-handling.md)*



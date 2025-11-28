# Exercise 03: Goroutines & Channels

## 🎯 Goal
Master Go's concurrency model - **goroutines** (lightweight threads) and **channels** (communication)!

---

## 📚 Key Concept: Java Threads vs Go Goroutines

### Java (Heavy Threads)
```java
// Each thread = ~1MB stack, OS-managed
Thread thread = new Thread(() -> {
    System.out.println("Running in thread");
});
thread.start();
thread.join();  // Wait for completion
```

### Go (Lightweight Goroutines)
```go
// Each goroutine = ~2KB stack, Go runtime managed
go func() {
    fmt.Println("Running in goroutine")
}()
// Can run millions of goroutines!
```

### Java Communication (Shared Memory)
```java
// Shared state + locks
synchronized(lock) {
    sharedData.add(item);
}
```

### Go Communication (Channels)
```go
// "Don't communicate by sharing memory; share memory by communicating"
ch := make(chan string)
go func() { ch <- "hello" }()  // Send
msg := <-ch                     // Receive
```

---

## 🏋️ Your Task: Build a Parallel URL Fetcher

Create a program that fetches multiple URLs **concurrently** and collects results.

### Requirements:

### 1. Create a Result struct
```go
type FetchResult struct {
    URL        string
    StatusCode int
    Duration   time.Duration
    Error      error
}
```

### 2. Create a function to "fetch" a URL (simulate with sleep)
```go
func fetchURL(url string) FetchResult {
    // Simulate network delay (random 100-500ms)
    // Return FetchResult with fake status code
}
```

### 3. Create a SEQUENTIAL fetcher
```go
func fetchSequential(urls []string) []FetchResult {
    // Fetch one by one, measure total time
}
```

### 4. Create a CONCURRENT fetcher using goroutines + channels
```go
func fetchConcurrent(urls []string) []FetchResult {
    // Spawn goroutine for each URL
    // Collect results via channel
    // Measure total time
}
```

### 5. In main(), compare both:
- Fetch 5 URLs sequentially
- Fetch 5 URLs concurrently
- Print results and time comparison

---

## 📁 Files to Create

```
03-goroutines-channels/
└── main.go
```

---

## ✅ Expected Output (example)
```
=== Sequential Fetching ===
Fetching: https://google.com ... 234ms (Status: 200)
Fetching: https://github.com ... 456ms (Status: 200)
Fetching: https://stackoverflow.com ... 123ms (Status: 200)
Fetching: https://reddit.com ... 345ms (Status: 200)
Fetching: https://twitter.com ... 278ms (Status: 200)
Sequential total time: 1436ms

=== Concurrent Fetching ===
Starting all fetches...
Received: https://stackoverflow.com - 123ms (Status: 200)
Received: https://google.com - 234ms (Status: 200)
Received: https://twitter.com - 278ms (Status: 200)
Received: https://reddit.com - 345ms (Status: 200)
Received: https://github.com - 456ms (Status: 200)
Concurrent total time: 456ms  ← Only as slow as slowest request!

⚡ Speedup: 3.15x faster with concurrency!
```

---

## 💡 Hints

### 1. Creating a goroutine
```go
go func() {
    // This runs concurrently
}()

// Or with a named function
go myFunction()
```

### 2. Creating and using channels
```go
// Create channel (can hold FetchResult values)
results := make(chan FetchResult)

// Send to channel (in goroutine)
go func() {
    results <- FetchResult{URL: "google.com", StatusCode: 200}
}()

// Receive from channel (blocks until data available)
result := <-results
```

### 3. Receiving multiple results
```go
// If you know how many results to expect
for i := 0; i < len(urls); i++ {
    result := <-results  // Blocks until a result arrives
    fmt.Println(result)
}
```

### 4. Simulating network delay
```go
import (
    "math/rand"
    "time"
)

func init() {
    rand.Seed(time.Now().UnixNano())  // Seed random (Go 1.20+ auto-seeds)
}

func fetchURL(url string) FetchResult {
    start := time.Now()
    
    // Random delay 100-500ms
    delay := time.Duration(100+rand.Intn(400)) * time.Millisecond
    time.Sleep(delay)
    
    return FetchResult{
        URL:        url,
        StatusCode: 200,
        Duration:   time.Since(start),
        Error:      nil,
    }
}
```

### 5. Measuring time
```go
start := time.Now()
// ... do work ...
elapsed := time.Since(start)
fmt.Printf("Took: %v\n", elapsed)
```

---

## 🎓 What You'll Learn

1. **Goroutines** - Lightweight concurrent execution
2. **Channels** - Safe communication between goroutines
3. **Blocking receives** - `<-ch` waits for data
4. **Fan-out pattern** - Spawn multiple goroutines for parallel work
5. **Performance gains** - Why concurrency matters

---

## ⚠️ Common Mistakes

```go
// ❌ BAD: Goroutine captures loop variable by reference
for _, url := range urls {
    go func() {
        fetchURL(url)  // Bug! All goroutines might see same url
    }()
}

// ✅ GOOD: Pass as parameter
for _, url := range urls {
    go func(u string) {
        fetchURL(u)
    }(url)
}

// ❌ BAD: Forgetting to receive all results (goroutine leak)
for _, url := range urls {
    go func(u string) {
        results <- fetchURL(u)
    }(url)
}
// Forgot to read from results channel!

// ✅ GOOD: Always receive what you send
for i := 0; i < len(urls); i++ {
    result := <-results
    fmt.Println(result)
}
```

---

## 🚀 Bonus Challenge (Optional)

Add a **timeout**: If any fetch takes more than 300ms, cancel it and return an error.

Hint: Use `select` with `time.After`:
```go
select {
case result := <-results:
    // Got result
case <-time.After(300 * time.Millisecond):
    // Timeout!
}
```

---

## ⏱️ Estimated Time: 25-30 minutes

When done, let me know! This is where Go really shines! 🌟


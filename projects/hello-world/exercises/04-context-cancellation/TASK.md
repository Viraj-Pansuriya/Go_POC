# Exercise 04: Context & Cancellation

## 🎯 Goal
Master Go's `context` package - essential for **timeouts**, **cancellation**, and **request-scoped data** in production apps!

---

## 📚 Key Concept: Why Context?

### The Problem
```go
// What if this takes 30 seconds? User already left!
func fetchData() {
    resp := http.Get("https://slow-api.com")  // Blocks forever?
}
```

### Java Approach
```java
// Future with timeout
Future<Response> future = executor.submit(() -> fetchData());
try {
    Response resp = future.get(5, TimeUnit.SECONDS);
} catch (TimeoutException e) {
    future.cancel(true);
}
```

### Go Approach (Context)
```go
// Context carries deadline, cancellation signal, and values
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result := fetchData(ctx)  // Function respects context deadline
```

---

## 📖 Context Types

| Function | Purpose | Example |
|----------|---------|---------|
| `context.Background()` | Root context, never cancels | Starting point |
| `context.TODO()` | Placeholder when unsure | Temporary use |
| `context.WithCancel(parent)` | Manual cancellation | User clicks "Stop" |
| `context.WithTimeout(parent, duration)` | Auto-cancel after time | API timeout |
| `context.WithDeadline(parent, time)` | Cancel at specific time | "Complete by 5 PM" |
| `context.WithValue(parent, key, val)` | Pass request-scoped data | User ID, trace ID |

---

## 🏋️ Your Task: Build a Timeout-Aware URL Fetcher

Extend your Exercise 03 fetcher to support **timeouts** and **cancellation**.

### Requirements:

### 1. Update FetchURL to accept context
```go
func FetchURL(ctx context.Context, url string, ch chan Result) {
    // Check if context is cancelled before/during work
    select {
    case <-ctx.Done():
        // Context cancelled or timed out
        ch <- Result{URL: url, Error: ctx.Err()}
        return
    default:
        // Continue with fetch
    }
    
    // ... do work, but periodically check ctx.Done()
}
```

### 2. Create timeout-aware Fetch method
```go
func (ur *UrlResponseFetcher) FetchWithTimeout(urls []string, timeout time.Duration) map[string]Result {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()  // Always call cancel to release resources!
    
    // ... launch goroutines with ctx
}
```

### 3. Handle partial results
- Some URLs may complete before timeout
- Some may timeout
- Collect whatever results you can!

### 4. In main(), demonstrate:
- Fetch with 200ms timeout (some should fail)
- Fetch with 2s timeout (all should succeed)
- Manual cancellation scenario

---

## 📁 Files to Create/Update

```
04-context-cancellation/
├── main.go
├── model/
│   └── result.go
└── service/
    └── fetcher.go
```

---

## ✅ Expected Output (example)
```
=== Fetch with 200ms timeout ===
✓ http://fast-api.com - 50ms (Status: 200)
✓ http://medium-api.com - 150ms (Status: 200)
✗ http://slow-api.com - TIMEOUT (context deadline exceeded)
✗ http://slower-api.com - TIMEOUT (context deadline exceeded)
Completed: 2/4, Timed out: 2/4
Total time: 200ms (capped by timeout)

=== Fetch with 2s timeout ===
✓ http://fast-api.com - 50ms (Status: 200)
✓ http://medium-api.com - 150ms (Status: 200)
✓ http://slow-api.com - 450ms (Status: 200)
✓ http://slower-api.com - 800ms (Status: 200)
Completed: 4/4, Timed out: 0/4
Total time: 800ms

=== Manual Cancellation ===
Started fetching...
Cancelled after 100ms!
Results before cancellation: 1/4
```

---

## 💡 Hints

### 1. The select statement (Go's switch for channels)
```go
select {
case result := <-ch:
    // Received a result
    fmt.Println("Got:", result)
case <-ctx.Done():
    // Context cancelled or timed out
    fmt.Println("Cancelled:", ctx.Err())
case <-time.After(1 * time.Second):
    // Timeout for this specific operation
    fmt.Println("Waited 1 second, giving up")
}
```

### 2. Simulating slow requests
```go
func FetchURL(ctx context.Context, url string, ch chan Result) {
    delay := time.Duration(rand.Intn(500)) * time.Millisecond
    
    select {
    case <-time.After(delay):
        // Work completed
        ch <- Result{URL: url, StatusCode: 200, Duration: delay}
    case <-ctx.Done():
        // Cancelled before completion
        ch <- Result{URL: url, Error: ctx.Err()}
    }
}
```

### 3. Collecting results with timeout
```go
func FetchWithTimeout(urls []string, timeout time.Duration) []Result {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    
    ch := make(chan Result, len(urls))  // Buffered channel!
    
    for _, url := range urls {
        go FetchURL(ctx, url, ch)
    }
    
    var results []Result
    for i := 0; i < len(urls); i++ {
        select {
        case result := <-ch:
            results = append(results, result)
        case <-ctx.Done():
            // Timeout reached, stop collecting
            return results
        }
    }
    return results
}
```

### 4. Always defer cancel()!
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()  // IMPORTANT: Releases resources even if timeout not reached
```

---

## 🎓 What You'll Learn

1. **context.Context** - The Go way to handle cancellation
2. **select statement** - Multiplexing on channels
3. **Timeouts** - Don't wait forever
4. **Graceful degradation** - Return partial results
5. **Resource cleanup** - Always `defer cancel()`

---

## 🔑 Key Pattern: Context Propagation

In real apps, context flows through the entire call chain:

```go
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()  // Get context from HTTP request
    
    user, err := userService.GetUser(ctx, userID)      // Pass context down
    orders, err := orderService.GetOrders(ctx, userID) // Pass context down
    
    // If client disconnects, ctx is cancelled automatically!
}

func (s *UserService) GetUser(ctx context.Context, id int) (*User, error) {
    // Check if cancelled
    if ctx.Err() != nil {
        return nil, ctx.Err()
    }
    
    return s.db.QueryContext(ctx, "SELECT * FROM users WHERE id = ?", id)
}
```

---

## ⚠️ Common Mistakes

```go
// ❌ BAD: Not checking context
func FetchURL(ctx context.Context, url string) Result {
    time.Sleep(5 * time.Second)  // Ignores cancellation!
    return Result{}
}

// ✅ GOOD: Respects context
func FetchURL(ctx context.Context, url string) Result {
    select {
    case <-time.After(5 * time.Second):
        return Result{StatusCode: 200}
    case <-ctx.Done():
        return Result{Error: ctx.Err()}
    }
}

// ❌ BAD: Forgetting to cancel
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
// Missing: defer cancel()  <- Resource leak!

// ❌ BAD: Using context.Background() everywhere
func inner() {
    ctx := context.Background()  // Wrong! Should receive from caller
}

// ✅ GOOD: Pass context through
func inner(ctx context.Context) {
    // Use the passed context
}
```

---

## ⏱️ Estimated Time: 30-35 minutes

This is a crucial concept for production Go! When done, let me know! 🚀


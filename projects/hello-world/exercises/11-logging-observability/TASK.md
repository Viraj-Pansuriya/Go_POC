# Exercise 11: Logging & Observability with Zap

## 🎯 Goal
Learn structured logging in Go using Zap - the Go equivalent of SLF4J/Logback with JSON output for production!

---

## 📚 Java (SLF4J/Logback) vs Go (Zap)

### Java with SLF4J
```java
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@Service
public class UserService {
    private static final Logger log = LoggerFactory.getLogger(UserService.class);
    
    public User getUser(Long id) {
        log.info("Fetching user with id={}", id);
        
        try {
            User user = repository.findById(id);
            log.debug("Found user: name={}, email={}", user.getName(), user.getEmail());
            return user;
        } catch (Exception e) {
            log.error("Failed to fetch user id={}", id, e);
            throw e;
        }
    }
}

// logback.xml for JSON output
<appender name="JSON" class="ch.qos.logback.core.ConsoleAppender">
    <encoder class="net.logstash.logback.encoder.LogstashEncoder"/>
</appender>
```

### Go with Zap
```go
import "go.uber.org/zap"

type UserService struct {
    logger *zap.Logger
    repo   UserRepository
}

func (s *UserService) GetUser(id uint) (*User, error) {
    s.logger.Info("fetching user", zap.Uint("user_id", id))
    
    user, err := s.repo.FindByID(id)
    if err != nil {
        s.logger.Error("failed to fetch user",
            zap.Uint("user_id", id),
            zap.Error(err),
        )
        return nil, err
    }
    
    s.logger.Debug("found user",
        zap.String("name", user.Name),
        zap.String("email", user.Email),
    )
    return user, nil
}

// Output (JSON - production ready for ELK/Splunk):
// {"level":"info","ts":1699123456.789,"caller":"service/user.go:15","msg":"fetching user","user_id":123}
```

---

## 🏗️ Why Zap?

| Feature | Zap | Standard log | Logrus |
|---------|-----|--------------|--------|
| **Performance** | ⚡ Fastest | 🐢 Slow | 🐌 Slower |
| **Structured** | ✅ Native | ❌ No | ✅ Yes |
| **Zero alloc** | ✅ Yes | ❌ No | ❌ No |
| **JSON output** | ✅ Built-in | ❌ No | ✅ Yes |
| **Log levels** | ✅ All | ⚠️ Limited | ✅ All |

**Zap is used by**: Uber, Kubernetes, CockroachDB, etcd

---

## 🚀 Your Tasks

### Task 1: Basic Zap Setup

**Install Zap:**
```bash
go get -u go.uber.org/zap
```

Create logger configuration:

```go
// logger/logger.go
package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// InitLogger creates a logger based on environment
func InitLogger(isDev bool) error {
    var err error
    
    if isDev {
        // Development: human-readable, colored output
        config := zap.NewDevelopmentConfig()
        config.EncoderConfig.EncodeLevel = zapcore.CapitalColoredLevelEncoder
        Log, err = config.Build()
    } else {
        // Production: JSON format for log aggregation
        config := zap.NewProductionConfig()
        config.EncoderConfig.TimeKey = "timestamp"
        config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
        Log, err = config.Build()
    }
    
    return err
}

// Sync flushes any buffered log entries
func Sync() {
    if Log != nil {
        _ = Log.Sync()
    }
}
```

### Task 2: Structured Logging Patterns

```go
// Different field types for type safety
logger.Info("user action",
    zap.String("action", "login"),
    zap.Int("user_id", 123),
    zap.Duration("latency", time.Since(start)),
    zap.Time("timestamp", time.Now()),
    zap.Bool("success", true),
    zap.Any("metadata", map[string]string{"ip": "1.2.3.4"}),
)

// Error with stack trace
logger.Error("database error",
    zap.Error(err),  // Automatically extracts error message
    zap.Stack("stacktrace"),
)

// Conditional logging
if logger.Core().Enabled(zap.DebugLevel) {
    // Only compute expensive debug info if debug is enabled
    logger.Debug("detailed info", zap.Any("data", expensiveComputation()))
}
```

### Task 3: Request Logger Middleware

Create Gin middleware for HTTP request logging:

```go
// middleware/request_logger.go
package middleware

import (
    "time"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

func RequestLogger(logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        query := c.Request.URL.RawQuery
        
        // Process request
        c.Next()
        
        // Log after request completes
        latency := time.Since(start)
        status := c.Writer.Status()
        
        // Choose log level based on status code
        logFunc := logger.Info
        if status >= 500 {
            logFunc = logger.Error
        } else if status >= 400 {
            logFunc = logger.Warn
        }
        
        logFunc("http request",
            zap.String("method", c.Request.Method),
            zap.String("path", path),
            zap.String("query", query),
            zap.Int("status", status),
            zap.Duration("latency", latency),
            zap.String("client_ip", c.ClientIP()),
            zap.String("user_agent", c.Request.UserAgent()),
            zap.Int("body_size", c.Writer.Size()),
        )
    }
}
```

### Task 4: Context-Aware Logging

Pass request context through handlers:

```go
// middleware/request_id.go
func RequestID(logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        
        // Create child logger with request context
        reqLogger := logger.With(
            zap.String("request_id", requestID),
            zap.String("path", c.Request.URL.Path),
        )
        
        // Store in context for handlers to use
        c.Set("logger", reqLogger)
        c.Header("X-Request-ID", requestID)
        
        c.Next()
    }
}

// handler/user_handler.go
func (h *UserHandler) GetUser(c *gin.Context) {
    // Get request-scoped logger
    logger := c.MustGet("logger").(*zap.Logger)
    
    logger.Info("fetching user")  // Automatically includes request_id!
    // ...
}
```

### Task 5: Log Levels & Configuration

```go
// logger/logger.go
type LogConfig struct {
    Level       string `mapstructure:"level"`       // debug, info, warn, error
    Format      string `mapstructure:"format"`      // json, console
    OutputPath  string `mapstructure:"output_path"` // stdout, /var/log/app.log
    Development bool   `mapstructure:"development"`
}

func InitLoggerWithConfig(cfg LogConfig) (*zap.Logger, error) {
    // Parse log level
    level, err := zapcore.ParseLevel(cfg.Level)
    if err != nil {
        level = zapcore.InfoLevel
    }
    
    // Build encoder config
    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.TimeKey = "timestamp"
    encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
    
    // Choose encoder (JSON for prod, Console for dev)
    var encoder zapcore.Encoder
    if cfg.Format == "console" || cfg.Development {
        encoderConfig.EncodeLevel = zapcore.CapitalColoredLevelEncoder
        encoder = zapcore.NewConsoleEncoder(encoderConfig)
    } else {
        encoder = zapcore.NewJSONEncoder(encoderConfig)
    }
    
    // Build core
    core := zapcore.NewCore(
        encoder,
        zapcore.AddSync(os.Stdout),  // Could add file output
        level,
    )
    
    // Build logger with options
    opts := []zap.Option{
        zap.AddCaller(),  // Add file:line info
        zap.AddStacktrace(zapcore.ErrorLevel),  // Stack trace on errors
    }
    
    if cfg.Development {
        opts = append(opts, zap.Development())
    }
    
    return zap.New(core, opts...), nil
}
```

### Task 6: Integration with Gin

```go
// main.go
func main() {
    // Load config
    cfg, _ := config.LoadConfig("./configs", "dev")
    
    // Initialize logger
    logger, err := logger.InitLoggerWithConfig(cfg.Log)
    if err != nil {
        log.Fatal("failed to init logger", err)
    }
    defer logger.Sync()
    
    // Replace Gin's default logger
    gin.SetMode(gin.ReleaseMode)
    r := gin.New()  // New() instead of Default() - no default middleware
    
    // Add our custom middleware
    r.Use(gin.Recovery())  // Panic recovery
    r.Use(middleware.RequestID(logger))
    r.Use(middleware.RequestLogger(logger))
    
    // Create handlers with logger
    userHandler := handler.NewUserHandler(logger, userService)
    
    r.GET("/users/:id", userHandler.GetUser)
    
    logger.Info("starting server",
        zap.String("host", cfg.Server.Host),
        zap.Int("port", cfg.Server.Port),
    )
    
    r.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}
```

---

## 📁 Project Structure

```
11-logging-observability/
├── main.go
├── go.mod
├── configs/
│   └── config.yaml
├── config/
│   └── config.go
├── logger/
│   └── logger.go          # Logger initialization
├── middleware/
│   ├── request_logger.go  # HTTP request logging
│   └── request_id.go      # Request ID tracking
├── handler/
│   └── user_handler.go    # Example handler
└── TASK.md
```

---

## 📄 Config File

```yaml
# configs/config.yaml
server:
  port: 8080

log:
  level: debug        # debug, info, warn, error
  format: console     # json, console
  development: true   # Enables dev-friendly output
```

---

## 💡 Key Patterns

### 1. Sugar Logger (Less Verbose)
```go
// Regular logger (type-safe, fastest)
logger.Info("user logged in",
    zap.String("email", email),
    zap.Int("user_id", id),
)

// Sugar logger (printf-style, slightly slower)
sugar := logger.Sugar()
sugar.Infof("user %s logged in", email)
sugar.Infow("user logged in", "email", email, "user_id", id)
```

### 2. Child Loggers with Context
```go
// Create child logger with persistent fields
userLogger := logger.With(
    zap.String("service", "user-service"),
    zap.String("version", "1.0.0"),
)

// All logs from userLogger include these fields automatically
userLogger.Info("processing request")
// Output: {"service":"user-service","version":"1.0.0","msg":"processing request"}
```

### 3. Sampling (High-Volume Logs)
```go
// Sample logs to reduce volume in production
config := zap.NewProductionConfig()
config.Sampling = &zap.SamplingConfig{
    Initial:    100,  // First 100 logs/sec
    Thereafter: 10,   // Then log every 10th
}
```

### 4. Multiple Outputs
```go
// Log to both stdout and file
core := zapcore.NewTee(
    zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level),
    zapcore.NewCore(encoder, zapcore.AddSync(logFile), level),
)
```

---

## 🆚 Java/SLF4J Comparison

| SLF4J/Logback | Zap | Notes |
|---------------|-----|-------|
| `LoggerFactory.getLogger()` | `zap.NewProduction()` | Create logger |
| `log.info("msg {}", val)` | `logger.Info("msg", zap.Any("key", val))` | Structured |
| `MDC.put("requestId", id)` | `logger.With(zap.String(...))` | Context |
| `@Slf4j` (Lombok) | Inject via constructor | No magic |
| `logback.xml` | `zap.Config{}` | Configuration |
| JSON encoder | Built-in | Production format |

---

## ✅ Expected Output

### Development Mode (Console)
```
2024-01-15T10:30:45.123+0530    INFO    middleware/request_logger.go:25    http request    {"method": "GET", "path": "/users/123", "status": 200, "latency": "2.5ms"}
2024-01-15T10:30:45.125+0530    DEBUG   handler/user_handler.go:18    fetching user    {"request_id": "abc-123", "user_id": 123}
```

### Production Mode (JSON)
```json
{"level":"info","timestamp":"2024-01-15T10:30:45.123Z","caller":"middleware/request_logger.go:25","msg":"http request","method":"GET","path":"/users/123","status":200,"latency":"2.5ms"}
{"level":"debug","timestamp":"2024-01-15T10:30:45.125Z","caller":"handler/user_handler.go:18","msg":"fetching user","request_id":"abc-123","user_id":123}
```

---

## 🎓 What You'll Learn

1. **Structured logging** - Key-value pairs instead of string concatenation
2. **Log levels** - Debug, Info, Warn, Error, Fatal
3. **Performance** - Zero-allocation logging
4. **JSON output** - Ready for ELK/Splunk/Datadog
5. **Request tracing** - Request ID propagation
6. **Context-aware logging** - Child loggers with fields
7. **Middleware pattern** - HTTP request/response logging

---

## 🔐 Best Practices

1. **Never log sensitive data** - Passwords, tokens, PII
2. **Use structured fields** - Not string interpolation
3. **Include request ID** - For distributed tracing
4. **Choose appropriate levels** - Debug for dev, Info for prod
5. **Sync before exit** - `defer logger.Sync()`
6. **Use child loggers** - Add context fields once

---

## ⏱️ Estimated Time: 25-30 minutes

Good logging saves hours of debugging! 📊


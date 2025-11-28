# Packages & Modules - Quick Reference

## 📦 Package Basics

Every Go file belongs to a package.

```go
// File: mypackage/utils.go
package mypackage

// Exported (public) - starts with uppercase
func PublicFunction() {}

// Unexported (private) - starts with lowercase
func privateFunction() {}

// Exported type
type User struct {
    Name  string  // Public field
    email string  // Private field
}
```

### Comparison

| Concept | Java | C++ | Go |
|---------|------|-----|-----|
| Visibility | `public`/`private`/`protected` | `public`/`private`/`protected` | Uppercase/lowercase |
| Package | `package com.example;` | `namespace` | `package example` |
| Import | `import com.example.*;` | `#include` | `import "example"` |

---

## 🗂️ Module System (go mod)

A **module** is a collection of packages with a `go.mod` file.

### Initialize a New Module

```bash
# Create a new module
go mod init github.com/username/myproject

# Creates go.mod file:
# module github.com/username/myproject
# go 1.21
```

### go.mod File

```go
module github.com/username/myproject

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    golang.org/x/sync v0.3.0
)

require (
    // Indirect dependencies (auto-managed)
    github.com/some/indirect v1.0.0 // indirect
)
```

**Java equivalent:** `pom.xml` (Maven) or `build.gradle` (Gradle)  
**C++ equivalent:** `CMakeLists.txt` + `conanfile.txt`

---

## 📥 Importing Packages

```go
// Single import
import "fmt"

// Multiple imports
import (
    "fmt"
    "os"
    "strings"
    
    // External packages
    "github.com/gin-gonic/gin"
    
    // Alias import
    mylog "github.com/sirupsen/logrus"
    
    // Dot import (use sparingly!)
    . "github.com/onsi/gomega"
    
    // Blank import (for side effects only)
    _ "github.com/lib/pq"
)
```

### Import Aliases

```go
import (
    "crypto/rand"
    mrand "math/rand"  // Alias to avoid conflict
)

func main() {
    rand.Read(buf)        // crypto/rand
    mrand.Intn(100)       // math/rand
}
```

---

## 🏗️ Project Structure

### Simple Project

```
myproject/
├── go.mod
├── go.sum
├── main.go
└── internal/
    └── config/
        └── config.go
```

### Larger Project (Recommended)

```
myproject/
├── go.mod
├── go.sum
├── cmd/
│   ├── server/
│   │   └── main.go      # go run ./cmd/server
│   └── cli/
│       └── main.go      # go run ./cmd/cli
├── internal/            # Private packages
│   ├── config/
│   ├── handler/
│   └── service/
├── pkg/                 # Public packages (can be imported)
│   └── utils/
├── api/                 # API definitions (OpenAPI, proto)
└── scripts/             # Build/deploy scripts
```

### Key Directories

| Directory | Purpose |
|-----------|---------|
| `cmd/` | Application entry points (`main` packages) |
| `internal/` | Private code (can't be imported by others) |
| `pkg/` | Public code (can be imported by others) |
| `api/` | API specs (protobuf, OpenAPI) |
| `web/` | Web assets |
| `scripts/` | Build, install, analysis scripts |

---

## 🔧 Go Workspace (Monorepo)

For multiple modules in one repo, use **go.work** (Go 1.18+).

```bash
# Initialize workspace
go work init

# Add modules
go work use ./projects/hello-world
go work use ./projects/api-server
```

### go.work File

```go
go 1.21

use (
    ./projects/hello-world
    ./projects/api-server
    ./pkg/shared
)
```

**Benefits:**
- Work on multiple modules simultaneously
- Local changes across modules without publishing
- One `go build` for everything

---

## 📋 Common go Commands

```bash
# Module management
go mod init <module-path>    # Initialize new module
go mod tidy                  # Add missing, remove unused deps
go mod download              # Download dependencies
go mod verify                # Verify dependencies
go mod graph                 # Print module dependency graph

# Get dependencies
go get github.com/pkg@latest # Add/update dependency
go get github.com/pkg@v1.2.3 # Specific version
go get -u ./...              # Update all dependencies

# Build and run
go build ./...               # Build all packages
go run .                     # Run current package
go run ./cmd/server          # Run specific package
go install ./...             # Build and install

# Testing
go test ./...                # Test all packages
go test -v ./...             # Verbose output
go test -cover ./...         # With coverage
go test -race ./...          # With race detection

# Tooling
go fmt ./...                 # Format code
go vet ./...                 # Static analysis
go doc fmt.Println           # View documentation
```

---

## 🔄 Dependency Versioning

Go uses **Semantic Import Versioning**.

```go
// v0.x.x or v1.x.x
import "github.com/pkg/errors"

// v2.x.x and above - path includes version!
import "github.com/pkg/errors/v2"
```

### Updating Dependencies

```bash
# Update to latest minor/patch
go get -u github.com/gin-gonic/gin

# Update to latest patch only
go get -u=patch github.com/gin-gonic/gin

# Update all
go get -u ./...

# Downgrade
go get github.com/gin-gonic/gin@v1.8.0

# Remove unused
go mod tidy
```

---

## 🏭 Internal Packages

The `internal` directory is special - packages inside cannot be imported from outside the module.

```
mymodule/
├── go.mod
├── internal/
│   └── secret/
│       └── secret.go    # Can only be imported within mymodule
└── pkg/
    └── public/
        └── public.go    # Can be imported by anyone
```

```go
// Inside mymodule - works
import "mymodule/internal/secret"

// Outside mymodule - compile error!
import "mymodule/internal/secret"  // ERROR
```

---

## 🌐 Standard Library Highlights

```go
// Common packages
import (
    "fmt"      // Formatted I/O
    "os"       // Operating system
    "io"       // I/O primitives
    "strings"  // String manipulation
    "strconv"  // String conversions
    "time"     // Time and duration
    "context"  // Context for cancellation
    "errors"   // Error creation
    "log"      // Logging
    "net/http" // HTTP client/server
    "encoding/json"  // JSON encoding
    "sync"     // Synchronization primitives
    "testing"  // Testing support
)
```

### Quick Examples

```go
// JSON
type User struct {
    Name  string `json:"name"`
    Email string `json:"email,omitempty"`
}
data, _ := json.Marshal(user)
json.Unmarshal(data, &user)

// HTTP
http.HandleFunc("/", handler)
http.ListenAndServe(":8080", nil)

// Time
time.Now()
time.Sleep(time.Second)
time.After(5 * time.Second)

// Context
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

---

## 🎯 Best Practices

1. **Keep packages focused** - one purpose per package

2. **Use internal/** - for implementation details

3. **Avoid package cycles** - Go doesn't allow them anyway

4. **Name packages well**:
   - Short, lowercase, single word if possible
   - `http`, not `httputil` or `httpHandler`
   - Package name shouldn't repeat in function names
   
   ```go
   // BAD
   package http
   func HTTPGet()  // http.HTTPGet() is redundant
   
   // GOOD
   package http
   func Get()      // http.Get() is clean
   ```

5. **Document exported items**:
   
   ```go
   // User represents a system user.
   type User struct {
       // Name is the user's display name.
       Name string
   }
   
   // NewUser creates a new User with the given name.
   func NewUser(name string) *User {
       return &User{Name: name}
   }
   ```

---

*Back to: [README.md](../README.md)*



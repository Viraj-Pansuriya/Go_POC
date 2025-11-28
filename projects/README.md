# Projects

This folder contains all your Go learning projects.

## 📁 Current Projects

| Project | Description | Status |
|---------|-------------|--------|
| [hello-world](./hello-world/) | Basic Go syntax and structure | ✅ Ready |

## 🆕 Creating a New Project

```bash
# 1. Create project directory
mkdir projects/my-new-project
cd projects/my-new-project

# 2. Initialize Go module
go mod init github.com/viraj/go-mono-repo/projects/my-new-project

# 3. Create main.go
touch main.go

# 4. Add to workspace (from repo root)
cd ../..
# Edit go.work to add: ./projects/my-new-project
```

Or use this one-liner from repo root:

```bash
PROJECT_NAME=my-project && \
mkdir -p projects/$PROJECT_NAME && \
cd projects/$PROJECT_NAME && \
go mod init github.com/viraj/go-mono-repo/projects/$PROJECT_NAME && \
echo 'package main\n\nimport "fmt"\n\nfunc main() {\n\tfmt.Println("Hello from '$PROJECT_NAME'")\n}' > main.go
```

Then add to `go.work`:
```go
use (
    ./projects/hello-world
    ./projects/my-project  // Add this line
)
```

## 📝 Project Ideas for Learning

### Beginner
- [ ] CLI calculator
- [ ] Todo list (in-memory)
- [ ] File reader/writer
- [ ] HTTP client (fetch API data)

### Intermediate
- [ ] REST API server
- [ ] Concurrent file downloader
- [ ] Chat server (TCP/WebSocket)
- [ ] JSON config parser

### Advanced
- [ ] Rate limiter
- [ ] Worker pool
- [ ] Caching layer
- [ ] gRPC service



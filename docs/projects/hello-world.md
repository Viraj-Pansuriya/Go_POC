# Hello World Project

## Overview

A simple "Hello, World!" program - the traditional first program in any language.

## Learning Goals

- [x] Understand Go file structure
- [x] Learn about `package main`
- [x] Use `func main()` as entry point
- [x] Import and use `fmt` package
- [x] Run a Go program

## How to Run

```bash
cd projects/hello-world
go run .

# Or build and run
go build -o bin/hello
./bin/hello
```

## Key Code Snippets

### Basic Structure

```go
package main  // Every executable needs this

import "fmt"  // Import standard library package

func main() {  // Entry point - no arguments, no return
    fmt.Println("Hello, World!")
}
```

### Comparison with Java

```java
// Java - more verbose
public class HelloWorld {
    public static void main(String[] args) {
        System.out.println("Hello, World!");
    }
}
```

### Comparison with C++

```cpp
// C++ - more ceremony
#include <iostream>

int main() {
    std::cout << "Hello, World!" << std::endl;
    return 0;
}
```

## Notes & Observations

1. **No class wrapping** - Unlike Java, no need to wrap everything in a class
2. **No return value** - Unlike C++, `main()` doesn't return an int
3. **Automatic imports** - Use `goimports` tool to auto-add imports
4. **No semicolons** - Go lexer inserts them automatically
5. **Braces required** - Opening brace must be on same line as `func`

## What's Next?

Try modifying the program to:
- Accept command line arguments using `os.Args`
- Use `fmt.Printf` for formatted output
- Add a function that returns a greeting

## Resources

- [Go Tour - Hello World](https://go.dev/tour/welcome/1)
- [Effective Go](https://go.dev/doc/effective_go)



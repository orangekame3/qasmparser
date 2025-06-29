# qasmparser - OpenQASM 3.0 Parser for Go

A robust, reusable Go package for parsing OpenQASM 3.0 quantum programs, built on ANTLR4. This package provides a clean API for building tools like formatters, linters, and other OpenQASM code analysis tools.

## Features

- **Complete OpenQASM 3.0 Support**: Based on official OpenQASM grammar
- **Clean AST**: Well-structured Abstract Syntax Tree with visitor pattern
- **Error Handling**: Comprehensive error reporting with position information
- **Flexible API**: Parse from strings, files, or readers
- **Performance**: Efficient parsing for large QASM files
- **Extensible**: Visitor pattern for custom AST traversal

## Installation

### As a Go Library

```bash
go get github.com/orangekame3/qasmparser
```

### CLI Tool

#### Homebrew (macOS/Linux)

```bash
brew install orangekame3/tap/qasmparser
```

#### Binary Releases

Download pre-built binaries from the [releases page](https://github.com/orangekame3/qasmparser/releases).

#### Docker

```bash
docker run --rm ghcr.io/orangekame3/qasmparser:latest --version
```

#### From Source

```bash
git clone https://github.com/orangekame3/qasmparser.git
cd qasmparser
task build-cli
```

## Quick Start

### Basic Parsing

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/orangekame3/qasmparser/parser"
)

func main() {
    qasm := `
    OPENQASM 3.0;
    include "stdgates.qasm";
    
    qubit[2] q;
    h q[0];
    cx q[0], q[1];
    `
    
    p := parser.NewParser()
    program, err := p.ParseString(qasm)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Parsed %d statements\n", len(program.Statements))
}
```

### Error Handling

```go
p := parser.NewParser()
result := p.ParseWithErrors(qasmCode)

if result.HasErrors() {
    for _, err := range result.Errors {
        fmt.Printf("Error at line %d: %s\n", err.Position.Line, err.Message)
    }
}

// Continue with partial AST if available
if result.Program != nil {
    // Process partial program
}
```

### AST Visitor Pattern

```go
type PrintVisitor struct {
    parser.BaseVisitor
}

func (v *PrintVisitor) VisitGateCall(node *parser.GateCall) interface{} {
    fmt.Printf("Gate: %s on %d qubits\n", node.Name, len(node.Qubits))
    return nil
}

// Usage
visitor := &PrintVisitor{}
parser.Walk(visitor, program)
```

## Project Structure

```bash
qasmparser/
├── parser/          # Core parser package
│   ├── ast.go      # AST node definitions
│   ├── parser.go   # Main parser interface
│   ├── visitor.go  # Visitor pattern implementation
│   └── errors.go   # Error handling
├── cmd/qasmparser/  # CLI tool
│   └── main.go     # CLI implementation
├── gen/parser/      # Generated ANTLR code
├── grammar/         # ANTLR grammar files
├── testdata/        # Test QASM files
├── examples/        # Usage examples
├── bin/             # Built binaries
└── Taskfile.yaml    # Build automation
```

## Build Setup

This package uses [Task](https://taskfile.dev) for build automation and requires ANTLR4 for code generation.

### Prerequisites

1. Install Task:

```bash
# macOS
brew install go-task/tap/go-task

# Other platforms: https://taskfile.dev/installation/
```

2. Install ANTLR4:

```bash
# macOS
brew install antlr

# Or use the automated installer
task install-antlr
```

### Development Setup

```bash
# Setup development environment
task setup

# Generate ANTLR parser files
task generate

# Run tests
task test

# Run tests with coverage
task test-coverage

# Build the package
task build

# Build CLI tool
task build-cli

# Run linters
task lint

# Run examples
task examples
```

## Supported OpenQASM 3.0 Features

### ✅ Fully Supported

- Version declarations (`OPENQASM 3.0;`)
- Include statements (`include "stdgates.qasm";`)
- Qubit declarations (`qubit q;`, `qubit[n] q;`)
- Classical declarations (`bit c;`, `int[32] i;`)
- Gate calls (`h q;`, `cx control, target;`)
- Parameterized gates (`rz(theta) q;`)
- Measurement (`measure q -> c;`)
- Basic expressions and arithmetic
- Comments (line and block)

### 🚧 Partial Support

- Gate definitions (basic syntax)
- Control flow (`if`, `for`, `while`)
- Function definitions (`def`)

### 📋 Planned

- Advanced type system
- Pulse-level programming
- Timing constructs
- Complex expressions

## API Reference

### Parser

```go
// Create parser
parser := parser.NewParser()
parser := parser.NewParserWithOptions(&parser.ParseOptions{
    StrictMode:      true,
    IncludeComments: false,
    ErrorRecovery:   true,
    MaxErrors:       10,
})

// Parse methods
program, err := parser.ParseString(content)
program, err := parser.ParseFile(filename)
program, err := parser.ParseReader(reader)
program, err := parser.ParseWithContext(ctx, content)

// Error handling
result := parser.ParseWithErrors(content)
err := parser.Validate(content)
```

### AST Nodes

Key AST node types:

- `Program` - Root node
- `Version` - OpenQASM version declaration
- `QuantumDeclaration` - Qubit declarations
- `ClassicalDeclaration` - Classical variable declarations
- `GateCall` - Gate applications
- `Measurement` - Measure statements
- `Include` - Include statements
- `Expression` types for all expressions

### Visitor Pattern

```go
type Visitor interface {
    VisitProgram(node *Program) interface{}
    VisitGateCall(node *GateCall) interface{}
    // ... other visit methods
}

// Base visitor with default implementations
type BaseVisitor struct{}

// Walk AST with visitor
result := parser.Walk(visitor, node)

// Depth-first traversal
depthFirst := parser.NewDepthFirstVisitor(visitor)
parser.Walk(depthFirst, program)
```

## CLI Tool

The package includes a command-line tool for parsing and analyzing QASM files:

```bash
# Build the CLI tool
task build-cli

# Parse QASM files
./bin/qasmparser parse file.qasm

# Validate syntax only
./bin/qasmparser validate file.qasm

# Show AST structure
./bin/qasmparser ast --format tree file.qasm

# Get program statistics
./bin/qasmparser stats --format json file.qasm

# Format QASM code (demonstration)
./bin/qasmparser format file.qasm

# Global options
./bin/qasmparser --help
```

### CLI Usage Examples

```bash
# Parse multiple files with verbose output
./bin/qasmparser parse -v file1.qasm file2.qasm

# Show AST in JSON format
./bin/qasmparser ast -f json --depth 3 circuit.qasm

# Get statistics for all QASM files in directory
./bin/qasmparser stats testdata/*.qasm

# Validate with strict mode
./bin/qasmparser validate --strict --max-errors 5 *.qasm

# Save output to file
./bin/qasmparser parse -o results.txt testdata/test_simple.qasm
```

## Examples

See the `examples/` directory for complete examples:

- `parse_simple/` - Basic parsing example
- `ast_visitor/` - AST visitor pattern usage
- `error_handling/` - Error handling patterns

## Integration

### With qasmfmt

```go
import "github.com/orangekame3/qasmparser/parser"

type Formatter struct {
    parser *parser.Parser
}

func (f *Formatter) Format(content string) (string, error) {
    program, err := f.parser.ParseString(content)
    if err != nil {
        return "", err
    }
    return f.formatProgram(program), nil
}
```

### Future qasmlint Integration

```go
type Linter struct {
    parser *parser.Parser
    rules  []Rule
}

func (l *Linter) Lint(content string) ([]Issue, error) {
    program, err := l.parser.ParseString(content)
    if err != nil {
        return nil, err
    }

    var issues []Issue
    for _, rule := range l.rules {
        issues = append(issues, rule.Check(program)...)
    }
    return issues, nil
}
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes and add tests
4. Run `task validate` to ensure all checks pass
5. Submit a pull request

## Development

```bash
# Clean and regenerate everything
task clean generate

# Watch for grammar changes (requires fswatch)
task dev

# Copy test data from qasmfmt
task copy-testdata

# Download official grammar files
task download-grammar
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Related Projects

- [qasmfmt](https://github.com/orangekame3/qasmfmt) - OpenQASM 3.0 formatter (uses this parser)
- [OpenQASM](https://github.com/openqasm/openqasm) - Official OpenQASM specification

## Acknowledgments

- Built with [ANTLR4](https://www.antlr.org/)
- Grammar based on [official OpenQASM 3.0 specification](https://openqasm.com/)
- Inspired by modern parser design patterns

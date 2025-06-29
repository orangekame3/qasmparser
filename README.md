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

```bash
go get github.com/orangekame3/qasmparser
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
type GateCountVisitor struct {
    parser.BaseVisitor
    count int
}

func (v *GateCountVisitor) VisitGateCall(node *parser.GateCall) interface{} {
    v.count++
    fmt.Printf("Gate: %s on %d qubits\n", node.Name, len(node.Qubits))
    return nil
}

// Usage
visitor := &GateCountVisitor{}
parser.Walk(visitor, program)
fmt.Printf("Total gates: %d\n", visitor.count)
```

## Project Structure

```bash
qasmparser/
├── parser/          # Core parser package
│   ├── ast.go      # AST node definitions
│   ├── parser.go   # Main parser interface
│   ├── visitor.go  # Visitor pattern implementation
│   └── errors.go   # Error handling
├── gen/parser/      # Generated ANTLR code
├── grammar/         # ANTLR grammar files
├── testdata/        # Test QASM files
├── examples/        # Usage examples
└── Taskfile.yaml    # Build automation
```

## Development Setup

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

### Development Commands

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

### Parser Creation

```go
// Create parser with default options
parser := parser.NewParser()

// Create parser with custom options
parser := parser.NewParserWithOptions(&parser.ParseOptions{
    StrictMode:      true,
    IncludeComments: false,
    ErrorRecovery:   true,
    MaxErrors:       10,
})
```

### Parse Methods

```go
// Parse from string
program, err := parser.ParseString(content)

// Parse from file
program, err := parser.ParseFile(filename)

// Parse with detailed error information
result := parser.ParseWithErrors(content)

// Quick validation (returns first error only)
err := parser.Validate(content)
```

### AST Node Types

Key AST node types:

- `Program` - Root node containing all statements
- `Version` - OpenQASM version declaration
- `QuantumDeclaration` - Qubit declarations (`qubit q;`)
- `ClassicalDeclaration` - Classical variable declarations (`bit c;`)
- `GateCall` - Gate applications (`h q;`)
- `Measurement` - Measure statements (`measure q -> c;`)
- `Include` - Include statements (`include "file.qasm";`)
- Various `Expression` types for literals, identifiers, and operations

### Visitor Pattern

```go
// Create custom visitor
type MyVisitor struct {
    parser.BaseVisitor
    // Add custom fields
}

func (v *MyVisitor) VisitGateCall(node *parser.GateCall) interface{} {
    // Custom logic for gate calls
    return nil
}

// Walk AST
visitor := &MyVisitor{}
result := parser.Walk(visitor, program)

// Use depth-first visitor for automatic child traversal
depthFirst := parser.NewDepthFirstVisitor(visitor)
parser.Walk(depthFirst, program)
```

## Examples

See the `examples/` directory for complete working examples:

- [`parse_simple/`](examples/parse_simple/) - Basic parsing and validation
- [`ast_visitor/`](examples/ast_visitor/) - AST visitor pattern usage with statistics
- [`error_handling/`](examples/error_handling/) - Comprehensive error handling patterns

To run examples:

```bash
# Run all examples
task examples

# Run specific example
go run examples/parse_simple/main.go
```

## Integration Examples

### With qasmfmt

This parser was originally developed for [qasmfmt](https://github.com/orangekame3/qasmfmt):

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

### Custom Analysis Tool

```go
type QuantumAnalyzer struct {
    parser *parser.Parser
}

func (a *QuantumAnalyzer) Analyze(content string) (*Report, error) {
    result := a.parser.ParseWithErrors(content)
    
    report := &Report{
        Errors: result.Errors,
    }
    
    if result.Program != nil {
        visitor := &AnalysisVisitor{report: report}
        parser.Walk(visitor, result.Program)
    }
    
    return report, nil
}
```

### Linter Integration

```go
type QASMLinter struct {
    parser *parser.Parser
    rules  []Rule
}

func (l *QASMLinter) Lint(content string) ([]Issue, error) {
    program, err := l.parser.ParseString(content)
    if err != nil {
        return nil, err
    }

    var issues []Issue
    for _, rule := range l.rules {
        visitor := &RuleVisitor{rule: rule}
        parser.Walk(visitor, program)
        issues = append(issues, visitor.Issues()...)
    }
    return issues, nil
}
```

## Testing

```bash
# Run tests
task test

# Run tests with coverage
task test-coverage

# Run benchmarks
task bench
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes and add tests
4. Run `task validate` to ensure all checks pass
5. Submit a pull request

### Development Workflow

```bash
# Clean and regenerate everything
task clean generate

# Watch for grammar changes (requires fswatch)
task dev

# Download official grammar files
task download-grammar

# Validate package before release
task validate
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
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial OpenQASM 3.0 parser implementation as a Go package
- Complete AST node definitions with visitor pattern
- Comprehensive error handling and reporting
- ANTLR4-based grammar for OpenQASM 3.0
- Three comprehensive examples showing package usage
- Task-based build automation
- Test suite with coverage reporting

### Features
- **Parser Library**: Reusable Go package for OpenQASM 3.0 parsing
- **AST Support**: Complete Abstract Syntax Tree with visitor pattern
- **Error Recovery**: Robust error handling with detailed position reporting
- **Flexible API**: Multiple parse methods (string, file, reader, with errors)
- **Extensible**: Visitor pattern for custom AST traversal and analysis
- **Package-focused**: Designed for integration into other tools

### Supported OpenQASM 3.0 Features

#### ✅ Fully Supported
- Version declarations (`OPENQASM 3.0;`)
- Include statements (`include "stdgates.qasm";`)
- Qubit declarations (`qubit q;`, `qubit[n] q;`)
- Classical declarations (`bit c;`, `int[32] i;`)
- Gate calls (`h q;`, `cx control, target;`)
- Parameterized gates (`rz(theta) q;`)
- Measurement (`measure q -> c;`)
- Basic expressions and arithmetic
- Comments (line and block)

#### 🚧 Partial Support
- Gate definitions (basic syntax)
- Control flow (`if`, `for`, `while`)
- Function definitions (`def`)

#### 📋 Planned
- Advanced type system
- Pulse-level programming
- Timing constructs
- Complex expressions

### Examples
- `parse_simple/` - Basic parsing and validation example
- `ast_visitor/` - AST visitor pattern usage with statistics
- `error_handling/` - Comprehensive error handling patterns

### API
- `NewParser()` - Create parser with default options
- `NewParserWithOptions()` - Create parser with custom configuration
- `ParseString()` - Parse QASM from string
- `ParseFile()` - Parse QASM from file
- `ParseWithErrors()` - Parse with detailed error collection
- `Validate()` - Quick validation (first error only)
- Visitor pattern with `Walk()` and `NewDepthFirstVisitor()`

### Notes
This package was separated from [qasmfmt](https://github.com/orangekame3/qasmfmt) to provide a reusable OpenQASM 3.0 parser for the Go ecosystem.
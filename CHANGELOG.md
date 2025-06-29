# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial OpenQASM 3.0 parser implementation
- Complete AST node definitions with visitor pattern
- CLI tool with multiple commands (parse, validate, ast, stats, format)
- Support for multiple output formats (text, JSON, tree)
- Comprehensive error handling and reporting
- ANTLR4-based grammar for OpenQASM 3.0
- Docker containerization support
- Homebrew formula for easy installation
- GitHub Actions CI/CD pipeline
- GoReleaser configuration for automated releases

### Features
- **Parser Library**: Reusable Go package for OpenQASM 3.0 parsing
- **CLI Tool**: Command-line interface with 5 subcommands
- **AST Support**: Complete Abstract Syntax Tree with visitor pattern
- **Error Recovery**: Robust error handling with detailed reporting
- **Multi-format Output**: Text, JSON, and tree visualization
- **Cross-platform**: Support for Linux, macOS, and Windows
- **Package Management**: Available via Homebrew, Docker, and binary releases

### Supported OpenQASM 3.0 Features
- Version declarations (`OPENQASM 3.0;`)
- Include statements (`include "stdgates.qasm";`)
- Qubit declarations (`qubit q;`, `qubit[n] q;`)
- Classical declarations (`bit c;`, `int[32] i;`)
- Gate calls (`h q;`, `cx control, target;`)
- Parameterized gates (`rz(theta) q;`)
- Measurement (`measure q -> c;`)
- Basic expressions and arithmetic
- Comments (line and block)
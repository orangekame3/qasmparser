# GitHub Copilot Code Review Instructions

## Review Focus Areas

Please pay special attention to the following areas during code review:

### Go Language Specific
- Proper use of Go Modules
- Error handling implementation
- Appropriate use of goroutines and channels
- Potential memory leaks
- Performance optimizations

### Security
- Input validation
- SQL injection prevention
- Proper handling of sensitive information
- Access control and authorization

### Code Quality
- Readability and maintainability
- Appropriate naming conventions
- Comment quality and documentation
- Test coverage
- Code duplication reduction

### QASM (Quantum Assembly Language) Parsing Specific
- QASM syntax parsing accuracy
- AST generation correctness
- Error reporting and recovery
- Parser performance optimization
- Memory management in parser operations

## Review Language
Please provide review comments in English.

## Excluded Items
The following files should be excluded from review:
- `go.sum`, `go.mod` (dependency files)
- `vendor/` directory
- Auto-generated files in `gen/` directory
- ANTLR generated parser files
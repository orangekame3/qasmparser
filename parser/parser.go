package parser

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

// ParseOptions configures the parser behavior
type ParseOptions struct {
	// StrictMode enables strict OpenQASM 3.0 compliance
	StrictMode bool

	// IncludeComments preserves comments in the AST
	IncludeComments bool

	// ErrorRecovery enables error recovery for partial parsing
	ErrorRecovery bool

	// MaxErrors limits the number of errors to collect
	MaxErrors int
}

// DefaultParseOptions returns default parsing options
func DefaultParseOptions() *ParseOptions {
	return &ParseOptions{
		StrictMode:      false,
		IncludeComments: true,
		ErrorRecovery:   true,
		MaxErrors:       100,
	}
}

// Parser represents the main OpenQASM 3.0 parser
type Parser struct {
	options *ParseOptions
}

// NewParser creates a new parser with default options
func NewParser() *Parser {
	return &Parser{
		options: DefaultParseOptions(),
	}
}

// NewParserWithOptions creates a parser with custom options
func NewParserWithOptions(opts *ParseOptions) *Parser {
	if opts == nil {
		opts = DefaultParseOptions()
	}
	return &Parser{
		options: opts,
	}
}

// ParseString parses QASM code from a string
func (p *Parser) ParseString(content string) (*Program, error) {
	result := p.ParseWithErrors(content)
	if result.HasErrors() {
		return result.Program, &result.Errors[0]
	}
	return result.Program, nil
}

// ParseReader parses QASM code from an io.Reader
func (p *Parser) ParseReader(reader io.Reader) (*Program, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return p.ParseString(string(content))
}

// ParseFile parses QASM code from a file
func (p *Parser) ParseFile(filename string) (*Program, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return p.ParseString(string(content))
}

// ParseWithContext parses with context for cancellation
func (p *Parser) ParseWithContext(ctx context.Context, content string) (*Program, error) {
	// Check if context is already cancelled
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// For now, we'll parse directly without context support
	// Future enhancement could implement context-aware parsing
	return p.ParseString(content)
}

// Validate validates QASM syntax without building full AST
func (p *Parser) Validate(content string) error {
	result := p.ParseWithErrors(content)
	if result.HasErrors() {
		return &result.Errors[0]
	}
	return nil
}

// ParseWithErrors returns partial results even with errors
func (p *Parser) ParseWithErrors(content string) *ParseResult {
	// Preprocess content to handle common issues
	content = p.preprocessContent(content)

	// Create input stream
	input := antlr.NewInputStream(content)

	// Create lexer (this will be replaced with generated code)
	lexer := p.createLexer(input)

	// Create error listener for lexer
	lexerErrors := NewErrorListener()
	lexer.RemoveErrorListeners()
	if !p.options.ErrorRecovery {
		lexer.AddErrorListener(lexerErrors)
	}

	// Create token stream
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create parser (this will be replaced with generated code)
	parser := p.createParser(stream)

	// Create error listener for parser
	parserErrors := NewErrorListener()
	parser.RemoveErrorListeners()
	if !p.options.ErrorRecovery {
		parser.AddErrorListener(parserErrors)
	}

	// Parse the program
	tree := p.parseProgram(parser)

	// Collect all errors
	allErrors := make([]ParseError, 0)
	allErrors = append(allErrors, lexerErrors.GetErrors()...)
	allErrors = append(allErrors, parserErrors.GetErrors()...)

	// Limit errors if specified
	if p.options.MaxErrors > 0 && len(allErrors) > p.options.MaxErrors {
		allErrors = allErrors[:p.options.MaxErrors]
	}

	// Convert parse tree to AST
	program := p.convertToAST(tree, content)

	return &ParseResult{
		Program: program,
		Errors:  allErrors,
	}
}

// preprocessContent handles common formatting issues
func (p *Parser) preprocessContent(content string) string {
	// Normalize line endings
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")

	// Ensure content ends with newline
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	return content
}

// createLexer creates the ANTLR lexer
// This is a placeholder - will be replaced with generated lexer
func (p *Parser) createLexer(input antlr.CharStream) antlr.Lexer {
	// This will be replaced with:
	// return parser.Newqasm3Lexer(input)
	panic("Generated lexer not available - run 'task generate' to create ANTLR files")
}

// createParser creates the ANTLR parser
// This is a placeholder - will be replaced with generated parser
func (p *Parser) createParser(stream antlr.TokenStream) antlr.Parser {
	// This will be replaced with:
	// return parser.Newqasm3Parser(stream)
	panic("Generated parser not available - run 'task generate' to create ANTLR files")
}

// parseProgram parses the root program
// This is a placeholder - will be replaced with generated parser
func (p *Parser) parseProgram(parser antlr.Parser) antlr.Tree {
	// This will be replaced with:
	// return parser.(*parser.qasm3Parser).Program()
	panic("Generated parser not available - run 'task generate' to create ANTLR files")
}

// convertToAST converts ANTLR parse tree to our AST
func (p *Parser) convertToAST(tree antlr.Tree, content string) *Program {
	// This will implement the conversion from ANTLR parse tree to our AST
	// For now, return a minimal program
	return &Program{
		BaseNode: BaseNode{
			Position: Position{Line: 1, Column: 1, Offset: 0},
			EndPos:   Position{Line: 1, Column: len(content), Offset: len(content)},
		},
		Statements: make([]Statement, 0),
		Comments:   make([]Comment, 0),
	}
}

// extractComments extracts comments from token stream
func (p *Parser) extractComments(stream antlr.TokenStream) []Comment {
	if !p.options.IncludeComments {
		return nil
	}

	comments := make([]Comment, 0)
	// This will be implemented to extract comments from the token stream
	return comments
}

// GetOptions returns the current parser options
func (p *Parser) GetOptions() *ParseOptions {
	return p.options
}

// SetOptions updates the parser options
func (p *Parser) SetOptions(opts *ParseOptions) {
	if opts != nil {
		p.options = opts
	}
}

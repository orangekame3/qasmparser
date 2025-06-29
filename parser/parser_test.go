package parser

import (
	"strings"
	"testing"
)

func TestNewParser(t *testing.T) {
	parser := NewParser()
	if parser == nil {
		t.Fatal("NewParser returned nil")
	}
	if parser.options == nil {
		t.Fatal("Parser options are nil")
	}
}

func TestNewParserWithOptions(t *testing.T) {
	opts := &ParseOptions{
		StrictMode:      true,
		IncludeComments: false,
		ErrorRecovery:   false,
		MaxErrors:       5,
	}
	
	parser := NewParserWithOptions(opts)
	if parser == nil {
		t.Fatal("NewParserWithOptions returned nil")
	}
	if parser.options != opts {
		t.Fatal("Parser options not set correctly")
	}
}

func TestDefaultParseOptions(t *testing.T) {
	opts := DefaultParseOptions()
	if opts == nil {
		t.Fatal("DefaultParseOptions returned nil")
	}
	
	if opts.StrictMode != false {
		t.Error("Expected StrictMode to be false")
	}
	if opts.IncludeComments != true {
		t.Error("Expected IncludeComments to be true")
	}
	if opts.ErrorRecovery != true {
		t.Error("Expected ErrorRecovery to be true")
	}
	if opts.MaxErrors != 100 {
		t.Error("Expected MaxErrors to be 100")
	}
}

func TestPreprocessContent(t *testing.T) {
	parser := NewParser()
	
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normalize CRLF",
			input:    "line1\r\nline2\r\n",
			expected: "line1\nline2\n",
		},
		{
			name:     "normalize CR",
			input:    "line1\rline2\r",
			expected: "line1\nline2\n",
		},
		{
			name:     "add final newline",
			input:    "line1\nline2",
			expected: "line1\nline2\n",
		},
		{
			name:     "preserve existing newline",
			input:    "line1\nline2\n",
			expected: "line1\nline2\n",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.preprocessContent(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestParseError(t *testing.T) {
	err := ParseError{
		Message:  "unexpected token",
		Position: Position{Line: 1, Column: 5},
		Type:     "syntax",
		Context:  "gate declaration",
	}
	
	expected := "syntax error at line 1, column 5: unexpected token (context: gate declaration)"
	if err.Error() != expected {
		t.Errorf("Expected %q, got %q", expected, err.Error())
	}
	
	// Test without context
	err.Context = ""
	expected = "syntax error at line 1, column 5: unexpected token"
	if err.Error() != expected {
		t.Errorf("Expected %q, got %q", expected, err.Error())
	}
}

func TestParseResult(t *testing.T) {
	result := &ParseResult{
		Program: &Program{},
		Errors: []ParseError{
			{Message: "error 1", Type: "syntax", Position: Position{Line: 1, Column: 1}},
			{Message: "error 2", Type: "semantic", Position: Position{Line: 2, Column: 1}},
		},
	}
	
	if !result.HasErrors() {
		t.Error("Expected HasErrors to return true")
	}
	
	messages := result.ErrorMessages()
	if len(messages) != 2 {
		t.Errorf("Expected 2 error messages, got %d", len(messages))
	}
	
	str := result.String()
	if !strings.Contains(str, "error 1") || !strings.Contains(str, "error 2") {
		t.Error("Expected string to contain both error messages")
	}
	
	// Test empty result
	emptyResult := &ParseResult{}
	if emptyResult.HasErrors() {
		t.Error("Expected HasErrors to return false for empty result")
	}
	if emptyResult.String() != "No errors" {
		t.Error("Expected 'No errors' for empty result")
	}
}

func TestErrorListener(t *testing.T) {
	listener := NewErrorListener()
	if listener == nil {
		t.Fatal("NewErrorListener returned nil")
	}
	
	if listener.HasErrors() {
		t.Error("New listener should have no errors")
	}
	
	// Simulate syntax error
	listener.SyntaxError(nil, nil, 1, 5, "unexpected token", nil)
	
	if !listener.HasErrors() {
		t.Error("Expected listener to have errors after SyntaxError call")
	}
	
	errors := listener.GetErrors()
	if len(errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(errors))
	}
	
	if errors[0].Type != "syntax" {
		t.Errorf("Expected error type 'syntax', got %q", errors[0].Type)
	}
	if errors[0].Position.Line != 1 {
		t.Errorf("Expected line 1, got %d", errors[0].Position.Line)
	}
	if errors[0].Position.Column != 5 {
		t.Errorf("Expected column 5, got %d", errors[0].Position.Column)
	}
}

func TestASTNodes(t *testing.T) {
	// Test Position
	pos := Position{Line: 1, Column: 5, Offset: 10}
	if pos.Line != 1 || pos.Column != 5 || pos.Offset != 10 {
		t.Error("Position fields not set correctly")
	}
	
	// Test BaseNode
	base := BaseNode{
		Position: Position{Line: 1, Column: 1},
		EndPos:   Position{Line: 1, Column: 10},
	}
	if base.Pos().Line != 1 {
		t.Error("BaseNode Pos() not working correctly")
	}
	if base.End().Column != 10 {
		t.Error("BaseNode End() not working correctly")
	}
	
	// Test Program
	program := &Program{
		BaseNode:   base,
		Statements: make([]Statement, 0),
		Comments:   make([]Comment, 0),
	}
	if program.String() != "Program" {
		t.Error("Program String() not working correctly")
	}
	
	// Test Version
	version := &Version{
		BaseNode: base,
		Number:   "3.0",
	}
	if version.String() != "Version: 3.0" {
		t.Error("Version String() not working correctly")
	}
	
	// Test QuantumDeclaration
	qDecl := &QuantumDeclaration{
		BaseNode:   base,
		Type:       "qubit",
		Identifier: "q",
	}
	if qDecl.String() != "QuantumDeclaration: q" {
		t.Error("QuantumDeclaration String() not working correctly")
	}
	
	// Test interface compliance
	var stmt Statement = qDecl
	stmt.StatementNode() // Should not panic
	
	var node Node = qDecl
	if node.Pos().Line != 1 {
		t.Error("QuantumDeclaration doesn't implement Node correctly")
	}
}

func TestVisitorPattern(t *testing.T) {
	// Test BaseVisitor
	visitor := &BaseVisitor{}
	
	program := &Program{
		BaseNode:   BaseNode{Position: Position{Line: 1, Column: 1}},
		Statements: make([]Statement, 0),
	}
	
	result := visitor.VisitProgram(program)
	if result != nil {
		t.Error("BaseVisitor should return nil")
	}
	
	// Test Walk function
	result = Walk(visitor, program)
	if result != nil {
		t.Error("Walk should return nil for BaseVisitor")
	}
	
	// Test Walk with nil node
	result = Walk(visitor, nil)
	if result != nil {
		t.Error("Walk should return nil for nil node")
	}
}

// Mock visitor for testing
type mockVisitor struct {
	BaseVisitor
	visitedNodes []string
}

func (m *mockVisitor) VisitProgram(node *Program) interface{} {
	m.visitedNodes = append(m.visitedNodes, "Program")
	return "visited program"
}

func (m *mockVisitor) VisitQuantumDeclaration(node *QuantumDeclaration) interface{} {
	m.visitedNodes = append(m.visitedNodes, "QuantumDeclaration")
	return "visited quantum declaration"
}

func TestMockVisitor(t *testing.T) {
	visitor := &mockVisitor{visitedNodes: make([]string, 0)}
	
	program := &Program{
		BaseNode:   BaseNode{Position: Position{Line: 1, Column: 1}},
		Statements: make([]Statement, 0),
	}
	
	result := Walk(visitor, program)
	if result != "visited program" {
		t.Errorf("Expected 'visited program', got %v", result)
	}
	
	if len(visitor.visitedNodes) != 1 || visitor.visitedNodes[0] != "Program" {
		t.Error("Visitor didn't visit Program node correctly")
	}
	
	qDecl := &QuantumDeclaration{
		BaseNode:   BaseNode{Position: Position{Line: 1, Column: 1}},
		Type:       "qubit",
		Identifier: "q",
	}
	
	result = Walk(visitor, qDecl)
	if result != "visited quantum declaration" {
		t.Errorf("Expected 'visited quantum declaration', got %v", result)
	}
	
	if len(visitor.visitedNodes) != 2 || visitor.visitedNodes[1] != "QuantumDeclaration" {
		t.Error("Visitor didn't visit QuantumDeclaration node correctly")
	}
}
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/fang"
	"github.com/orangekame3/qasmparser/parser"
)

var (
	// Version, Commit, and BuildDate are set via ldflags during build
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

type CLI struct {
	// Global flags
	Verbose    bool   `help:"Enable verbose output" short:"v"`
	Output     string `help:"Output file (default: stdout)" short:"o"`
	Format     string `help:"Output format: text, json, tree" short:"f" default:"text"`
	Strict     bool   `help:"Enable strict mode"`
	MaxErrors  int    `help:"Maximum number of errors to report" default:"100"`
	NoComments bool   `help:"Exclude comments from output"`
	NoRecovery bool   `help:"Disable error recovery"`
}

type ParseCmd struct {
	Files      []string `arg:"" help:"QASM files to parse" required:""`
	ErrorsOnly bool     `help:"Only show errors, suppress successful parsing messages" short:"e"`
}

type ValidateCmd struct {
	Files []string `arg:"" help:"QASM files to validate" required:""`
}

type ASTCmd struct {
	Files   []string `arg:"" help:"QASM files to show AST for" required:""`
	Depth   int      `help:"Maximum depth to display (-1 for unlimited)" short:"d" default:"-1"`
	Compact bool     `help:"Compact output format" short:"c"`
}

type StatsCmd struct {
	Files []string `arg:"" help:"QASM files to analyze" required:""`
}

type FormatCmd struct {
	Files []string `arg:"" help:"QASM files to format" required:""`
}

func main() {
	ctx := context.Background()
	cli := &CLI{}

	app := &fang.Program{
		Name:        "qasmparser",
		Description: "OpenQASM 3.0 parser and analyzer",
		Version:     Version,
		Commands: []*fang.Command{
			{
				Name:        "parse",
				Description: "Parse QASM files and report errors",
				Exec:        cli.runParse,
			},
			{
				Name:        "validate", 
				Description: "Validate QASM files (quick syntax check)",
				Exec:        cli.runValidate,
			},
			{
				Name:        "ast",
				Description: "Display AST structure",
				Exec:        cli.runAST,
			},
			{
				Name:        "stats",
				Description: "Show program statistics",
				Exec:        cli.runStats,
			},
			{
				Name:        "format",
				Description: "Format QASM files (demonstration)",
				Exec:        cli.runFormat,
			},
		},
	}

	if err := app.Run(ctx, cli); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func (c *CLI) createParser() *parser.Parser {
	opts := &parser.ParseOptions{
		StrictMode:      c.Strict,
		IncludeComments: !c.NoComments,
		ErrorRecovery:   !c.NoRecovery,
		MaxErrors:       c.MaxErrors,
	}
	return parser.NewParserWithOptions(opts)
}

func (c *CLI) writeOutput(content string) error {
	if c.Output == "" {
		fmt.Print(content)
		return nil
	}

	return os.WriteFile(c.Output, []byte(content), 0644)
}

func (c *CLI) runParse(ctx context.Context, cmd *ParseCmd) error {
	p := c.createParser()

	var allResults []ParseFileResult
	hasErrors := false

	for _, filename := range cmd.Files {
		result := c.parseFile(p, filename)
		allResults = append(allResults, result)

		if result.Error != nil || result.ParseResult.HasErrors() {
			hasErrors = true
		}

		if cmd.ErrorsOnly && !result.ParseResult.HasErrors() && result.Error == nil {
			continue
		}

		output := c.formatParseResult(filename, result)
		if err := c.writeOutput(output); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
	}

	if c.Format == "json" {
		jsonOutput, err := json.MarshalIndent(allResults, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		return c.writeOutput(string(jsonOutput) + "\n")
	}

	if hasErrors {
		os.Exit(1)
	}
	return nil
}

func (c *CLI) runValidate(ctx context.Context, cmd *ValidateCmd) error {
	p := c.createParser()
	hasErrors := false

	for _, filename := range cmd.Files {
		content, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: failed to read file: %v\n", filename, err)
			hasErrors = true
			continue
		}

		if err := p.Validate(string(content)); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", filename, err)
			hasErrors = true
		} else if c.Verbose {
			fmt.Printf("%s: ✓ valid\n", filename)
		}
	}

	if hasErrors {
		os.Exit(1)
	}
	return nil
}

func (c *CLI) runAST(ctx context.Context, cmd *ASTCmd) error {
	p := c.createParser()

	for _, filename := range cmd.Files {
		result := c.parseFile(p, filename)

		if result.Error != nil {
			fmt.Fprintf(os.Stderr, "%s: failed to read file: %v\n", filename, result.Error)
			continue
		}

		if result.ParseResult.HasErrors() && c.Verbose {
			fmt.Fprintf(os.Stderr, "%s: parsed with %d errors\n", filename, len(result.ParseResult.Errors))
		}

		if result.ParseResult.Program != nil {
			output := c.formatAST(filename, result.ParseResult.Program, cmd.Depth, cmd.Compact)
			if err := c.writeOutput(output); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}
		}
	}

	return nil
}

func (c *CLI) runStats(ctx context.Context, cmd *StatsCmd) error {
	p := c.createParser()

	for _, filename := range cmd.Files {
		result := c.parseFile(p, filename)

		if result.Error != nil {
			fmt.Fprintf(os.Stderr, "%s: failed to read file: %v\n", filename, result.Error)
			continue
		}

		if result.ParseResult.Program != nil {
			stats := c.analyzeProgram(result.ParseResult.Program)
			output := c.formatStats(filename, stats)
			if err := c.writeOutput(output); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}
		}
	}

	return nil
}

func (c *CLI) runFormat(ctx context.Context, cmd *FormatCmd) error {
	p := c.createParser()

	for _, filename := range cmd.Files {
		result := c.parseFile(p, filename)

		if result.Error != nil {
			fmt.Fprintf(os.Stderr, "%s: failed to read file: %v\n", filename, result.Error)
			continue
		}

		if result.ParseResult.HasErrors() {
			fmt.Fprintf(os.Stderr, "%s: cannot format file with parse errors\n", filename)
			continue
		}

		if result.ParseResult.Program != nil {
			formatted := c.formatProgram(result.ParseResult.Program)
			
			if c.Output == "" {
				fmt.Printf("=== %s ===\n%s\n", filename, formatted)
			} else {
				if err := c.writeOutput(formatted); err != nil {
					return fmt.Errorf("failed to write output: %w", err)
				}
			}
		}
	}

	return nil
}

type ParseFileResult struct {
	Filename    string              `json:"filename"`
	Error       error               `json:"error,omitempty"`
	ParseResult *parser.ParseResult `json:"parse_result,omitempty"`
}

func (c *CLI) parseFile(p *parser.Parser, filename string) ParseFileResult {
	content, err := os.ReadFile(filename)
	if err != nil {
		return ParseFileResult{
			Filename: filename,
			Error:    err,
		}
	}

	result := p.ParseWithErrors(string(content))
	return ParseFileResult{
		Filename:    filename,
		ParseResult: result,
	}
}

func (c *CLI) formatParseResult(filename string, result ParseFileResult) string {
	var sb strings.Builder

	if c.Format == "json" {
		return "" // JSON output handled separately
	}

	sb.WriteString(fmt.Sprintf("=== %s ===\n", filename))

	if result.Error != nil {
		sb.WriteString(fmt.Sprintf("Error: %v\n", result.Error))
		return sb.String()
	}

	if result.ParseResult.HasErrors() {
		sb.WriteString(fmt.Sprintf("Found %d errors:\n", len(result.ParseResult.Errors)))
		for i, err := range result.ParseResult.Errors {
			sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, err.Error()))
		}
	} else {
		sb.WriteString("✓ Parsing successful\n")
	}

	if result.ParseResult.Program != nil {
		sb.WriteString(fmt.Sprintf("Statements: %d\n", len(result.ParseResult.Program.Statements)))
		if result.ParseResult.Program.Version != nil {
			sb.WriteString(fmt.Sprintf("Version: %s\n", result.ParseResult.Program.Version.Number))
		}
		if len(result.ParseResult.Program.Comments) > 0 {
			sb.WriteString(fmt.Sprintf("Comments: %d\n", len(result.ParseResult.Program.Comments)))
		}
	}

	sb.WriteString("\n")
	return sb.String()
}

func (c *CLI) formatAST(filename string, program *parser.Program, depth int, compact bool) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("=== AST for %s ===\n", filename))
	
	if c.Format == "json" {
		if jsonData, err := json.MarshalIndent(program, "", "  "); err == nil {
			sb.WriteString(string(jsonData))
		} else {
			sb.WriteString("Error serializing AST to JSON")
		}
	} else if c.Format == "tree" {
		sb.WriteString(c.formatASTTree(program, 0, depth, compact))
	} else {
		sb.WriteString(c.formatASTNode(program, 0, depth, compact))
	}
	
	sb.WriteString("\n")
	return sb.String()
}

func (c *CLI) formatASTNode(node parser.Node, currentDepth, maxDepth int, compact bool) string {
	if maxDepth >= 0 && currentDepth > maxDepth {
		return ""
	}

	indent := strings.Repeat("  ", currentDepth)
	var sb strings.Builder

	switch n := node.(type) {
	case *parser.Program:
		sb.WriteString(fmt.Sprintf("%sProgram (%d statements)\n", indent, len(n.Statements)))
		if n.Version != nil && (maxDepth < 0 || currentDepth < maxDepth) {
			sb.WriteString(c.formatASTNode(n.Version, currentDepth+1, maxDepth, compact))
		}
		for i, stmt := range n.Statements {
			if maxDepth >= 0 && currentDepth+1 > maxDepth {
				break
			}
			if i > 10 && compact {
				sb.WriteString(fmt.Sprintf("%s  ... (%d more statements)\n", indent, len(n.Statements)-i))
				break
			}
			sb.WriteString(c.formatASTNode(stmt, currentDepth+1, maxDepth, compact))
		}
	default:
		sb.WriteString(fmt.Sprintf("%s%s\n", indent, node.String()))
	}

	return sb.String()
}

func (c *CLI) formatASTTree(node parser.Node, currentDepth, maxDepth int, compact bool) string {
	if maxDepth >= 0 && currentDepth > maxDepth {
		return ""
	}

	var sb strings.Builder
	prefix := ""
	
	for i := 0; i < currentDepth; i++ {
		if i == currentDepth-1 {
			prefix += "├── "
		} else {
			prefix += "│   "
		}
	}

	switch n := node.(type) {
	case *parser.Program:
		sb.WriteString(fmt.Sprintf("%sProgram (%d statements)\n", prefix, len(n.Statements)))
		if n.Version != nil && (maxDepth < 0 || currentDepth < maxDepth) {
			sb.WriteString(c.formatASTTree(n.Version, currentDepth+1, maxDepth, compact))
		}
		for i, stmt := range n.Statements {
			if maxDepth >= 0 && currentDepth+1 > maxDepth {
				break
			}
			if i > 10 && compact {
				sb.WriteString(fmt.Sprintf("%s└── ... (%d more statements)\n", prefix, len(n.Statements)-i))
				break
			}
			sb.WriteString(c.formatASTTree(stmt, currentDepth+1, maxDepth, compact))
		}
	default:
		sb.WriteString(fmt.Sprintf("%s%s\n", prefix, node.String()))
	}

	return sb.String()
}

type ProgramStats struct {
	TotalStatements       int      `json:"total_statements"`
	QuantumDeclarations   int      `json:"quantum_declarations"`
	ClassicalDeclarations int      `json:"classical_declarations"`
	GateCalls             int      `json:"gate_calls"`
	Measurements          int      `json:"measurements"`
	Includes              int      `json:"includes"`
	Comments              int      `json:"comments"`
	UniqueGates           []string `json:"unique_gates"`
	MaxQubitIndex         int      `json:"max_qubit_index,omitempty"`
	Version               string   `json:"version,omitempty"`
}

func (c *CLI) analyzeProgram(program *parser.Program) ProgramStats {
	stats := ProgramStats{
		TotalStatements: len(program.Statements),
		Comments:        len(program.Comments),
		UniqueGates:     make([]string, 0),
		MaxQubitIndex:   -1,
	}

	if program.Version != nil {
		stats.Version = program.Version.Number
	}

	gateSet := make(map[string]bool)

	// Simplified analysis - would use visitor pattern in full implementation
	for _, stmt := range program.Statements {
		switch s := stmt.(type) {
		case *parser.QuantumDeclaration:
			stats.QuantumDeclarations++
		case *parser.ClassicalDeclaration:
			stats.ClassicalDeclarations++
		case *parser.GateCall:
			stats.GateCalls++
			if !gateSet[s.Name] {
				gateSet[s.Name] = true
				stats.UniqueGates = append(stats.UniqueGates, s.Name)
			}
		case *parser.Measurement:
			stats.Measurements++
		case *parser.Include:
			stats.Includes++
		}
	}

	return stats
}

func (c *CLI) formatStats(filename string, stats ProgramStats) string {
	var sb strings.Builder

	if c.Format == "json" {
		if jsonData, err := json.MarshalIndent(stats, "", "  "); err == nil {
			sb.WriteString(string(jsonData))
		}
	} else {
		sb.WriteString(fmt.Sprintf("=== Statistics for %s ===\n", filename))
		
		if stats.Version != "" {
			sb.WriteString(fmt.Sprintf("OpenQASM version: %s\n", stats.Version))
		}
		
		sb.WriteString(fmt.Sprintf("Total statements: %d\n", stats.TotalStatements))
		sb.WriteString(fmt.Sprintf("  ├── Quantum declarations: %d\n", stats.QuantumDeclarations))
		sb.WriteString(fmt.Sprintf("  ├── Classical declarations: %d\n", stats.ClassicalDeclarations))
		sb.WriteString(fmt.Sprintf("  ├── Gate calls: %d\n", stats.GateCalls))
		sb.WriteString(fmt.Sprintf("  ├── Measurements: %d\n", stats.Measurements))
		sb.WriteString(fmt.Sprintf("  └── Includes: %d\n", stats.Includes))
		
		if stats.Comments > 0 {
			sb.WriteString(fmt.Sprintf("Comments: %d\n", stats.Comments))
		}
		
		if len(stats.UniqueGates) > 0 {
			sb.WriteString(fmt.Sprintf("Unique gates (%d): %s\n", len(stats.UniqueGates), strings.Join(stats.UniqueGates, ", ")))
		}
		
		if stats.MaxQubitIndex >= 0 {
			sb.WriteString(fmt.Sprintf("Max qubit index: %d\n", stats.MaxQubitIndex))
		}
	}

	sb.WriteString("\n")
	return sb.String()
}

func (c *CLI) formatProgram(program *parser.Program) string {
	// Basic program formatting - this is a simplified version
	var sb strings.Builder

	if program.Version != nil {
		sb.WriteString(fmt.Sprintf("OPENQASM %s;\n\n", program.Version.Number))
	}

	// Group includes at the top
	for _, stmt := range program.Statements {
		if include, ok := stmt.(*parser.Include); ok {
			sb.WriteString(fmt.Sprintf("include \"%s\";\n", include.Path))
		}
	}

	if c.hasIncludes(program.Statements) {
		sb.WriteString("\n")
	}

	// Then declarations
	for _, stmt := range program.Statements {
		switch s := stmt.(type) {
		case *parser.QuantumDeclaration:
			if s.Size != nil {
				sb.WriteString(fmt.Sprintf("%s[size] %s;\n", s.Type, s.Identifier))
			} else {
				sb.WriteString(fmt.Sprintf("%s %s;\n", s.Type, s.Identifier))
			}
		case *parser.ClassicalDeclaration:
			if s.Size != nil {
				sb.WriteString(fmt.Sprintf("%s[size] %s;\n", s.Type, s.Identifier))
			} else {
				sb.WriteString(fmt.Sprintf("%s %s;\n", s.Type, s.Identifier))
			}
		}
	}

	if c.hasDeclarations(program.Statements) {
		sb.WriteString("\n")
	}

	// Then gate calls and other statements
	for _, stmt := range program.Statements {
		switch s := stmt.(type) {
		case *parser.GateCall:
			sb.WriteString(fmt.Sprintf("%s qubits;\n", s.Name))
		case *parser.Measurement:
			sb.WriteString("measure qubit -> target;\n")
		case *parser.Include, *parser.QuantumDeclaration, *parser.ClassicalDeclaration:
			// Already handled above
		default:
			sb.WriteString(fmt.Sprintf("// %s\n", stmt.String()))
		}
	}

	return sb.String()
}

func (c *CLI) hasIncludes(statements []parser.Statement) bool {
	for _, stmt := range statements {
		if _, ok := stmt.(*parser.Include); ok {
			return true
		}
	}
	return false
}

func (c *CLI) hasDeclarations(statements []parser.Statement) bool {
	for _, stmt := range statements {
		switch stmt.(type) {
		case *parser.QuantumDeclaration, *parser.ClassicalDeclaration:
			return true
		}
	}
	return false
}
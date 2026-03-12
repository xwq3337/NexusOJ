package formatter

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// CommandFormatter is a base type for formatters that execute external commands
type CommandFormatter struct {
	name      string
	command   string
	languages []string
	args      []string
	timeout   time.Duration
}

// NewCommandFormatter creates a new command-based formatter
func NewCommandFormatter(name, command string, languages []string, args []string) *CommandFormatter {
	return &CommandFormatter{
		name:      name,
		command:   command,
		languages: languages,
		args:      args,
		timeout:   30 * time.Second, // Default timeout
	}
}

// Name returns the formatter name
func (f *CommandFormatter) Name() string {
	return f.name
}

// Languages returns supported languages
func (f *CommandFormatter) Languages() []string {
	return f.languages
}

// Available checks if the formatter command is available
func (f *CommandFormatter) Available() bool {
	_, err := exec.LookPath(f.command)
	return err == nil
}

// SetTimeout sets the timeout for formatter execution
func (f *CommandFormatter) SetTimeout(timeout time.Duration) {
	f.timeout = timeout
}

// Format formats code using the external command
func (f *CommandFormatter) Format(ctx context.Context, code string, options FormatOptions) (string, error) {
	if strings.TrimSpace(code) == "" {
		return "", ErrEmptyCode
	}

	// Create context with timeout
	cmdCtx, cancel := context.WithTimeout(ctx, f.timeout)
	defer cancel()

	// Build command arguments
	args := f.buildArgs(options)

	// Create command
	cmd := exec.CommandContext(cmdCtx, f.command, args...)

	// Set up stdin and stdout
	cmd.Stdin = strings.NewReader(code)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run command
	err := cmd.Run()
	if err != nil {
		return "", &FormatExecError{
			Command: f.command,
			Args:    args,
			Stderr:  stderr.String(),
			Err:     err,
		}
	}

	result := stdout.String()

	// Handle formatters that output to file
	return f.postProcess(result, code), nil
}

// buildArgs constructs command arguments
func (f *CommandFormatter) buildArgs(options FormatOptions) []string {
	args := append([]string{}, f.args...)

	// Add any formatter-specific options
	if len(options.FormatterArgs) > 0 {
		args = append(args, options.FormatterArgs...)
	}

	return args
}

// postProcess handles any post-processing needed
func (f *CommandFormatter) postProcess(output, original string) string {
	return strings.TrimSpace(output)
}

// FormatExecError represents an error during formatter execution
type FormatExecError struct {
	Command string
	Args    []string
	Stderr  string
	Err     error
}

func (e *FormatExecError) Error() string {
	return fmt.Sprintf("exec error: %s %v: %s\nstderr: %s", e.Command, e.Args, e.Err, e.Stderr)
}

func (e *FormatExecError) Unwrap() error {
	return e.Err
}

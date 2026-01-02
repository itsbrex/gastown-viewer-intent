package beads

import (
	"errors"
	"fmt"
)

// BDNotFoundError indicates the bd CLI is not installed or not in PATH.
type BDNotFoundError struct{}

func (e *BDNotFoundError) Error() string {
	return "bd CLI not found in PATH. Install from https://github.com/intent-solutions-io/beads"
}

// NotInitializedError indicates beads is not initialized in the directory.
type NotInitializedError struct {
	Message string
}

func (e *NotInitializedError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("beads not initialized: %s", e.Message)
	}
	return "beads not initialized. Run 'bd init' in your project directory"
}

// NotFoundError indicates the requested issue was not found.
type NotFoundError struct {
	ID string
}

func (e *NotFoundError) Error() string {
	if e.ID != "" {
		return fmt.Sprintf("issue not found: %s", e.ID)
	}
	return "issue not found"
}

// ExecutionError indicates a bd command failed.
type ExecutionError struct {
	Command string
	Stderr  string
	Err     error
}

func (e *ExecutionError) Error() string {
	if e.Stderr != "" {
		return fmt.Sprintf("bd %s failed: %s", e.Command, e.Stderr)
	}
	return fmt.Sprintf("bd %s failed: %v", e.Command, e.Err)
}

func (e *ExecutionError) Unwrap() error {
	return e.Err
}

// ParseError indicates failure to parse bd output.
type ParseError struct {
	Command string
	Err     error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("failed to parse bd %s output: %v", e.Command, e.Err)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

// IsBDNotFoundError checks if the error indicates bd is not installed.
func IsBDNotFoundError(err error) bool {
	var e *BDNotFoundError
	return errors.As(err, &e)
}

// IsNotInitializedError checks if the error indicates beads is not initialized.
func IsNotInitializedError(err error) bool {
	var e *NotInitializedError
	return errors.As(err, &e)
}

// IsNotFoundError checks if the error indicates an issue was not found.
func IsNotFoundError(err error) bool {
	var e *NotFoundError
	return errors.As(err, &e)
}

// IsParseError checks if the error is a parse error.
func IsParseError(err error) bool {
	var e *ParseError
	return errors.As(err, &e)
}

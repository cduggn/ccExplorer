package flags

import (
	"fmt"
	"strings"
)

// Validator defines the interface for validating and parsing flag values
type Validator[T any] interface {
	Validate(value string) (T, error)
	AllowedValues() []string
	Type() string
}

// Flag represents a generic command-line flag with type-safe validation
type Flag[T any, V Validator[T]] struct {
	value     T
	validator V
	isSet     bool
}

// NewFlag creates a new generic flag with the specified validator
func NewFlag[T any, V Validator[T]](validator V) *Flag[T, V] {
	return &Flag[T, V]{
		validator: validator,
	}
}

// Set implements the pflag.Value interface for command-line parsing
func (f *Flag[T, V]) Set(value string) error {
	parsed, err := f.validator.Validate(value)
	if err != nil {
		return err
	}
	f.value = parsed
	f.isSet = true
	return nil
}

// String implements the pflag.Value interface
func (f *Flag[T, V]) String() string {
	return fmt.Sprintf("%v", f.value)
}

// Type implements the pflag.Value interface
func (f *Flag[T, V]) Type() string {
	return f.validator.Type()
}

// Value returns the parsed and validated value
func (f *Flag[T, V]) Value() T {
	return f.value
}

// IsSet returns true if the flag was explicitly set
func (f *Flag[T, V]) IsSet() bool {
	return f.isSet
}

// ValidationError represents a validation error with context
type ValidationError struct {
	Field    string
	Value    string
	Allowed  []string
	Message  string
}

func (e ValidationError) Error() string {
	if len(e.Allowed) > 0 {
		return fmt.Sprintf("invalid %s: %s. Must be one of: %s", 
			e.Field, e.Value, strings.Join(e.Allowed, ", "))
	}
	return fmt.Sprintf("invalid %s: %s. %s", e.Field, e.Value, e.Message)
}

// Commonly used constraint types for AWS resources
type AWSResourceType interface {
	~string
}

type DimensionType string
type TagType string
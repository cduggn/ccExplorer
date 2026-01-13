package writer

import (
	"io"
)

// Writer defines a generic interface for transforming input data to output and rendering it
type Writer[TInput, TOutput any] interface {
	Transform(input TInput) (TOutput, error)
	Render(output TOutput) error
}

// Renderer defines a generic interface for rendering output data  
type Renderer[T any] interface {
	Render(data T) error
}

// Transformer defines a generic interface for data transformation
type Transformer[TInput, TOutput any] interface {
	Transform(input TInput) (TOutput, error)
}

// TableRenderer provides generic table rendering capabilities
type TableRenderer[T any] struct {
	output io.Writer
	config TableConfig[T]
}

// TableConfig defines configuration for table rendering
type TableConfig[T any] struct {
	Columns    []Column[T]
	Style      TableStyle
	HeaderFunc func() []string
	RowFunc    func(T) []string
	FooterFunc func(data []T) []string
}

// Column defines a table column with generic data extraction
type Column[T any] struct {
	Header    string
	Extractor func(T) string
	Width     int
}

// TableStyle defines styling options for tables
type TableStyle struct {
	ShowHeader    bool
	ShowFooter    bool
	ShowBorders   bool
	AlternateRows bool
}

// NewTableRenderer creates a new generic table renderer
func NewTableRenderer[T any](output io.Writer, config TableConfig[T]) *TableRenderer[T] {
	return &TableRenderer[T]{
		output: output,
		config: config,
	}
}

// Render implements the Renderer interface for tables
func (r *TableRenderer[T]) Render(data []T) error {
	// Implementation will be added when we replace specific table implementations
	return nil
}

// CompositeWriter combines transformation and rendering in a single type
type CompositeWriter[TInput, TOutput any] struct {
	transformer Transformer[TInput, TOutput]
	renderer    Renderer[TOutput]
}

// NewCompositeWriter creates a writer that combines transformation and rendering
func NewCompositeWriter[TInput, TOutput any](
	transformer Transformer[TInput, TOutput],
	renderer Renderer[TOutput],
) *CompositeWriter[TInput, TOutput] {
	return &CompositeWriter[TInput, TOutput]{
		transformer: transformer,
		renderer:    renderer,
	}
}

// Transform transforms input data to output format
func (w *CompositeWriter[TInput, TOutput]) Transform(input TInput) (TOutput, error) {
	return w.transformer.Transform(input)
}

// Render renders the output data
func (w *CompositeWriter[TInput, TOutput]) Render(output TOutput) error {
	return w.renderer.Render(output)
}

// Write provides a convenience method that transforms and renders in one call
func (w *CompositeWriter[TInput, TOutput]) Write(input TInput) error {
	output, err := w.Transform(input)
	if err != nil {
		return err
	}
	return w.Render(output)
}
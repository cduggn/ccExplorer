package writer

// Renderer defines a generic interface for rendering output data
type Renderer[T any] interface {
	Render(data T) error
}

// Transformer defines a generic interface for data transformation
type Transformer[TInput, TOutput any] interface {
	Transform(input TInput) (TOutput, error)
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

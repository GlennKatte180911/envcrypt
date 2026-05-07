package env

// Pipeline chains multiple transformation functions over a slice of Entry,
// applying each stage in order and returning the final result.
//
// Each stage is a function that accepts and returns []Entry, making it easy
// to compose existing helpers such as Filter, Transform, Redact, and Validate.
type PipelineStage func([]Entry) ([]Entry, error)

// Pipeline holds an ordered sequence of stages to apply to entries.
type Pipeline struct {
	stages []PipelineStage
}

// NewPipeline creates an empty Pipeline.
func NewPipeline() *Pipeline {
	return &Pipeline{}
}

// Pipe appends a stage to the pipeline and returns the same Pipeline so
// calls can be chained fluently.
func (p *Pipeline) Pipe(stage PipelineStage) *Pipeline {
	p.stages = append(p.stages, stage)
	return p
}

// PipeFunc is a convenience wrapper that promotes a plain transformation
// function ([]Entry) []Entry into a PipelineStage with a nil error.
func PipeFunc(fn func([]Entry) []Entry) PipelineStage {
	return func(entries []Entry) ([]Entry, error) {
		return fn(entries), nil
	}
}

// Run executes all stages in order against entries.  If any stage returns a
// non-nil error the pipeline stops and that error is returned together with
// whatever entries were produced up to that point.
func (p *Pipeline) Run(entries []Entry) ([]Entry, error) {
	current := make([]Entry, len(entries))
	copy(current, entries)

	for _, stage := range p.stages {
		var err error
		current, err = stage(current)
		if err != nil {
			return current, err
		}
	}
	return current, nil
}

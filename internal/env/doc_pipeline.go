// Package env — Pipeline
//
// Pipeline provides a composable, ordered chain of transformation stages for
// []Entry slices.  Each stage is a function with the signature:
//
//	func([]Entry) ([]Entry, error)
//
// Stages are appended with Pipe and executed in order by Run.  If any stage
// returns an error the pipeline halts immediately.
//
// # Basic usage
//
//	out, err := env.NewPipeline().
//		Pipe(env.PipeFunc(env.UppercaseKeys)).
//		Pipe(env.PipeFunc(env.DropEmpty)).
//		Pipe(env.PipeFunc(func(e []env.Entry) []env.Entry {
//			return env.Filter(e, env.WithPrefix("APP_"))
//		})).
//		Run(entries)
//
// PipeFunc is a convenience wrapper for plain (non-error-returning)
// transformation functions so they can be used as pipeline stages without
// boilerplate.
//
// The input slice passed to Run is never mutated; each stage receives a copy
// of the current working set.
package env

package batch

import "context"

// Batch is the interface to process in context.
type Batch interface {
	Run(context.Context)
}

// Worker is the func to process
type Worker func(context.Context)

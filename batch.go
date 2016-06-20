package batch

import "golang.org/x/net/context"

// Batch is the interface to process in context.
type Batch interface {
	Run(context.Context)
}

// Worker is the func to process
type Worker func(context.Context)

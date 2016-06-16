package batch

import (
	"golang.org/x/net/context"
)

// Result of the Worker
type Result int

const (
	// NotProcessed means the Worker did NOT processed anything.
	NotProcessed Result = iota
	// Processed means the Worker processed something.
	Processed
)

// Worker is the body of the batch.
type Worker func(ctx context.Context) Result

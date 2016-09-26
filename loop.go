package batch

import "context"

// Loop the Worker with the interval.
type Loop struct {
	Worker func(context.Context)
	Sig
}

// Run loops the worker
func (l *Loop) Run(ctx context.Context) {
	ctx = l.context(ctx)
	for {
		l.Worker(ctx)

		select {
		case <-ctx.Done():
			return
		default:
			continue
		}
	}
}

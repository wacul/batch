package batch

import "golang.org/x/net/context"

// Loop the Worker with the interval.
type Loop struct {
	Worker func(context.Context)
	Sig
}

// Run loops the worker
func (l *Loop) Run(ctx context.Context) {
	for {
		l.Sig.run(ctx, l.Worker)

		select {
		case <-ctx.Done():
			return
		default:
			continue
		}
	}
}

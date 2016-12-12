package batch

import "context"

// Signal the Worker with the cancelation
type Signal struct {
	Worker func(context.Context)
	Sig
}

// Run worker with the cancelable context
func (s *Signal) Run(ctx context.Context) {
	ctx, cancel := s.context(ctx)
	defer cancel()
	s.Worker(ctx)
}

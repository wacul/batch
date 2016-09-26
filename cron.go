package batch

import (
	"context"
	"time"

	"github.com/gorhill/cronexpr"
)

// Cron will process the Worker on the schedule
type Cron struct {
	Worker    func(context.Context)
	Expr      string
	Location  *time.Location
	Immediate bool // Run the worker immidiately.
	Once      bool // Run the worker at once.
	Sig
}

// Run the Worker on the schedule
func (c *Cron) Run(ctx context.Context) {
	ctx = c.context(ctx)

	immediate := make(chan struct{}, 1)
	if c.Immediate {
		immediate <- struct{}{}
	}

	for {
		loc := time.UTC
		if c.Location != nil {
			loc = c.Location
		}

		now := time.Now().In(loc)
		next := cronexpr.MustParse(c.Expr).Next(now)
		wait := next.Sub(now)
		if next.IsZero() {
			return
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(wait):
			c.Worker(ctx)
		case <-immediate:
			c.Worker(ctx)
		}
		if c.Once {
			break
		}
	}
}

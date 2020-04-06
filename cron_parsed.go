package batch

import (
	"context"
	"time"

	"github.com/gorhill/cronexpr"
)

type cronParsed struct {
	Worker    func(context.Context)
	Expr      *cronexpr.Expression
	Location  *time.Location
	Immediate bool // Run the worker immidiately.
	Once      bool // Run the worker at once.
	Sig
}

// Run the Worker on the schedule
func (c *cronParsed) Run(ctx context.Context) {
	ctx, cancel := c.context(ctx)
	defer cancel()

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
		next := c.Expr.Next(now)
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

package batch

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorhill/cronexpr"
	"golang.org/x/net/context"
)

// Cron will process the Worker on the schedule
type Cron struct {
	Expr      string
	Worker    Worker
	Location  *time.Location
	Immediate bool // Run the worker immidiately.
	Once      bool // Run the worker at once.
}

func (c *Cron) do(ctx context.Context) {
	immediate := make(chan struct{})
	if c.Immediate {
		close(immediate)
	}
	for {
		loc := time.UTC
		if c.Location != nil {
			loc = c.Location
		}
		nextTime := cronexpr.MustParse(c.Expr).Next(time.Now().In(loc))
		wait := nextTime.Sub(time.Now().In(loc))
		if nextTime.IsZero() {
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

// Run the Worker on the schedule.
// Spending the schedule in the Worker, next will be skipped.
func (c *Cron) Run() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGHUP,
	)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-sigCh:
			fmt.Println("cron terminated by signal")
			cancel()
		}
	}()
	c.do(ctx)
}

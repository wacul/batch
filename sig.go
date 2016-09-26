package batch

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Sig the Workers in parallel with the interval.
type Sig struct {
	IgnoreSig bool
}

// Run a batch worker.
// Receiving a SIGTERM, it waits for all of the Worker finished and stop.
// The context canceled, Workers should finish safely and quickly.
func (l *Sig) context(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancel := context.WithCancel(ctx)
	if !l.IgnoreSig {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh,
			syscall.SIGTERM,
			syscall.SIGINT,
			syscall.SIGQUIT,
			syscall.SIGHUP,
		)
		go func() {
			select {
			case <-sigCh:
				cancel()
			}
		}()
	}
	return ctx
}

package batch

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/net/context"
)

var defaultParallels = 1

// Delayer decides a interval depending on the batch result.
type Delayer func(processed RunStatus) time.Duration

var defaultDelayFunc = func(status RunStatus) time.Duration {
	if status == NotRun {
		return 1000 * time.Millisecond
	}
	return 0
}

// Loop the Workers in parallel with the interval.
type Loop struct {
	Worker       Worker
	Delayer      Delayer
	Parallels    int
	IgnoreSignal bool
	lock         sync.Mutex
}

func (l *Loop) delayer() Delayer {
	d := l.Delayer
	if l.Delayer == nil {
		return defaultDelayFunc
	}
	// synchronized delayer
	return func(status RunStatus) time.Duration {
		l.lock.Lock()
		defer l.lock.Unlock()
		return d(status)
	}
}

func (l *Loop) do(ctx context.Context) {
	p := l.Parallels
	if p < 1 {
		p = defaultParallels
	}

	wg := sync.WaitGroup{}
	wg.Add(p)
	for i := 0; i < p; i++ {
		go func() {
			defer wg.Done()
			for {
				run := l.Worker(ctx)
				select {
				case <-time.After(l.delayer()(run)):
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	wg.Wait()
}

// Run a batch worker.
// Receiving a SIGTERM, it waits for all of the Worker finished and stop.
// The context canceled, Workers should finish safely and quickly.
func (l *Loop) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	if !l.IgnoreSignal {
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
	l.do(ctx)
}

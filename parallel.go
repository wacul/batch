package batch

import (
	"sync"

	"golang.org/x/net/context"
)

// Parallel processes the Worker in parallel.
type Parallel struct {
	Worker    func(context.Context)
	Parallels int
	Sig
}

// Run proccesses the worker in parallel.
func (p *Parallel) Run(ctx context.Context) {
	n := p.Parallels
	if n < 1 {
		n = 1
	}

	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			p.Sig.run(ctx, p.Worker)
			defer wg.Done()
		}()
	}
	wg.Wait()
}

package batch

import (
	"context"
	"sync"
)

// Parallel processes the Worker in parallel.
type Parallel struct {
	Worker    func(context.Context)
	Parallels int
	Sig
}

// Run proccesses the worker in parallel.
func (p *Parallel) Run(ctx context.Context) {
	ctx, cancel := p.context(ctx)
	defer cancel()
	n := p.Parallels
	if n < 1 {
		n = 1
	}

	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			p.Worker(ctx)
			defer wg.Done()
		}()
	}
	wg.Wait()
}

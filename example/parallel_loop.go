package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/wacul/batch"
	"github.com/wacul/ptr"
	"golang.org/x/net/context"
)

func nonsense(ctx context.Context) {
	name := nameFromContext(ctx)
	counter := counterFromContext(ctx)
	interval := intervalFromContext(ctx)
	*counter++
	fmt.Printf("I'm %s @%d (%s)\n", name, *counter, interval.String())

	select {
	case <-time.After(interval):
		return
	case <-ctx.Done():
		return
	}
}

func main() {
	para := &batch.Parallel{
		Parallels: 5,
		Worker: func(ctx context.Context) {
			name := getRandomName()
			ctx = nameContext(ctx, name)
			counter := 0
			ctx = counterContext(ctx, &counter)
			interval := getRandomInterval()
			ctx = intervalContext(ctx, interval)

			loop := &batch.Loop{
				Worker: nonsense,
			}
			loop.Run(ctx)
		},
	}
	para.Run(nil)
}

var (
	r = rand.New(rand.NewSource(time.Now().Unix()))
)

func getRandomInterval() time.Duration {
	return time.Duration(rand.Intn(5)+3) * time.Second
}

func getRandomName() string {
	res, err := http.Get("http://uinames.com/api/?region=canada")
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(res.Body)
	var prof struct {
		Name string `json:"name"`
	}
	if err := decoder.Decode(&prof); err != nil {
		panic(err)
	}
	return prof.Name
}

var (
	nameKey     = ptr.String("")
	counterKey  = ptr.String("")
	intervalKey = ptr.String("")
)

func nameFromContext(ctx context.Context) string {
	return ctx.Value(nameKey).(string)
}

func counterFromContext(ctx context.Context) *int {
	return ctx.Value(counterKey).(*int)
}

func intervalFromContext(ctx context.Context) time.Duration {
	return ctx.Value(intervalKey).(time.Duration)
}

func nameContext(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, nameKey, v)
}

func counterContext(ctx context.Context, v *int) context.Context {
	return context.WithValue(ctx, counterKey, v)
}

func intervalContext(ctx context.Context, v time.Duration) context.Context {
	return context.WithValue(ctx, intervalKey, v)
}

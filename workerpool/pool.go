package workerpool

import (
	"context"
	"fmt"
	"sync"

	"github.com/bhuvana-chinnadurai/cli-http-workerpool-go/client"
)

type Result struct {
	Value interface{}
	Error error
}

type Pool struct {
	maxParallelCount int
	URLs             chan string
	Results          chan Result
}

const (
	DefaultMaxParallelWorkers = 10
)

func New(maxParallelCount int) Pool {
	if maxParallelCount == 0 {
		maxParallelCount = DefaultMaxParallelWorkers
	}
	return Pool{
		maxParallelCount: maxParallelCount,
		URLs:             make(chan string, maxParallelCount),
		Results:          make(chan Result, maxParallelCount),
	}
}

func (wp Pool) Run(ctx context.Context, client client.Client) {
	wg := &sync.WaitGroup{}
	for i := 0; i < wp.maxParallelCount; i++ {
		wg.Add(1)
		go makeRequest(ctx, client, wg, wp.URLs, wp.Results)
	}
	wg.Wait()
	close(wp.Results)
}

func makeRequest(ctx context.Context, client client.Client, wg *sync.WaitGroup, urls <-chan string, results chan<- Result) {
	defer wg.Done()
	for {
		select {
		case url, ok := <-urls:
			if !ok {
				return
			}
			result, err := client.GetResult(url)
			if err != nil {
				results <- Result{Error: fmt.Errorf("error occurred while gettting response for the url: '%s' : '%s' ", url, err.Error())}
			} else {
				results <- Result{Value: result}
			}
		case <-ctx.Done():
			results <- Result{Error: fmt.Errorf("error while waiting for response: %s", ctx.Err())}
			return
		}
	}

}

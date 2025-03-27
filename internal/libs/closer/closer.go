package closer

import (
	"context"
	"fmt"
	"strings"
	"sync"
)


type closeFunc = func(context.Context) error

// Closer represents a type that captures close functions for connections. 
type Closer struct {
	sync.Mutex
	funcs []closeFunc
}

// NewCloser creates a new Closer
func NewCloser() *Closer {
	return &Closer{}
}

// Add function adds a new closer function 
func (c *Closer) Add(f closeFunc) {
	c.Lock()
	defer c.Unlock()
	c.funcs = append(c.funcs, f)
}

// Close function runs all closer functions
func (c *Closer) Close(ctx context.Context) error {
	c.Lock()
	defer c.Unlock()

	var (
		msgs    = make([]string, 0, len(c.funcs))
		waitgroup      sync.WaitGroup
		errorCh = make(chan error, len(c.funcs))
		done    = make(chan struct{})
	)

	for i := len(c.funcs) - 1; i >= 0; i-- {
		waitgroup.Add(1)
		go func(f closeFunc) {
			defer waitgroup.Done()
			if err := f(ctx); err != nil {
				errorCh <- err
			}
		}(c.funcs[i])
	}

	go func() {
		waitgroup.Wait()
		close(done)
		close(errorCh)
	}()

	select {
	case <-done:
		break
	case <-ctx.Done():
		return fmt.Errorf("shutdown timeout: %v", ctx.Err())
	}

	for err := range errorCh {
		msgs = append(msgs, fmt.Sprintf("[!] %v", err))
	}

	if len(msgs) > 0 {
		return fmt.Errorf(
			"shutdown completed with errors:\n%s",
			strings.Join(msgs, "\n"),
		)
	}

	return nil
}
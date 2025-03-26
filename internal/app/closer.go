package app

import (
	"context"
	"fmt"
	"strings"
	"sync"
)


type Func = func(context.Context) error

type Closer struct {
	sync.Mutex
	funcs []Func
}


func NewCloser() *Closer {
	return &Closer{}
}

func (c *Closer) Add(f Func) {
	c.Lock()
	defer c.Unlock()
	c.funcs = append(c.funcs, f)
}


func (c *Closer) Close(ctx context.Context) error {
	c.Lock()
	defer c.Unlock()

	var (
		msgs    = make([]string, 0, len(c.funcs))
		wg      sync.WaitGroup
		errorCh = make(chan error, len(c.funcs))
		done    = make(chan struct{})
	)

	for i := len(c.funcs) - 1; i >= 0; i-- {
		wg.Add(1)
		go func(f Func) {
			defer wg.Done()
			if err := f(ctx); err != nil {
				errorCh <- err
			}
		}(c.funcs[i])
	}

	go func() {
		wg.Wait()
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
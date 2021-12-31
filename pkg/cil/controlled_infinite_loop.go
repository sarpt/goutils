package cil

import (
	"context"
	"sync"
)

type Cfg struct {
	AfterLoopCb     func()
	Cb              func() error
	Result          chan<- error
	ShouldTerminate func(error) bool
	TerminateCb     func(error)
}

// ControlledInfiniteLoop loops infinitely, calling 'Cb' and 'AfterLoopCb' until provided ctx or error-handling tells it to stop.
// 'Cb' and 'Result' are mandatory to work. When loop terminates due to ctx being done, the reason will be provided in the 'Result'.
// If 'ShouldTerminate' is provided, it handles error-based cancellation.
// If 'ShouldTerminate' is not provided, the error-based cancellation is done based on whether `Cb` returns nil or error.
// 'TerminateCb' is called (when provided) on defer just after the internal go-routine that loops infinitely has finished
// (note: if loop is closed by ctx, the `TerminateCb` will receive nil as error the context-based cancelation reason will be provided to Result).
// 'AfterLoopCb' is called (when provided) after every single `Cb` loop execution that was not terminated.
func ControlledInfiniteLoop(ctx context.Context, cfg Cfg) {
	loopTerminate := make(chan error, 1) // size of 1 since we don't want to leave goroutine blocked indifinitely if ctx.Done is caught earlier

	terminate := false
	terminateLock := &sync.RWMutex{}

	go func() {
		var err error
		defer close(loopTerminate)
		if cfg.TerminateCb != nil {
			defer cfg.TerminateCb(err)
		}

		for {
			err = cfg.Cb()
			shouldTerminate := cfg.ShouldTerminate != nil && cfg.ShouldTerminate(err)
			uncontrolledError := cfg.ShouldTerminate == nil && err != nil
			if shouldTerminate || uncontrolledError {
				loopTerminate <- err
				break
			}

			terminateLock.RLock()
			shouldStop := terminate
			terminateLock.RUnlock()

			if shouldStop {
				break
			}

			if cfg.AfterLoopCb != nil {
				cfg.AfterLoopCb()
			}
		}
	}()

	select {
	case <-ctx.Done():
		terminateLock.Lock()
		terminate = true
		terminateLock.Unlock()
		cfg.Result <- ctx.Err()
	case err := <-loopTerminate:
		cfg.Result <- err
	}
}

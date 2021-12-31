package cil_test

import (
	"context"
	"errors"
	"testing"

	"github.com/sarpt/goutils/pkg/cil"
)

func TestCil_TerminateWithContextCancel(t *testing.T) {
	// given
	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan error)
	cfg := cil.Cfg{
		Cb: func() error {
			defer cancel()
			return nil
		},
		Result: result,
	}

	// when
	go cil.ControlledInfiniteLoop(ctx, cfg)

	err := <-result

	// then
	if err != context.Canceled {
		t.Fatalf("Reason in result is '%s', expected '%s'", err, context.Canceled)
	}
}

func TestCil_TerminateWithErrorFromCallback(t *testing.T) {
	// given
	ctx := context.Background()
	expectedErr := errors.New("error from callback")
	result := make(chan error)
	cfg := cil.Cfg{
		Cb: func() error {
			return expectedErr
		},
		Result: result,
	}

	// when
	go cil.ControlledInfiniteLoop(ctx, cfg)

	err := <-result

	// then
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Result returned error '%s', expected '%s'", err, expectedErr)
	}
}

func TestCil_AfterLoopCallbackCalledWithShouldTerminateCheck(t *testing.T) {
	// given
	ctx := context.Background()
	result := make(chan error)
	count := 0
	expectedCount := 10
	cfg := cil.Cfg{
		AfterLoopCb: func() {
			count += 1
		},
		ShouldTerminate: func(err error) bool {
			return count == expectedCount
		},
		Cb: func() error {
			return nil
		},
		Result: result,
	}

	// when
	go cil.ControlledInfiniteLoop(ctx, cfg)

	<-result

	// then
	if count != expectedCount {
		t.Fatalf("Expected count to be '%d', it is: '%d'", expectedCount, count)
	}
}

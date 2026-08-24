package main

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeTrayLifecycle struct {
	runErr  error
	stopped chan struct{}
	once    sync.Once
}

func newFakeTrayLifecycle(runErr error) *fakeTrayLifecycle {
	return &fakeTrayLifecycle{runErr: runErr, stopped: make(chan struct{})}
}

func (f *fakeTrayLifecycle) Run(ctx context.Context) error {
	if f.runErr != nil {
		return f.runErr
	}
	select {
	case <-f.stopped:
		return nil
	case <-ctx.Done():
		return nil
	}
}

func (f *fakeTrayLifecycle) Stop() {
	f.once.Do(func() { close(f.stopped) })
}

func TestPrintUsage(t *testing.T) {
	err := run([]string{"help"})
	assert.NoError(t, err)

	err = run([]string{"--help"})
	assert.NoError(t, err)

	err = run([]string{"-h"})
	assert.NoError(t, err)
}

func TestRunVersion(t *testing.T) {
	err := run([]string{"version"})
	assert.NoError(t, err)

	err = run([]string{"-version"})
	assert.NoError(t, err)
}

func TestUnknownCommand(t *testing.T) {
	err := run([]string{"nonexistent-command-xyz"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown command")
}

func TestRunTrayLifecycleStopsTrayWhenServerFails(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	tray := newFakeTrayLifecycle(nil)
	wantErr := errors.New("listen tcp: address already in use")

	err := runTrayLifecycle(ctx, cancel, func(context.Context) error {
		return wantErr
	}, tray)

	assert.ErrorIs(t, err, wantErr)
	select {
	case <-tray.stopped:
	default:
		t.Fatal("tray was not stopped after the server failed")
	}
}

func TestRunTrayLifecycleCancelsServerWhenTrayFails(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	wantErr := errors.New("tray event loop failed")
	tray := newFakeTrayLifecycle(wantErr)
	serverStopped := make(chan struct{})

	err := runTrayLifecycle(ctx, cancel, func(ctx context.Context) error {
		<-ctx.Done()
		close(serverStopped)
		return nil
	}, tray)

	assert.ErrorIs(t, err, wantErr)
	select {
	case <-serverStopped:
	default:
		t.Fatal("server was not cancelled after the tray failed")
	}
}

package gdx_test

import (
	"bytes"
	"context"
	"io"
	"strings"
	"syscall"
	"testing"
	"time"

	diagx "github.com/retrovibed/gdx"
	"github.com/stretchr/testify/require"
)

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error { return nil }

func TestDumpRoutinesInto(t *testing.T) {
	t.Run("writes goroutine stacks and closes dst", func(t *testing.T) {
		var buf bytes.Buffer
		err := diagx.DumpRoutinesInto(nopWriteCloser{&buf})
		require.NoError(t, err)
		require.Contains(t, buf.String(), "goroutine")
	})
}

func TestDumpRoutines(t *testing.T) {
	t.Run("writes to a temp file and returns its path", func(t *testing.T) {
		path, err := diagx.DumpRoutines()
		require.NoError(t, err)
		require.True(t, strings.HasSuffix(path, ".trace") || path == "stderr")
	})
}

func TestOnSignal(t *testing.T) {
	t.Run("runs do on signal and stops when ctx is done", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		calls := make(chan struct{}, 1)
		done := make(chan struct{})
		go func() {
			diagx.OnSignal(ctx, func(ctx context.Context) error {
				calls <- struct{}{}
				return nil
			}, syscall.SIGUSR1)
			close(done)
		}()

		time.Sleep(10 * time.Millisecond) // let signal.Notify register before we raise.
		require.NoError(t, syscall.Kill(syscall.Getpid(), syscall.SIGUSR1))

		select {
		case <-calls:
		case <-time.After(time.Second):
			t.Fatal("do was not invoked in time")
		}

		cancel()

		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("OnSignal did not return after ctx was cancelled")
		}
	})
}

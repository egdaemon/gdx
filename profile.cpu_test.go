package gdx_test

import (
	"context"
	"io"
	"runtime/trace"
	"testing"
	"time"

	"github.com/egdaemon/gdx"
	"github.com/stretchr/testify/require"
)

func TestCPU(t *testing.T) {
	t.Run("streams a gzip-encoded pprof profile", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		out, err := io.ReadAll(gdx.CPU(ctx))
		require.NoError(t, err)
		require.True(t, len(out) >= 2)
		require.Equal(t, gzipMagic, out[:2])
	})

	t.Run("does not leave runtime/trace running", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		_, err := io.ReadAll(gdx.CPU(ctx))
		require.NoError(t, err)

		require.NoError(t, trace.Start(io.Discard))
		trace.Stop()
	})
}

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

var gzipMagic = []byte{0x1f, 0x8b}

func TestProfile(t *testing.T) {
	t.Run("dispatches cpu to a gzip-encoded pprof profile", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		out, err := io.ReadAll(gdx.Profile(ctx, gdx.ProfileMode_cpu))
		require.NoError(t, err)
		require.True(t, len(out) >= 2)
		require.Equal(t, gzipMagic, out[:2])
	})

	t.Run("dispatches heap and mem to the same gzip-encoded heap profile", func(t *testing.T) {
		for _, mode := range []gdx.ProfileMode{gdx.ProfileMode_heap, gdx.ProfileMode_mem} {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
			out, err := io.ReadAll(gdx.Profile(ctx, mode))
			cancel()
			require.NoError(t, err)
			require.True(t, len(out) >= 2)
			require.Equal(t, gzipMagic, out[:2])
		}
	})

	t.Run("dispatches allocs to a gzip-encoded allocs profile", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		out, err := io.ReadAll(gdx.Profile(ctx, gdx.ProfileMode_allocs))
		require.NoError(t, err)
		require.True(t, len(out) >= 2)
		require.Equal(t, gzipMagic, out[:2])
	})

	t.Run("dispatches block to a gzip-encoded block profile", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		out, err := io.ReadAll(gdx.Profile(ctx, gdx.ProfileMode_block))
		require.NoError(t, err)
		require.True(t, len(out) >= 2)
		require.Equal(t, gzipMagic, out[:2])
	})

	t.Run("unknown mode errors on read", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		_, err := io.ReadAll(gdx.Profile(ctx, gdx.ProfileMode(99)))
		require.ErrorContains(t, err, "unknown profile mode")
	})
}

func TestCPU(t *testing.T) {
	t.Run("streams a gzip-encoded pprof profile", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		out, err := io.ReadAll(gdx.CPU(ctx))
		require.NoError(t, err)
		require.True(t, len(out) >= 2)
		require.Equal(t, gzipMagic, out[:2])
	})

	t.Run("does not leave runtime/trace running (genieql-lineage bug fix)", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		_, err := io.ReadAll(gdx.CPU(ctx))
		require.NoError(t, err)

		// if CPU had left a trace running (the genieql-lineage bug), this
		// Start would fail with "trace: is enabled".
		require.NoError(t, trace.Start(io.Discard))
		trace.Stop()
	})
}

func TestMemory(t *testing.T) {
	t.Run("streams a gzip-encoded heap profile once ctx ends", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		out, err := io.ReadAll(gdx.Memory(ctx))
		require.NoError(t, err)
		require.True(t, len(out) >= 2)
		require.Equal(t, gzipMagic, out[:2])
	})
}

func TestHeap(t *testing.T) {
	t.Run("streams a gzip-encoded heap profile once ctx ends", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		out, err := io.ReadAll(gdx.Heap(ctx))
		require.NoError(t, err)
		require.True(t, len(out) >= 2)
		require.Equal(t, gzipMagic, out[:2])
	})
}

func TestAllocs(t *testing.T) {
	t.Run("streams a gzip-encoded allocs profile once ctx ends", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		out, err := io.ReadAll(gdx.Allocs(ctx))
		require.NoError(t, err)
		require.True(t, len(out) >= 2)
		require.Equal(t, gzipMagic, out[:2])
	})
}

func TestBlock(t *testing.T) {
	t.Run("streams a gzip-encoded block profile for the duration of ctx", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		out, err := io.ReadAll(gdx.Block(ctx))
		require.NoError(t, err)
		require.True(t, len(out) >= 2)
		require.Equal(t, gzipMagic, out[:2])
	})
}

func TestTrace(t *testing.T) {
	t.Run("streams a non-pprof runtime/trace execution trace", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		out, err := io.ReadAll(gdx.Trace(ctx))
		require.NoError(t, err)
		require.True(t, len(out) >= 2)
		require.NotEqual(t, gzipMagic, out[:2])
	})
}

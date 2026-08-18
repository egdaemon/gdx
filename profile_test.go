package gdx_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/egdaemon/gdx"
	"github.com/stretchr/testify/require"
)

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

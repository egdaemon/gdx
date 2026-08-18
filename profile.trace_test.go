package gdx_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/egdaemon/gdx"
	"github.com/stretchr/testify/require"
)

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

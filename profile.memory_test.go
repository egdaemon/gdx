package gdx_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/egdaemon/gdx"
	"github.com/stretchr/testify/require"
)

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

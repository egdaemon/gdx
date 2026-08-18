package gdx_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/egdaemon/gdx"
	"github.com/stretchr/testify/require"
)

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

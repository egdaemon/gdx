package gdx_test

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/egdaemon/gdx"
	"github.com/stretchr/testify/require"
)

func TestRecordFile(t *testing.T) {
	t.Run("writes cpu profile to disk", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "cpu.pprof")

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := gdx.RecordFile(ctx, path, gdx.Profile(ctx, gdx.ProfileMode_cpu))
		require.NoError(t, err)

		data, err := os.ReadFile(path)
		require.NoError(t, err)
		require.True(t, len(data) >= 2)
		require.Equal(t, gzipMagic, data[:2])
	})

	t.Run("creates parent directory", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "a", "b", "c", "heap.pprof")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := gdx.RecordFile(ctx, path, gdx.Heap(ctx))
		require.NoError(t, err)

		data, err := os.ReadFile(path)
		require.NoError(t, err)
		require.True(t, len(data) >= 2)
		require.Equal(t, gzipMagic, data[:2])
	})

	t.Run("writes heap profile on context cancel", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "heap.pprof")

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(10 * time.Millisecond)
			cancel()
		}()

		err := gdx.RecordFile(ctx, path, gdx.Heap(ctx))
		require.NoError(t, err)

		data, err := os.ReadFile(path)
		require.NoError(t, err)
		require.True(t, len(data) >= 2)
		require.Equal(t, gzipMagic, data[:2])
	})

	t.Run("fails on unwritable path", func(t *testing.T) {
		err := gdx.RecordFile(context.Background(), "/nonexistent/dir/profile.pprof", io.NopCloser(nil))
		require.Error(t, err)
	})
}

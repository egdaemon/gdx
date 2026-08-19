package gdx_test

import (
	"context"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/egdaemon/gdx"
	"github.com/stretchr/testify/require"
)

func TestUnixServeServesHTTP(t *testing.T) {
	dir := t.TempDir()
	socketPath := filepath.Join(dir, "test.sock")

	go gdx.UnixServe(t.Context(), socketPath)

	// Give the server time to start and serve.
	require.Eventually(t, func() bool {
		conn, err := net.DialTimeout("unix", socketPath, 500*time.Millisecond)
		if err != nil {
			return false
		}
		conn.Close()
		return true
	}, 2*time.Second, 50*time.Millisecond, "server did not start in time")

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", socketPath)
			},
		},
	}

	t.Run("/debug/vars responds on unix socket", func(t *testing.T) {
		resp, err := client.Get("http://unixsocket/debug/vars")
		require.NoError(t, err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Contains(t, string(body), `"cmdline"`)
	})

	t.Run("/debug/goroutines responds on unix socket", func(t *testing.T) {
		resp, err := client.Get("http://unixsocket/debug/goroutines")
		require.NoError(t, err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Contains(t, string(body), "goroutine")
	})

	t.Run("removes stale socket file", func(t *testing.T) {
		// Simulate a crashed server by leaving a stale socket file behind.
		staleDir := t.TempDir()
		stalePath := filepath.Join(staleDir, "stale.sock")
		f, err := os.OpenFile(stalePath, os.O_CREATE, 0o755)
		require.NoError(t, err)
		f.Close()

		require.FileExists(t, stalePath)

		ctx, cancel := context.WithCancel(context.Background())
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			gdx.UnixServe(ctx, stalePath)
		}()

		// Server should have cleaned up the stale socket and be serving.
		require.Eventually(t, func() bool {
			conn, err := net.DialTimeout("unix", stalePath, 500*time.Millisecond)
			if err != nil {
				return false
			}
			conn.Close()
			return true
		}, 2*time.Second, 50*time.Millisecond, "server did not start after removing stale socket")

		cancel()
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
		case <-time.After(3 * time.Second):
			t.Fatal("UnixServe did not stop after context cancellation")
		}

		// After shutdown the socket file is removed by net package.
		_, err = os.Stat(stalePath)
		require.True(t, os.IsNotExist(err), "socket file still exists after shutdown")
	})
}

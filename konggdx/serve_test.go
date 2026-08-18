package konggdx

import (
	"context"
	"expvar"
	"io"
	"net"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/alecthomas/kong"
	"github.com/egdaemon/gdx/gdxapi"
	"github.com/stretchr/testify/require"
)

func TestServeRun(t *testing.T) {
	t.Run("serves the gdx debug surface until ctx is cancelled", func(t *testing.T) {
		var serveTestExpvar = expvar.NewString("konggdx_serve_test_expvar_marker")
		serveTestExpvar.Set("hello")

		path := filepath.Join(t.TempDir(), "test.socket")

		l, err := net.Listen("unix", path)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())

		done := make(chan error, 1)
		go func() {
			done <- (Serve{}).run(ctx, l)
		}()

		client := gdxapi.NewUnixClient(path)

		var resp *http.Response
		require.Eventually(t, func() bool {
			resp, err = client.Get("http://unix/debug/vars")
			return err == nil
		}, 2*time.Second, 10*time.Millisecond)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		require.NoError(t, err)
		require.Contains(t, string(body), `"konggdx_serve_test_expvar_marker": "hello"`)

		cancel()
		require.NoError(t, <-done)
	})

	t.Run("--socket flag reaches Serve.Socket via kong", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "flag.socket")

		var cli struct {
			Serve Serve `cmd:""`
		}

		parser := kong.Must(&cli, kong.Vars{"vars_gdx_socket": filepath.Join(t.TempDir(), "unused.socket")})
		kctx, err := parser.Parse([]string{"serve", "--socket", path})
		require.NoError(t, err)
		require.Equal(t, path, cli.Serve.Socket)
		_ = kctx
	})
}

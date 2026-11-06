package konggdx

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/egdaemon/gdx/gdxapi"
	"github.com/stretchr/testify/require"
)

func TestGoroutinesRun(t *testing.T) {
	t.Run("streams the goroutine dump bytes to out", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "test.socket")
		l, err := net.Listen("unix", path)
		require.NoError(t, err)

		srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/debug/goroutines", r.URL.Path)
			w.Write([]byte("fake-goroutine-dump"))
		}))
		srv.Listener = l
		srv.Start()
		defer srv.Close()

		client := gdxapi.NewUnixClient(path)

		var out bytes.Buffer
		cmd := Goroutines{}
		require.NoError(t, cmd.run(context.Background(), client, &out))
		require.Equal(t, "fake-goroutine-dump", out.String())
	})

	t.Run("returns an error on a non-2xx response", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "test.socket")
		l, err := net.Listen("unix", path)
		require.NoError(t, err)

		srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "boom", http.StatusInternalServerError)
		}))
		srv.Listener = l
		srv.Start()
		defer srv.Close()

		client := gdxapi.NewUnixClient(path)

		var out bytes.Buffer
		cmd := Goroutines{}
		err = cmd.run(context.Background(), client, &out)
		require.Error(t, err)
		require.ErrorContains(t, err, "boom")
	})
}

// Do must work over a client for any transport, the socket is only used by Run.
func TestGoroutinesDo(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/debug/goroutines", r.URL.Path)
		_, _ = w.Write([]byte("remote-dump"))
	}))
	defer srv.Close()

	c := gdxapi.NewClient(func(ctx context.Context, _, _ string) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, "tcp", srv.Listener.Addr().String())
	})

	path := filepath.Join(t.TempDir(), "out")
	var cli struct {
		Goroutines Goroutines `cmd:""`
	}
	parser := kong.Must(&cli, kong.Vars{"vars_gdx_socket": "unused", "vars_gdx_default_output": "-"})
	_, err := parser.Parse([]string{"goroutines", "--output", path})
	require.NoError(t, err)

	require.NoError(t, cli.Goroutines.Do(context.Background(), c))
	raw, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, "remote-dump", string(raw))
}

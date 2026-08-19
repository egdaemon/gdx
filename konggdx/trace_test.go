package konggdx

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/alecthomas/kong"
	"github.com/egdaemon/gdx/gdxapi"
	"github.com/stretchr/testify/require"
)

func TestTraceRun(t *testing.T) {
	t.Run("streams the trace bytes to out", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "test.socket")
		l, err := net.Listen("unix", path)
		require.NoError(t, err)

		srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/debug/trace", r.URL.Path)
			require.Equal(t, "5", r.URL.Query().Get("duration"))
			w.Write([]byte("fake-trace-bytes"))
		}))
		srv.Listener = l
		srv.Start()
		defer srv.Close()

		client := gdxapi.NewUnixClient(path)

		var out bytes.Buffer
		cmd := Trace{Duration: 5 * time.Second}
		require.NoError(t, cmd.run(context.Background(), client, &out))
		require.Equal(t, "fake-trace-bytes", out.String())
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
		cmd := Trace{Duration: 5 * time.Second}
		err = cmd.run(context.Background(), client, &out)
		require.Error(t, err)
		require.ErrorContains(t, err, "boom")
	})

	t.Run("--duration flag reaches Trace.Duration via kong", func(t *testing.T) {
		var cli struct {
			Trace Trace `cmd:""`
		}

		parser := kong.Must(&cli, kong.Vars{"vars_gdx_socket": filepath.Join(t.TempDir(), "unused.socket")})
		_, err := parser.Parse([]string{"trace", "--duration", "5s"})
		require.NoError(t, err)
		require.Equal(t, 5*time.Second, cli.Trace.Duration)
	})
}

package gdx_test

import (
	"expvar"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/egdaemon/gdx"
	"github.com/stretchr/testify/require"
)

func TestNewHTTPFn(t *testing.T) {
	srv := httptest.NewServer(gdx.NewHTTPFn(gdx.Options().WithDefaultDuration(0)))
	defer srv.Close()

	t.Run("/debug/vars includes registered expvar.Vars", func(t *testing.T) {
		var serverTestExpvar = expvar.NewString("diagx_test_server_expvar_marker")

		serverTestExpvar.Set("hello")

		resp, err := srv.Client().Get(srv.URL + "/debug/vars")
		require.NoError(t, err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Contains(t, string(body), `"diagx_test_server_expvar_marker": "hello"`)
	})

	t.Run("/debug/goroutines returns a non-empty stack dump", func(t *testing.T) {
		resp, err := srv.Client().Get(srv.URL + "/debug/goroutines")
		require.NoError(t, err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Contains(t, string(body), "goroutine")
	})

	t.Run("/debug/profile/cpu returns a gzip-encoded pprof profile", func(t *testing.T) {
		resp, err := srv.Client().Get(srv.URL + "/debug/profile/cpu?duration=1")
		require.NoError(t, err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.True(t, len(body) >= 2)
		require.Equal(t, []byte{0x1f, 0x8b}, body[:2])
	})

	t.Run("/debug/profile/{unknown mode} returns 500", func(t *testing.T) {
		resp, err := srv.Client().Get(srv.URL + "/debug/profile/nope")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, 500, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Contains(t, string(body), "unknown profile mode")
	})

	t.Run("/debug/trace returns a non-empty, non-pprof execution trace", func(t *testing.T) {
		resp, err := srv.Client().Get(srv.URL + "/debug/trace?duration=1")
		require.NoError(t, err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.True(t, len(body) >= 2)
		require.NotEqual(t, []byte{0x1f, 0x8b}, body[:2])
	})
}

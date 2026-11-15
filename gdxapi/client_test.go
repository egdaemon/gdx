package gdxapi_test

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/egdaemon/gdx/gdxapi"
	"github.com/stretchr/testify/require"
)

// the client must work over any transport, not only unix sockets.
func TestNewClientUsesProvidedDialer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/debug/goroutines", r.URL.Path)
		_, _ = w.Write([]byte("dump"))
	}))
	defer srv.Close()

	var dialed int
	c := gdxapi.NewClient(func(ctx context.Context, network, addr string) (net.Conn, error) {
		dialed++
		return (&net.Dialer{}).DialContext(ctx, "tcp", srv.Listener.Addr().String())
	})

	var out bytes.Buffer
	require.NoError(t, gdxapi.Fetch(context.Background(), c, "/debug/goroutines", time.Second, &out))
	require.Equal(t, "dump", out.String())
	require.Equal(t, 1, dialed)
}

func TestFetchReportsDialFailure(t *testing.T) {
	c := gdxapi.NewClient(func(context.Context, string, string) (net.Conn, error) {
		return nil, net.ErrClosed
	})

	err := gdxapi.Fetch(context.Background(), c, "/debug/goroutines", time.Second, &bytes.Buffer{})
	require.ErrorIs(t, err, net.ErrClosed)
}

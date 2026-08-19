package gdx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAutosocket(t *testing.T) {
	t.Run("handles deeply nested binary paths", func(t *testing.T) {
		// Regression: AutoSocket previously passed os.Args[0] verbatim to
		// RuntimeDirectory. When the binary lived at a deeply nested path like
		// /home/james/development/retrovibed/.../retrovibed, that full path was
		// used as the subdirectory under $XDG_RUNTIME_DIR, producing a unix socket
		// path that exceeded the platform limit and failed with "bind: invalid
		// argument".
		//
		// autosocket must strip the binary path to its basename so the resulting
		// socket path is short enough to bind.
		deeplyNested := "/home/james/development/retrovibed/retrovibed/console/build/linux/x64/debug/bundle/retrovibed/retrovibed"

		path := autosocket(deeplyNested)

		// The socket path should use only the binary name as the subdirectory,
		// not the full nested path.  Verify the immediate parent of gdx.socket
		// is just "retrovibed".
		last := path[len(path)-len("retrovibed/gdx.socket"):]
		require.Equal(t, "retrovibed/gdx.socket", last)
	})

	t.Run("already a bare name works identically", func(t *testing.T) {
		path := autosocket("myapp")

		last := path[len(path)-len("myapp/gdx.socket"):]
		require.Equal(t, "myapp/gdx.socket", last)
	})
}

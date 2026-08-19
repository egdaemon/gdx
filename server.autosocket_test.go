package gdx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAutosocket(t *testing.T) {
	t.Run("handles deeply nested binary paths", func(t *testing.T) {
		// autosocket must strip the binary path to its basename so the resulting
		// socket path is short enough to bind.
		deeplyNested := "/home/user/development/derp/derp/console/build/linux/x64/debug/bundle/derp/derp"

		path := autosocket(deeplyNested)

		last := path[len(path)-len("derp/gdx.socket"):]
		require.Equal(t, "derp/gdx.socket", last)
	})

	t.Run("already a bare name works identically", func(t *testing.T) {
		path := autosocket("myapp")

		last := path[len(path)-len("myapp/gdx.socket"):]
		require.Equal(t, "myapp/gdx.socket", last)
	})
}

package konggdx

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResolveSocket(t *testing.T) {
	// live pids that exist for the duration of the test.
	self := os.Getpid()
	parent := os.Getppid()

	t.Run("existing hint is used as is", func(t *testing.T) {
		dir := t.TempDir()
		hint := filepath.Join(dir, "gdx.socket")
		require.NoError(t, os.WriteFile(hint, nil, 0600))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "gdx."+strconv.Itoa(parent)+".socket"), nil, 0600))

		got, err := resolveSocket(hint, self, false, strings.NewReader(""), &bytes.Buffer{})
		require.NoError(t, err)
		require.Equal(t, hint, got)
	})

	t.Run("filters own pid and uses the single remaining socket", func(t *testing.T) {
		dir := t.TempDir()
		other := filepath.Join(dir, "gdx."+strconv.Itoa(parent)+".socket")
		require.NoError(t, os.WriteFile(filepath.Join(dir, "gdx."+strconv.Itoa(self)+".socket"), nil, 0600))
		require.NoError(t, os.WriteFile(other, nil, 0600))

		got, err := resolveSocket(filepath.Join(dir, "gdx.socket"), self, false, strings.NewReader(""), &bytes.Buffer{})
		require.NoError(t, err)
		require.Equal(t, other, got)
	})

	t.Run("skips sockets of dead processes", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "gdx.4194304.socket"), nil, 0600))

		_, err := resolveSocket(filepath.Join(dir, "gdx.socket"), self, false, strings.NewReader(""), &bytes.Buffer{})
		require.ErrorContains(t, err, "no gdx sockets found")
	})

	t.Run("prompts when multiple sockets remain", func(t *testing.T) {
		dir := t.TempDir()
		paths := []string{
			filepath.Join(dir, "gdx."+strconv.Itoa(self)+".socket"),
			filepath.Join(dir, "gdx."+strconv.Itoa(parent)+".socket"),
		}
		for _, p := range paths {
			require.NoError(t, os.WriteFile(p, nil, 0600))
		}
		sort.Strings(paths)

		// -1 as our pid so neither socket is filtered; an invalid answer is re-prompted.
		var out bytes.Buffer
		got, err := resolveSocket(filepath.Join(dir, "gdx.socket"), -1, true, strings.NewReader("bogus\n2\n"), &out)
		require.NoError(t, err)
		require.Contains(t, out.String(), "multiple gdx sockets found")
		require.Equal(t, paths[1], got)
	})

	t.Run("non interactive with multiple sockets errors listing them", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "gdx."+strconv.Itoa(self)+".socket"), nil, 0600))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "gdx."+strconv.Itoa(parent)+".socket"), nil, 0600))

		_, err := resolveSocket(filepath.Join(dir, "gdx.socket"), -1, false, strings.NewReader(""), &bytes.Buffer{})
		require.ErrorContains(t, err, "select one with --socket")
	})
}

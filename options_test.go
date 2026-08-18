package gdx

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOptions(t *testing.T) {
	t.Run("zero value applies stdlib default duration", func(t *testing.T) {
		cfg := Options().apply()
		require.Equal(t, 10*time.Second, cfg.defaultDuration)
	})

	t.Run("WithDefaultDuration overrides the default", func(t *testing.T) {
		cfg := Options().WithDefaultDuration(90 * time.Second).apply()
		require.Equal(t, 90*time.Second, cfg.defaultDuration)
	})

	t.Run("FromEnv falls back to 10s when unset", func(t *testing.T) {
		t.Setenv("DIAGX_DEFAULT_DURATION", "")
		cfg := Options().FromEnv().apply()
		require.Equal(t, 10*time.Second, cfg.defaultDuration)
	})

	t.Run("FromEnv reads DIAGX_DEFAULT_DURATION", func(t *testing.T) {
		t.Setenv("DIAGX_DEFAULT_DURATION", "45s")
		cfg := Options().FromEnv().apply()
		require.Equal(t, 45*time.Second, cfg.defaultDuration)
	})

	t.Run("chained calls compose in order, last write wins", func(t *testing.T) {
		t.Setenv("DIAGX_DEFAULT_DURATION", "45s")
		cfg := Options().WithDefaultDuration(5 * time.Second).FromEnv().apply()
		require.Equal(t, 45*time.Second, cfg.defaultDuration)

		cfg = Options().FromEnv().WithDefaultDuration(5 * time.Second).apply()
		require.Equal(t, 5*time.Second, cfg.defaultDuration)
	})
}

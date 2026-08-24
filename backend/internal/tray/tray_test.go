package tray

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTray(t *testing.T) {
	opts := Options{
		Version:      "1.0.0",
		DashboardURL: "http://localhost:20180",
		Port:         20180,
		Logger:       slog.Default(),
	}

	tray := New(opts)
	require.NotNil(t, tray)
	assert.Equal(t, "1.0.0", tray.opts.Version)
	assert.Equal(t, "http://localhost:20180", tray.opts.DashboardURL)
	assert.Equal(t, 20180, tray.opts.Port)
}

func TestDefaultIconEmbedded(t *testing.T) {
	assert.NotEmpty(t, DefaultIcon, "Default embedded tray icon should not be empty")
	// Verify PNG header magic bytes \x89PNG\r\n\x1a\n
	require.True(t, len(DefaultIcon) >= 8)
	assert.Equal(t, byte(0x89), DefaultIcon[0])
	assert.Equal(t, byte('P'), DefaultIcon[1])
	assert.Equal(t, byte('N'), DefaultIcon[2])
	assert.Equal(t, byte('G'), DefaultIcon[3])
}

func TestTrayStopIdempotent(t *testing.T) {
	tray := New(Options{
		Version:      "1.0.0",
		DashboardURL: "http://localhost:20180",
		Port:         20180,
	})
	// Stop should be safe to call multiple times even if tray is nil
	assert.NotPanics(t, func() {
		tray.Stop()
		tray.Stop()
	})
}

func TestTrayCannotAttachAfterEarlyStop(t *testing.T) {
	tray := New(Options{})
	tray.Stop()

	assert.False(t, tray.attach(nil), "a tray stopped during startup must not enter its event loop")
}

func TestTrayAttachAndStopAreRaceSafe(t *testing.T) {
	for range 100 {
		tray := New(Options{})
		started := make(chan struct{})
		done := make(chan struct{})

		go func() {
			close(started)
			tray.attach(nil)
			close(done)
		}()

		<-started
		tray.Stop()
		<-done
	}
}

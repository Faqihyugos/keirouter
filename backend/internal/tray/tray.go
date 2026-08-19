package tray

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"runtime"
	"sync"

	"github.com/gogpu/systray"
)

//go:embed icon.png
var DefaultIcon []byte

// Options configures the system tray.
type Options struct {
	Version      string
	DashboardURL string
	Port         int
	OnOpenURL    func(url string)
	OnQuit       func()
	Logger       *slog.Logger
}

// Tray manages the KeiRouter system tray lifecycle.
type Tray struct {
	opts    Options
	tray    *systray.SystemTray
	mu      sync.Mutex
	stopped bool
}

// New creates a new system tray instance with the given options.
func New(opts Options) *Tray {
	if opts.Logger == nil {
		opts.Logger = slog.Default()
	}
	return &Tray{
		opts: opts,
	}
}

// Run initializes and displays the system tray icon, builds the context menu,
// and blocks running the platform event loop. Callers should call this on the
// main OS thread (after runtime.LockOSThread()).
func (t *Tray) Run(ctx context.Context) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	initPlatform()

	st := systray.New()
	t.tray = st

	// Configure Icon
	if len(DefaultIcon) > 0 {
		st.SetIcon(DefaultIcon)
	}

	tooltip := fmt.Sprintf("KeiRouter v%s (%s)", t.opts.Version, t.opts.DashboardURL)
	st.SetTooltip(tooltip)

	// Build Context Menu
	menu := systray.NewMenu()

	// Title / Header item (disabled)
	titleItem := menu.Add(fmt.Sprintf("KeiRouter v%s", t.opts.Version), nil)
	if titleItem != nil {
		titleItem.SetDisabled(true)
	}

	statusItem := menu.Add(fmt.Sprintf("Status: Running (Port %d)", t.opts.Port), nil)
	if statusItem != nil {
		statusItem.SetDisabled(true)
	}

	menu.AddSeparator()

	// Open Dashboard
	menu.Add("Open Dashboard", func() {
		if t.opts.OnOpenURL != nil {
			t.opts.OnOpenURL(t.opts.DashboardURL)
		}
	})

	// Documentation
	menu.Add("Documentation", func() {
		if t.opts.OnOpenURL != nil {
			t.opts.OnOpenURL("https://github.com/mydisha/keirouter")
		}
	})

	menu.AddSeparator()

	// Quit
	menu.Add("Quit KeiRouter", func() {
		t.Stop()
		if t.opts.OnQuit != nil {
			t.opts.OnQuit()
		}
	})

	st.SetMenu(menu)

	// Left click on tray icon opens the dashboard
	st.OnClick(func() {
		if t.opts.OnOpenURL != nil {
			t.opts.OnOpenURL(t.opts.DashboardURL)
		}
	})

	st.Show()

	// Best-effort startup notification
	st.ShowNotification("KeiRouter", fmt.Sprintf("Running in background on %s", t.opts.DashboardURL))

	// Listen for ctx cancellation in background to tear down tray loop
	go func() {
		<-ctx.Done()
		t.Stop()
	}()

	t.opts.Logger.Info("system tray initialized", "url", t.opts.DashboardURL)

	// Run message loop (blocks until t.Stop() / st.Remove())
	if err := st.Run(); err != nil {
		t.opts.Logger.Warn("tray event loop returned error", "err", err)
		return err
	}
	return nil
}

// Stop removes the tray icon and exits the event loop safely.
func (t *Tray) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.stopped {
		return
	}
	t.stopped = true
	if t.tray != nil {
		t.tray.Remove()
	}
}

package cli

import (
	"testing"
	"time"

	"github.com/zaydek/retro/pkg/expect"
)

func TestWatch(t *testing.T) {
	var cmd WatchCommand

	cmd = parseWatchArguments()
	expect.Expect(t, cmd, WatchCommand{
		Cached:    false,
		Poll:      250 * time.Millisecond,
		Port:      8000,
		SourceMap: false,
		Paths:     []string{"pages"},
	})

	cmd = parseWatchArguments("src", "components", "pages")
	expect.Expect(t, cmd, WatchCommand{
		Cached:    false,
		Poll:      250 * time.Millisecond,
		Port:      8000,
		SourceMap: false,
		Paths:     []string{"src", "components", "pages"},
	})

	cmd = parseWatchArguments("--cached")
	expect.Expect(t, cmd, WatchCommand{
		Cached:    true,
		Poll:      250 * time.Millisecond,
		Port:      8000,
		SourceMap: false,
		Paths:     []string{"pages"},
	})

	cmd = parseWatchArguments("--poll=500ms")
	expect.Expect(t, cmd, WatchCommand{
		Cached:    false,
		Poll:      500 * time.Millisecond,
		Port:      8000,
		SourceMap: false,
		Paths:     []string{"pages"},
	})

	cmd = parseWatchArguments("--port=8080")
	expect.Expect(t, cmd, WatchCommand{
		Cached:    false,
		Poll:      250 * time.Millisecond,
		Port:      8080,
		SourceMap: false,
		Paths:     []string{"pages"},
	})

	cmd = parseWatchArguments("--source-map")
	expect.Expect(t, cmd, WatchCommand{
		Cached:    false,
		Poll:      250 * time.Millisecond,
		Port:      8000,
		SourceMap: true,
		Paths:     []string{"pages"},
	})

	cmd = parseWatchArguments("--cached", "--poll=500ms", "--port=8080", "--source-map", "src", "components", "pages")
	expect.Expect(t, cmd, WatchCommand{
		Cached:    true,
		Poll:      500 * time.Millisecond,
		Port:      8080,
		SourceMap: true,
		Paths:     []string{"src", "components", "pages"},
	})
}

func TestBuild(t *testing.T) {
	var cmd BuildCommand

	cmd = parseBuildArguments()
	expect.Expect(t, cmd, BuildCommand{
		Cached:    false,
		SourceMap: false,
	})

	cmd = parseBuildArguments("--cached")
	expect.Expect(t, cmd, BuildCommand{
		Cached:    true,
		SourceMap: false,
	})

	cmd = parseBuildArguments("--source-map")
	expect.Expect(t, cmd, BuildCommand{
		Cached:    false,
		SourceMap: true,
	})

	cmd = parseBuildArguments("--cached", "--source-map")
	expect.Expect(t, cmd, BuildCommand{
		Cached:    true,
		SourceMap: true,
	})
}

func TestServe(t *testing.T) {
	var cmd ServeCommand

	cmd = parseServeArguments()
	expect.Expect(t, cmd, ServeCommand{Port: 8000})

	cmd = parseServeArguments("--port=8080")
	expect.Expect(t, cmd, ServeCommand{Port: 8080})
}
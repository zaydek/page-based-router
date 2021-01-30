package cli

import "time"

// CreateCommand describes retro create flags and arguments.
type CreateCommand struct {
	Template string

	// Arguments
	Directory string
}

// WatchCommand describes retro watch flags and arguments.
type WatchCommand struct {
	Cached    bool
	Poll      time.Duration
	Port      int
	SourceMap bool

	// Arguments
	Paths []string
}

// BuildCommand describes retro build flags and arguments.
type BuildCommand struct {
	Cached    bool
	SourceMap bool
}

// ServeCommand describes retro serve flags and arguments.
type ServeCommand struct {
	Port int
}
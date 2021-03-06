package run

import (
	"testing"

	"github.com/zaydek/retro/pkg/expect"
)

func TestRun(t *testing.T) {
	stdout, err := Cmd("echo", "Hello, world!")
	if err != nil {
		t.Error(err)
	}
	expect.DeepEqual(t, string(stdout), "Hello, world!\n")
}

func TestRunError(t *testing.T) {
	_, err := Cmd("bad")
	if err == nil {
		t.Error("expected an error")
	}
	expect.DeepEqual(t, err.Error(), "exec: \"bad\": executable file not found in $PATH")
}

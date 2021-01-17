package main

import (
	_ "embed"

	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/zaydek/retro/static"
)

// Retro is a namespace for commands.
type Retro struct {
	stdout io.Writer // In writer
	stderr io.Writer // Out writer
}

func newRetro(stdout, stderr io.Writer) Retro {
	return Retro{stdout: stdout, stderr: stderr}
}

//go:embed usage
var usageMessage string

func (r Retro) help() {
	fmt.Fprintln(r.stdout, usageMessage)
}

func (r Retro) unknown(cmd string) {
	fmt.Fprintln(r.stderr, fmt.Sprintf("unknown command %s; try retro help", cmd))
}

func (r Retro) init(root string) {
	for _, asset := range static.Assets {
		dir := path.Join(root, path.Dir(asset.Path))
		if dir != "." {
			if err := os.MkdirAll(path.Join(root, dir), 0755); err != nil {
				log.Fatal("an unexpected error occurred; %w", err)
			}
		}
		if err := ioutil.WriteFile(path.Join(root, asset.Path), []byte(asset.Contents), 0644); err != nil {
			log.Fatal("an unexpected error occurred; %w", err)
		}
	}
}

func (r Retro) watch() {
	// ...
}

func (r Retro) build() {
	// ...
}

func (r Retro) serve() {
	// ...
}
package main

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type copyInfo struct {
	source string
	target string
}

func copyDir(from_dir, to_dir string, excludes []string) error {
	// Sweep for source and dest
	var infos []copyInfo
	err := filepath.WalkDir(from_dir, func(source string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		for _, exclude := range excludes {
			if source == exclude {
				return nil
			}
		}
		info := copyInfo{
			source: source,
			target: filepath.Join(to_dir, source),
		}
		infos = append(infos, info)
		return nil
	})
	if err != nil {
		return err
	}

	// Copy sources to targets
	for _, info := range infos {
		if dir := filepath.Dir(info.target); dir != "." {
			if err := os.MkdirAll(dir, PERM_DIR); err != nil {
				return err
			}
		}
		source, err := os.Open(info.source)
		if err != nil {
			return err
		}
		target, err := os.Create(info.target)
		if err != nil {
			return err
		}
		if _, err := io.Copy(target, source); err != nil {
			return err
		}
		source.Close()
		target.Close()
	}
	return nil
}

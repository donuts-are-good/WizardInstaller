package main

import (
	"embed"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//go:embed Application/*
var data embed.FS

func dumpTo(destination string) error {
	paths, err := walkFS(data)
	if err != nil {
		return fmt.Errorf("file enumeration failed: %w", err)
	}

	err = os.MkdirAll(destination, 0o755)
	if err != nil {
		return fmt.Errorf("failed to create %q: %w", destination, err)
	}

	for _, file := range paths {
		err = dumpBytes(data, file, filepath.Join(destination, file))
		if err != nil {
			return fmt.Errorf("failed to write %q to %q: %w", file, destination, err)
		}
	}

	return nil
}

func walkFS(f embed.FS) ([]string, error) {
	paths := []string{"Application"}
	for i := 0; i < len(paths); i++ {
		file, err := f.Open(paths[i])
		if err != nil {
			return nil, fmt.Errorf("found file does not exist? Wat?")
		}

		fi, err := file.Stat()
		if err != nil {
			return nil, fmt.Errorf("failed to stat file? %w", err)
		}

		if !fi.IsDir() {
			continue
		}

		entries, err := f.ReadDir(paths[i])
		if err != nil {
			return nil, fmt.Errorf("failed to list content of blob: %w", err)
		}

		for _, entry := range entries {
			paths = append(paths, path.Join(paths[i], entry.Name()))
		}
	}

	paths = paths[1:]
	for i, path := range paths {
		paths[i] = strings.TrimPrefix(path, "Application/")
	}

	return paths, nil
}

func dumpBytes(container embed.FS, inpath, outpath string) error {
	file, err := container.Open("Application/" + inpath)
	if err != nil {
		return fmt.Errorf("found file does not exist? Wat?")
	}

	fi, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file? %w", err)
	}

	if fi.IsDir() {
		return os.MkdirAll(outpath, 0o755)
	}

	fh, err := os.OpenFile(outpath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open destination file: %w", err)
	}
	defer fh.Close()

	_, err = io.Copy(fh, file)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

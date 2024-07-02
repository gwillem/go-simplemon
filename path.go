package simplemon

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func ageDaysOfNewestFile(glob string) (time.Duration, error) {
	files, err := filepath.Glob(glob)
	if err != nil {
		return 0, err
	}
	if len(files) == 0 {
		return 0, fmt.Errorf("no files found at %s", glob)
	}

	var newest time.Time
	for _, f := range files {
		info, err := os.Stat(f)
		if err != nil {
			return 0, err
		}

		if info.ModTime().After(newest) {
			newest = info.ModTime()
		}
	}

	return time.Since(newest), nil
}

func allDirsUnder(root string) (dirs []string) {
	err := filepath.Walk(root, func(path string, info os.FileInfo, e error) error {
		if e != nil && errors.Is(e, os.ErrNotExist) {
			return nil
		}
		if e != nil {
			fmt.Println(e)
			return e
		}
		if info.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})
	if err != nil {
		fmt.Println("outer:", err)
	}
	return dirs
}

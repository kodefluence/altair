package usecase

import (
	"io/fs"
	"path"
	"strconv"
	"strings"
)

// maxSourceVersion scans FS at sourcePath for golang-migrate style filenames
// (e.g. `4_create_table.up.sql`) and returns the highest leading integer.
// Returns 0 if the directory is empty or unreadable — callers treat that as
// "no target version", which matches golang-migrate's ErrNilVersion semantics.
func maxSourceVersion(fsys fs.FS, sourcePath string) uint {
	entries, err := fs.ReadDir(fsys, path.Clean(sourcePath))
	if err != nil {
		return 0
	}
	var max uint
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		underscore := strings.IndexByte(name, '_')
		if underscore <= 0 {
			continue
		}
		n, err := strconv.ParseUint(name[:underscore], 10, 64)
		if err != nil {
			continue
		}
		if uint(n) > max {
			max = uint(n)
		}
	}
	return max
}

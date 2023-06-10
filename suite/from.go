// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package suite

import (
	"os"
	"path/filepath"

	"github.com/jaypipes/gdt-core/scenario"
	"github.com/samber/lo"
)

var (
	validFileExts = []string{".yaml", ".yml"}
)

// FromDir reads the supplied directory path and returns a Suite representing
// the suite of test cases in that directory.
func FromDir(dirPath string) (*Suite, error) {
	// List YAML files in the directory and parse each into a testable unit
	s := New(WithPath(dirPath))

	if err := filepath.Walk(
		dirPath,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			suffix := filepath.Ext(path)
			if !lo.Contains(validFileExts, suffix) {
				return nil
			}
			f, err := os.Open(path)

			if err != nil {
				return err
			}
			defer f.Close()

			tc, err := scenario.FromReader(f, path)
			if err != nil {
				return err
			}
			s.Append(tc)
			return nil
		},
	); err != nil {
		return nil, err
	}
	return s, nil
}

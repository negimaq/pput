package convert

import (
	"path/filepath"
	"strings"
)

func getFilenameWithoutExtension(name string) string {
	return strings.TrimSuffix(name, filepath.Ext(name))
}

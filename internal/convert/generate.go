package convert

import (
	"fmt"
	"os"
)

type generateJPGImages struct {
	Converter
}

func (c generateJPGImages) run() error {
	// Create output dir
	if err := os.MkdirAll(c.OutputDirPath); err != nil {
		return err
	}

	// Convert images to JPG files
	// imagemagick must be installed
	args := []string{
		"-format jpg",
		"-quality 80%",
		"-define jpeg:extent=800kb",
		"-strip",
		"-interlace Plane",
		"-path",
		c.OutputDirPath,
		filepath.Join(c.InputDirPath, "*.*"),
	}
	if err := exec.Command("mogrify", args...).Run(); err != nil {
		return err
	}

	files, err := os.ReadDir(c.InputDirPath)
	if err != nil {
		return err
	}

	// Rename JPG files
	for i, v := range files {
		oldPath := filepath.Join(c.OutputDirPath, v)
		newPath := filepath.Join(c.OutputDirPath, v+"."+fmt.Sprintf("%05d.jpg", i+1))

		if err := os.Rename(oldPath, newPath); err != nil {
			return err
		}
	}

	return nil
}

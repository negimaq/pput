package convert

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type generateJPGImages struct {
	Converter
}

func (c generateJPGImages) run() error {
	// Create output dir
	if err := os.MkdirAll(c.OutputDirPath, 0755); err != nil {
		return err
	}

	// Convert images to JPG files
	// imagemagick must be installed
	args := []string{
		"-format", "jpg",
		"-quality", "80%",
		"-define", "jpeg:extent=800kb",
		"-strip",
		"-interlace", "Plane",
		"-path", c.OutputDirPath,
		filepath.Join(c.InputDirPath, "*.*"),
	}
	out, err := exec.Command("mogrify", args...).CombinedOutput()
	if err != nil {
		log.Print(string(out))
		return err
	}

	files, err := os.ReadDir(c.InputDirPath)
	if err != nil {
		return err
	}

	// Rename JPG files
	for i, v := range files {
		oldPath := filepath.Join(c.OutputDirPath, getFilenameWithoutExtension(v.Name())+".jpg")
		newPath := filepath.Join(c.OutputDirPath, fmt.Sprintf("%05d.jpg", i+1))

		if err := os.Rename(oldPath, newPath); err != nil {
			return err
		}
	}

	return nil
}

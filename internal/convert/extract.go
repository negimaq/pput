package convert

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type extractPDFImages struct {
	Converter
}

func (c extractPDFImages) run() error {
	files, err := os.ReadDir(c.InputDirPath)
	if err != nil {
		return err
	}

	path := filepath.Join(c.InputDirPath, files[0].Name())

	if err := os.Chdir(c.InputDirPath); err != nil {
		return err
	}

	// Extract images from PDF file
	// poppler-utils must be installed
	args := []string{
		"-all",
		path,
		"image",
	}
	out, err := exec.Command("pdfimages", args...).CombinedOutput()
	if err != nil {
		log.Print(string(out))
		return err
	}

	// Remove PDF file
	if err := os.Remove(path); err != nil {
		return err
	}

	return nil
}

package convert

import (
	"filepath"
	"os"
	"os/exec"
)

type extractPDFImages struct {
	Converter
}

func (c extractPDFImages) run() error {
	files, err := os.ReadDir(c.InputDirPath)
	if err != nil {
		return err
	}

	if err := os.Chdir(c.InputDirPath); err != nil {
		return err
	}

	// Extract images from PDF file
	// poppler-utils must be installed
	args := []string{
		"-all",
		filepath.Join(c.InputDirPath, files[0].Name()),
		"image",
	}
	if err := exec.Command("pdfimages", args...).Run(); err != nil {
		return err
	}

	// Remove PDF file
	if err := os.Remove(path); err != nil {
		return err
	}

	return nil
}

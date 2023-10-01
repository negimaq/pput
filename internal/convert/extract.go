package convert

import (
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
)

type extractPDFImages struct {
	Converter
}

func (c extractPDFImages) run() error {
	slog.Info("start runner", "type", "extractPDFImages", "path", c.InputDirPath)

	entries, err := os.ReadDir(c.InputDirPath)
	if err != nil {
		return err
	}

	path := filepath.Join(c.InputDirPath, entries[0].Name())

	if err := os.Chdir(c.InputDirPath); err != nil {
		return err
	}

	// Extract images from PDF file
	// poppler-utils must be installed
	slog.Debug("extracting images from PDF file", "path", path)
	args := []string{
		"-all",
		path,
		"image",
	}
	out, err := exec.Command("pdfimages", args...).CombinedOutput()
	if err != nil {
		slog.Error("failed to extract images from PDF file", "path", path, "msg", string(out))
		return err
	}

	// Remove PDF file
	if err := os.Remove(path); err != nil {
		return err
	}

	slog.Info("successfully terminate runner", "type", "extractPDFImages", "path", c.InputDirPath)

	return nil
}

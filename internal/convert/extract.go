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

// Extract images from PDF file
// poppler-utils must be installed
func (c extractPDFImages) run() error {
	slog.Debug("start runner", "type", "extractPDFImages", "path", c.InputDirPath)

	entries, err := os.ReadDir(c.InputDirPath)
	if err != nil {
		return err
	}

	p := filepath.Join(c.InputDirPath, entries[0].Name())

	if err := os.Chdir(c.InputDirPath); err != nil {
		return err
	}

	slog.Info("extracting images from PDF file", "path", p)
	args := []string{
		"-all",
		p,
		"image",
	}
	out, err := exec.Command("pdfimages", args...).CombinedOutput()
	if err != nil {
		slog.Error("failed to extract images from PDF file", "path", p, "msg", string(out))
		return err
	}
	slog.Info("images extracted from PDF file", "path", p)

	// Remove PDF file
	if err := os.Remove(p); err != nil {
		return err
	}
	slog.Debug("PDF file removed", "path", p)

	slog.Debug("successfully terminate runner", "type", "extractPDFImages", "path", c.InputDirPath)

	return nil
}

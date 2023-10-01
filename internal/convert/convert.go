// Package converts image files and extracts image files from PDF files.
package convert

import (
	"log/slog"
	"os"
	"path/filepath"
)

type runner interface {
	run() error
}

type Converter struct {
	InputDirPath  string
	OutputDirPath string
	Concurrency   int
}

func (c Converter) Run() error {
	slog.Info("start converter", "inputDirPath", c.InputDirPath, "outputDirPath", c.OutputDirPath, "concurrency", c.Concurrency)

	files, err := os.ReadDir(c.InputDirPath)
	if err != nil {
		return err
	}

	if len(files) == 1 && filepath.Ext(files[0].Name()) == ".pdf" {
		extractRunner := &extractPDFImages{
			Converter: c,
		}
		if err := extractRunner.run(); err != nil {
			return err
		}
	}

	deleteRunner := &deleteWhiteImages{
		Converter: c,
	}
	if err := deleteRunner.run(); err != nil {
		return err
	}

	generateRunner := &generateJPGImages{
		Converter: c,
	}
	if err := generateRunner.run(); err != nil {
		return err
	}

	slog.Info("successfully terminate converter", "inputDirPath", c.InputDirPath, "outputDirPath", c.OutputDirPath, "concurrency", c.Concurrency)

	return nil
}

// Package converts image files and extracts image files from PDF files.
package convert

import (
	"errors"
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
	Mode          string
	IsRename      bool
}

func (c Converter) Run() error {
	slog.Debug("start converter", "inputDirPath", c.InputDirPath, "outputDirPath", c.OutputDirPath, "concurrency", c.Concurrency, "mode", c.Mode, "isRename", c.IsRename)

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

	switch c.Mode {
	case "jpg":
		generateRunner := &generateJPGImages{
			Converter: c,
		}
		if err := generateRunner.run(); err != nil {
			return err
		}
	default:
		return errors.New("specified mode does not exist")
	}

	slog.Debug("successfully terminate converter", "inputDirPath", c.InputDirPath, "outputDirPath", c.OutputDirPath, "concurrency", c.Concurrency, "mode", c.Mode, "isRename", c.IsRename)

	return nil
}

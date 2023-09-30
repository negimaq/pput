// Package converts image files and extracts image files from PDF files.
package convert

import (
	"os"
	"path/filepath"
)

type runner interface {
	run() error
}

type Converter struct {
	InputDirPath  string
	OutputDirPath string
}

func (c Converter) Run() error {
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

	generateRunner := &generateJPGImages{
		Converter: c,
	}
	if err := generateRunner.run(); err != nil {
		return err
	}

	return nil
}

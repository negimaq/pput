// Package uploads image files and video files to cloud storage.
package upload

import (
	"errors"
	"log/slog"
	"sync"
)

type runner interface {
	run() error
}

type Uploader struct {
	DirPath       string
	ParentDirName string
	ChildDirName  string
	Mode          string
	Root          string
	User          string
	Password      string
	Mutex         sync.Mutex
}

func (u Uploader) Run() error {
	slog.Debug("start uploader", "dirPath", u.DirPath, "mode", u.Mode)

	switch u.Mode {
	case "nextcloud":
		uploadRunner := &uploadToNextcloud{
			Uploader: u,
		}
		if err := uploadRunner.run(); err != nil {
			return err
		}
	default:
		return errors.New("specified mode does not exist")
	}

	slog.Debug("successfully terminate uploader", "dirPath", u.DirPath, "mode", u.Mode)

	return nil
}

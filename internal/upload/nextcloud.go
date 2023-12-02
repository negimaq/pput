package upload

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/studio-b12/gowebdav"
)

type uploadToNextcloud struct {
	Uploader
}

// Check if the file name is included in the file list
func checkFileExist(fileList []os.FileInfo, targetFileName string) bool {
	for _, f := range fileList {
		if f.Name() == targetFileName {
			return true
		}
	}

	return false
}

func (u uploadToNextcloud) run() error {
	slog.Debug("start runner", "type", "uploadToNextcloud", "path", u.DirPath, "parentDirName", u.ParentDirName, "childDirName", u.ChildDirName)

	targetDirPath := filepath.Join(u.ParentDirName, u.ChildDirName)

	c := gowebdav.NewClient(u.Root, u.User, u.Password)
	if err := c.Connect(); err != nil {
		return err
	}

	if err := c.MkdirAll(targetDirPath, 0644); err != nil {
		return err
	}
	fis, err := c.ReadDir(targetDirPath)
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(u.DirPath)
	if err != nil {
		return err
	}

	for _, e := range entries {
		p := filepath.Join(u.DirPath, e.Name())
		if e.IsDir() {
			slog.Warn("the entry is directory (skip)", "path", p)
			continue
		}

		if checkFileExist(fis, e.Name()) {
			slog.Debug("the entry already uploaded (skip)", "path", p)
			continue
		}

		f, err := os.Open(p)
		if err != nil {
			return err
		}
		defer f.Close()

		tp := filepath.Join(targetDirPath, e.Name())
		if err := c.WriteStream(tp, f, 0644); err != nil {
			return err
		}
		slog.Info("file uploaded to nextcloud", "path", tp)
	}

	slog.Debug("successfully terminate runner", "type", "uploadToNextcloud", "path", u.DirPath, "parentDirName", u.ParentDirName, "childDirName", u.ChildDirName)

	return nil
}

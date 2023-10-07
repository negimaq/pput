package convert

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

type deleteWhiteImages struct {
	Converter
}

// Delete white images
// imagemagick must be installed
func (c deleteWhiteImages) run() error {
	slog.Debug("start runner", "type", "deleteWhiteImages", "path", c.InputDirPath, "concurrency", c.Concurrency)

	entries, err := os.ReadDir(c.InputDirPath)
	if err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(context.Background())
	sem := make(chan struct{}, c.Concurrency)

	for _, e := range entries {
		if !e.IsDir() {
			sem <- struct{}{}
			inputFilePath := filepath.Join(c.InputDirPath, e.Name())

			eg.Go(func() error {
				defer func() { <-sem }()
				select {
				case <-ctx.Done():
					slog.Warn("white image deletion canceled", "path", inputFilePath)
					return nil
				default:
					slog.Debug("checking white image", "path", inputFilePath)
					args := []string{
						"-format", "%[fx:255*mean]",
						inputFilePath,
					}
					out, err := exec.Command("identify", args...).Output()
					if err != nil {
						slog.Error("failed to check white image", "path", inputFilePath)
						return err
					}
					mean := string(out)
					if mean == "255" {
						if err := os.Remove(inputFilePath); err != nil {
							slog.Error("failed to delete white image", "path", inputFilePath)
							return err
						}
						slog.Info("white image deleted", "path", inputFilePath, "mean", mean)
					}
					return nil
				}
			})
		}
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	slog.Debug("successfully terminate runner", "type", "deleteWhiteImages", "path", c.InputDirPath, "concurrency", c.Concurrency)

	return nil
}

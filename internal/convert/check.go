package convert

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

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
			inputFilePath := filepath.Join(c.InputDirPath, e.Name())
			if e.Name() == "Thumbs.db" {
				if err := os.Remove(inputFilePath); err != nil {
					return err
				}
				slog.Info("unnecessary file deleted", "path", inputFilePath)
				continue
			}

			sem <- struct{}{}
			eg.Go(func() error {
				defer func() { <-sem }()
				select {
				case <-ctx.Done():
					slog.Warn("white image deletion canceled", "path", inputFilePath)
					return nil
				default:
					args := []string{
						"-format", "%[fx:255*mean]",
						inputFilePath,
					}
					out, err := exec.Command("identify", args...).Output()
					if err != nil {
						slog.Error("failed to check white image", "path", inputFilePath)
						return err
					}
					mean, err := strconv.ParseFloat(string(out), 64)
					if err != nil {
						return err
					}
					if 254.9 <= mean {
						if err := os.Remove(inputFilePath); err != nil {
							slog.Error("failed to delete white image", "path", inputFilePath)
							return err
						}
						slog.Info("white image deleted", "path", inputFilePath, "mean", mean)
					} else {
						slog.Debug("white image checked", "path", inputFilePath, "mean", mean)
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

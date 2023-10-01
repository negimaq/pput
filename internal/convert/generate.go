package convert

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/sync/errgroup"
)

type generateJPGImages struct {
	Converter
}

func (c generateJPGImages) run() error {
	slog.Info("start runner", "type", "generateJPGImages", "inputDirPath", c.InputDirPath, "outputDirPath", c.OutputDirPath, "concurrency", c.Concurrency)

	// Create output dir
	if err := os.MkdirAll(c.OutputDirPath, 0755); err != nil {
		return err
	}

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
					slog.Warn("JPG image generation canceled", "path", inputFilePath)
					return nil
				default:
					// Check image size
					// imagemagick must be installed
					slog.Debug("checking image size", "path", inputFilePath)
					CheckArgs := []string{
						"-format", "%w,%h",
						inputFilePath,
					}
					out, err := exec.Command("identify", CheckArgs...).Output()
					if err != nil {
						slog.Error("failed to check white image", "path", inputFilePath)
						return err
					}
					size := strings.Split(string(out), ",")
					width, err := strconv.Atoi(size[0])
					if err != nil {
						slog.Error("failed to get width", "path", inputFilePath)
						return err
					}
					height, err := strconv.Atoi(size[1])
					if err != nil {
						slog.Error("failed to get height", "path", inputFilePath)
						return err
					}
					sizeArgs := make([]string, 0, 14)
					if width <= height {
						if 2000 < height {
							sizeArgs = append(sizeArgs, "-resize", "x2000")
						}
					} else {
						if 2000 < width {
							sizeArgs = append(sizeArgs, "-resize", "2000x")
						}
					}

					// Convert image to JPG file
					// imagemagick must be installed
					slog.Debug("generating JPG image", "path", inputFilePath)
					otherArgs := []string{
						"-format", "jpg",
						"-quality", "80%",
						"-define", "jpeg:extent=800kb",
						"-strip",
						"-interlace", "Plane",
						"-path", c.OutputDirPath,
						inputFilePath,
					}
					convertArgs := append(sizeArgs, otherArgs...)
					out, err = exec.Command("mogrify", convertArgs...).CombinedOutput()
					if err != nil {
						slog.Error("failed to generate JPG image", "path", inputFilePath, "msg", string(out))
						return err
					}
					return nil
				}
			})
		}
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	// Rename JPG files
	for i, v := range entries {
		oldPath := filepath.Join(c.OutputDirPath, getFilenameWithoutExtension(v.Name())+".jpg")
		newPath := filepath.Join(c.OutputDirPath, fmt.Sprintf("%05d.jpg", i+1))

		if err := os.Rename(oldPath, newPath); err != nil {
			return err
		}
	}

	slog.Info("successfully terminate runner", "type", "generateJPGImages", "inputDirPath", c.InputDirPath, "outputDirPath", c.OutputDirPath, "concurrency", c.Concurrency)

	return nil
}

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

// Check image size
// imagemagick must be installed
func getImageSize(filePath string) (int, int, error) {
	checkArgs := []string{
		"-format", "%w,%h",
		filePath,
	}
	out, err := exec.Command("identify", checkArgs...).Output()
	if err != nil {
		slog.Error("failed to check image size", "path", filePath)
		return 0, 0, err
	}
	size := strings.Split(string(out), ",")
	width, err := strconv.Atoi(size[0])
	if err != nil {
		slog.Error("failed to get width", "path", filePath)
		return 0, 0, err
	}
	height, err := strconv.Atoi(size[1])
	if err != nil {
		slog.Error("failed to get height", "path", filePath)
		return 0, 0, err
	}
	slog.Debug("image size retrieved", "path", filePath, "width", width, "height", height)

	return width, height, nil
}

type generateJPGImages struct {
	Converter
}

// Convert image to JPG file
// imagemagick must be installed
func (c generateJPGImages) run() error {
	slog.Debug("start runner", "type", "generateJPGImages", "inputDirPath", c.InputDirPath, "outputDirPath", c.OutputDirPath, "concurrency", c.Concurrency)

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
					width, height, err := getImageSize(inputFilePath)
					if err != nil {
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
					out, err := exec.Command("mogrify", convertArgs...).CombinedOutput()
					if err != nil {
						slog.Error("failed to generate JPG image", "path", inputFilePath, "msg", string(out))
						return err
					}
					slog.Debug("JPG image generated", "path", inputFilePath)

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
	slog.Debug("JPG files renamed", "path", c.OutputDirPath)

	slog.Debug("successfully terminate runner", "type", "generateJPGImages", "inputDirPath", c.InputDirPath, "outputDirPath", c.OutputDirPath, "concurrency", c.Concurrency)

	return nil
}

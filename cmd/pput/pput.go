package main

import (
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/negimaq/pput/internal/convert"
)

var (
	inputDir  = "/input"
	outputDir = "/output"
)

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func main() {
	numConverters, err := strconv.Atoi(getEnv("NUM_CONVERTERS", "4"))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	concurrency, err := strconv.Atoi(getEnv("CONCURRENCY", "16"))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	slog.Info("start pput", "numConverters", numConverters, "concurrency", concurrency)

	entries, err := os.ReadDir(inputDir)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	wg := &sync.WaitGroup{}
	sem := make(chan struct{}, numConverters)

	for _, e := range entries {
		if e.IsDir() {
			sem <- struct{}{}

			wg.Add(1)
			go func(dirName string) {
				defer wg.Done()
				defer func() { <-sem }()

				c := convert.Converter{
					InputDirPath:  filepath.Join(inputDir, dirName),
					OutputDirPath: filepath.Join(outputDir, dirName),
					Concurrency:   concurrency,
				}
				if err := c.Run(); err != nil {
					slog.Error(err.Error())
					os.Exit(1)
				}
			}(e.Name())
		}
	}
	wg.Wait()
}

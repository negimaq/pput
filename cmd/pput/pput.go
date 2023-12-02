package main

import (
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lmittmann/tint"

	"github.com/negimaq/pput/internal/convert"
	"github.com/negimaq/pput/internal/upload"
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
	isDebug, err := strconv.ParseBool(getEnv("DEBUG", "false"))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	var level slog.Level
	if isDebug {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stdout, &tint.Options{
			Level:      level,
			TimeFormat: time.RFC3339,
		}),
	))
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

	converterMode := getEnv("CONVERTER_MODE", "jpg")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	uploaderMode := getEnv("UPLOADER_MODE", "nextcloud")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	root := getEnv("ROOT", "")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	} else if root == "" {
		slog.Error("failed to get ROOT")
		os.Exit(1)
	}

	user := getEnv("USER", "")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	} else if user == "" {
		slog.Error("failed to get USER")
		os.Exit(1)
	}

	password := getEnv("PASSWORD", "")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	} else if password == "" {
		slog.Error("failed to get PASSWORD")
		os.Exit(1)
	}

	slog.Debug("start pput", "numConverters", numConverters, "concurrency", concurrency)

	parentEntries, err := os.ReadDir(inputDir)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	wg := &sync.WaitGroup{}
	sem := make(chan struct{}, numConverters)

	for _, pe := range parentEntries {
		if !pe.IsDir() {
			slog.Debug("the entry is not directory (skip)", "path", filepath.Join(inputDir, pe.Name()))
			continue
		}
		childEntries, err := os.ReadDir(filepath.Join(inputDir, pe.Name()))
		if err != nil {
			slog.Warn(err.Error())
			continue
		}

		for _, ce := range childEntries {
			if !ce.IsDir() {
				slog.Info("the entry is not directory (skip)", "path", filepath.Join(inputDir, pe.Name(), ce.Name()))
				continue
			}

			sem <- struct{}{}
			wg.Add(1)
			go func(parentDirName, childDirName string, isRename bool) {
				defer wg.Done()
				defer func() { <-sem }()

				c := convert.Converter{
					InputDirPath:  filepath.Join(inputDir, parentDirName, childDirName),
					OutputDirPath: filepath.Join(outputDir, parentDirName, childDirName),
					Concurrency:   concurrency,
					Mode:          converterMode,
					IsRename:      isRename,
				}
				if err := c.Run(); err != nil {
					slog.Warn(err.Error())
					return
				}

				u := upload.Uploader{
					DirPath:       filepath.Join(outputDir, parentDirName, childDirName),
					ParentDirName: parentDirName,
					ChildDirName:  childDirName,
					Mode:          uploaderMode,
					Root:          root,
					User:          user,
					Password:      password,
				}
				if err := u.Run(); err != nil {
					slog.Warn(err.Error())
					return
				}
			}(pe.Name(), ce.Name(), !strings.HasPrefix(ce.Name(), "_"))
		}
	}
	wg.Wait()
}

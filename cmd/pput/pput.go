package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/negimaq/pput/internal/convert"
)

var (
	concurrency int
	inputDir    = "/input"
	outputDir   = "/output"
)

func init() {
	flag.IntVar(&concurrency, "c", 4, "concurrency for converting and uploading (default: 4)")
}

func main() {
	flag.Parse()

	entries, err := os.ReadDir(inputDir)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	wg := &sync.WaitGroup{}
	sem := make(chan struct{}, concurrency)

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
				}
				if err := c.Run(); err != nil {
					log.Fatal(err)
				}
			}(e.Name())
		}
	}
	wg.Wait()
}

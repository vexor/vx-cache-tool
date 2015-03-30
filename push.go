package main

import (
	"github.com/codegangsta/cli"
	"log"
	"os"
	"path/filepath"
)

func doPush(c *cli.Context) {
	log.Println("pushing: starting")

	// isCacheChanged()
	// pushUrl := c.Args()[0]

	log.Println("pushing: finishing")
}

func isCacheChanged() bool {
	if fileExists(cfg.FetchTar) {
		return true
	} else {
		withEachRegularFile(
			func(path string, mtime int64) {
				log.Println("===>>>", path, mtime)
			},
		)
		return false
	}
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else {
		return false
	}
}

func withEachRegularFile(f func(path string, mtime int64)) {
	for path, mtime := range mtimes {
		filepath.Walk(
			path,
			func(file string, finfo os.FileInfo, err error) error {
				if finfo.Mode().IsRegular() {
					f(file, mtime)
				}
				return nil
			},
		)
	}
}

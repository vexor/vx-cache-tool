package main

import (
	"github.com/codegangsta/cli"
	"log"
	"os"
	"path/filepath"
)

func doPush(c *cli.Context) {
	log.Println("pushing: starting")

	if isCacheChanged() {
		log.Println("changes detected, packing new archive")

		saveMd5Sums()

		args := []string{cfg.Md5File}
		for dir := range mtimes {
			args = append(args, dir)
		}
		tar(nil, "c", cfg.PushTar, args...)
	} else {
		log.Println("nothing changed, not updating cache")
	}
	// pushUrl := c.Args()[0]

	log.Println("pushing: finishing")
}

func isCacheChanged() (res bool) {
	res = false
	if !fileExists(cfg.FetchTar) {
		res = true
	} else {
		withEachRegularFile(
			func(path string, mtime int64) {
				if isFileUnchanged(path, mtime) {
					return
				} else {
					log.Println(path, "was modified")
					res = true
				}
			},
		)
	}
	return
}

func isFileUnchanged(file string, mtime int64) bool {
	if isMtimeUnchanged(file, mtime) {
		return true
	} else {
		if md5sum, ok := getMd5Sums()[file]; ok {
			md5, _ := fileMd5(file)
			return md5sum == md5
		} else {
			return false
		}
	}
}

func isMtimeUnchanged(file string, mtime int64) bool {
	fi, _ := os.Stat(file)
	return fi.ModTime().Unix() <= mtime
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

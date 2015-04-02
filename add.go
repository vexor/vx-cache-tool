package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"path/filepath"
	"time"
)

// Implements 'add' subcommand
func doAdd(c *cli.Context) {
	paths := c.Args()

	for _, path := range paths {
		path, err := filepath.Abs(path)
		if err != nil {
			continue
		}

		fmt.Printf("adding %s to cache\n", path)
		os.MkdirAll(path, 0755)
		mtimes[path] = time.Now().Unix()

		warner := func() {
			fmt.Println(path, "is not yet cached")
		}
		tar(warner, "x", cfg.FetchTar, path)
	}

	mtimes.store()
}

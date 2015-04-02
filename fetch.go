package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"os/exec"
)

func doFetch(c *cli.Context) {
	fetchAccessUrls := c.Args()
	isFetched := false

	for _, fetchAccessUrl := range fetchAccessUrls {
		storageUrl, err := getStorageUrl(fetchAccessUrl)
		if err == nil {
			err = fetchCacheArchive(storageUrl)
			if err == nil {
				isFetched = true
			}
		}
	}

	if isFetched {
		fmt.Println("cache archive fetched")
		fmt.Println("extracting checksums")
		warner := func() {
			fmt.Println("checksums aren't yet calculated, skipping")
		}
		tar(warner, "x", cfg.FetchTar, cfg.Md5File)
	} else {
		fmt.Println("could not download cache")
		os.Remove(cfg.FetchTar)
	}
}

func fetchCacheArchive(storageUrl string) (err error) {
	args := []string{
		"-m", "30",
		"-L",
		"--tcp-nodelay",
		"-f",
		"-s",
		storageUrl,
		"-o", cfg.FetchTar,
	}
	cmd := exec.Command("curl", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("FAILED: %s => %s", cmd.Args, out)
	}
	return
}

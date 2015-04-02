package main

import (
	"github.com/codegangsta/cli"
	"log"
	"os"
	"os/exec"
)

func doFetch(c *cli.Context) {
	log.Println("fetching: started")

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
		log.Println("cache archive fetched")
		log.Println("extracting checksums")
		warner := func() {
			log.Println("checksums aren't yet calculated, skipping")
		}
		tar(warner, "x", cfg.FetchTar, cfg.Md5File)
	} else {
		log.Printf("could not download cache")
		os.Remove(cfg.FetchTar)
	}

	log.Println("fetching: finished")
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
		log.Printf("FAILED: %s => %s", cmd.Args, out)
	}
	return
}

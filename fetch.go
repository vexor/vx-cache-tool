package main

import (
	"github.com/codegangsta/cli"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func doFetch(c *cli.Context) {
	log.Println("fetching: started")

	fetchUrls := c.Args()
	isFetched := false

	for _, fetchUrl := range fetchUrls {
		storageUrl, err := getStorageUrl(fetchUrl)
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

func getStorageUrl(fetchUrl string) (string, error) {
	log.Println("requesting cache archive location")
	res, err := http.Get(fetchUrl)
	if err != nil {
		log.Println(err)
		return "", err
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Println(err)
		return "", err
	}
	storageUrl := string(body)
	log.Println("received cache location at", storageUrl)
	return storageUrl, nil
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

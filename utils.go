package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func check(e error, blah ...string) {
	if e != nil {
		fmt.Fprintln(os.Stderr, "FAIL: =>", e, "=>", blah)
		os.Exit(1)
	}
}

func getStorageUrl(accessUrl string) (string, error) {
	res, err := http.Get(accessUrl)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return "", err
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return "", err
	}
	storageUrl := string(body)
	return storageUrl, nil
}

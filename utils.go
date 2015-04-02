package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func getStorageUrl(accessUrl string) (string, error) {
	log.Println("requesting cache archive location")
	res, err := http.Get(accessUrl)
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

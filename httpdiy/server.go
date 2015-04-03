package main

import (
	"github.com/codegangsta/cli"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

func runServer(c *cli.Context) {
	log.Println("server: starting at localhost:4242; press CTRL+C to exit")

	r := mux.NewRouter()

	r.HandleFunc("/urls/r123456", fetchUrlHandler)
	r.HandleFunc("/storage/fetch.tgz", fetchFileHandler)
	r.HandleFunc("/urls/w123456", pushUrlHandler)
	r.HandleFunc("/storage{params:.*}", pushFileHandler)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe("localhost:4242", nil))

	log.Println("server: finishing")
}

func fetchUrlHandler(w http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, req.URL)
	io.WriteString(w, "http://localhost:4242/storage/fetch.tgz")
}

func fetchFileHandler(w http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, req.URL)
	http.ServeFile(w, req, "files/fetch.tgz")
}

func pushUrlHandler(w http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, req.URL)
	io.WriteString(w, "http://localhost:4242/storage")
}

func pushFileHandler(w http.ResponseWriter, req *http.Request) {
	log.Println(
		req.Method,
		req.Proto,
		req.RequestURI,
		req.Header,
		req.Form,
		req.MultipartForm,
	)
}

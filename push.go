package main

import (
	"bytes"
	b64 "encoding/base64"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
)

type ChunkInfo struct {
	Chunk  string
	Digest string
}

func doPush(c *cli.Context) {
	if isCacheChanged() {
		fmt.Println("changes detected, packing new archive")

		saveMd5Sums()

		args := []string{cfg.Md5File}
		for dir := range mtimes {
			args = append(args, dir)
		}
		tar(nil, "c", cfg.PushTar, args...)

		pushAccessUrl := c.Args()[0]
		storageUrl, err := getStorageUrl(pushAccessUrl)
		if err == nil {
			pushCacheArchive(storageUrl)
		} else {
			fmt.Println("failed to retrieve cache url")
		}
	} else {
		fmt.Println("nothing changed, not updating cache")
	}
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
					fmt.Println(path, "was modified")
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

func pushCacheArchive(uri string) {
	chunks := splitPushTar()

	args := []string{
		"-XPUT",
		"-H", "x-ms-blob-type: BlockBlob",
		"-s",
		"-S",
		"-m", "60",
		"-d", "",
		uri,
	}
	cmd := exec.Command("curl", args...)
	out, err := cmd.CombinedOutput()
	check(err, string(out))

	fmt.Print("uploading chunks ")
	res := false
	for _, chunk := range chunks {
		chunkUrl := fmt.Sprintf("%s&comp=block&blockid=%s", uri, url.QueryEscape(chunk.Digest))
		args := []string{
			"-s",
			"-S",
			"-m", "60",
			"-T",
			chunk.Chunk,
			chunkUrl,
		}
		cmd := exec.Command("curl", args...)
		_, err := cmd.CombinedOutput()
		if err == nil {
			fmt.Print(".")
			res = true
		} else {
			res = false
			break
		}
	}
	if res {
		fmt.Println(" OK")
	} else {
		fmt.Println(" FAIL")
		os.Exit(1)
	}

	commitChunks(chunks, uri)
}

func commitChunks(chunks []ChunkInfo, uri string) {
	var xmlBuffer bytes.Buffer
	chunksFile := filepath.Join(cfg.BackupDir, "chunks.xml")

	xmlBuffer.WriteString("<?xml version=\"1.0\" encoding=\"utf-8\"?>\n")
	xmlBuffer.WriteString("<BlockList>\n")
	for _, chunk := range chunks {
		xmlBuffer.WriteString(fmt.Sprintf("<Latest>%s</Latest>\n", chunk.Digest))
	}
	xmlBuffer.WriteString("</BlockList>\n")
	ioutil.WriteFile(chunksFile, xmlBuffer.Bytes(), 0644)

	fmt.Print("committing chunks")
	args := []string{
		"-s",
		"-S",
		"-m", "60",
		"-H", "x-ms-version: 2011-08-18",
		"-T",
		chunksFile,
		fmt.Sprintf("%s&comp=blocklist", uri),
	}
	cmd := exec.Command("curl", args...)
	_, err := cmd.CombinedOutput()
	if err == nil {
		fmt.Println(" OK")
	} else {
		fmt.Println(" FAIL")
	}
}

func splitPushTar() []ChunkInfo {
	res := []ChunkInfo{}
	args := []string{
		"-a", "8",
		"-b", "4m",
		cfg.PushTar,
		fmt.Sprintf("%s.", cfg.PushTar),
	}
	cmd := exec.Command("split", args...)

	wd, _ := os.Getwd()
	defer os.Chdir(wd)

	os.Chdir(cfg.BackupDir)
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("could not split archive to chunks")
		return res
	}

	chunks, _ := filepath.Glob(fmt.Sprintf("%s.*", cfg.PushTar))
	for _, chunk := range chunks {
		md5Digest, _ := fileMd5(chunk)
		encodedDigest := b64.StdEncoding.EncodeToString([]byte(md5Digest))

		res = append(res, ChunkInfo{
			Chunk:  chunk,
			Digest: encodedDigest,
		})
	}
	return res
}

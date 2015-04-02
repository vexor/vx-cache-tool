package main

import (
	"crypto/md5"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
)

var (
	md5sums Md5Sums
)

type Md5Sums map[string]string

func getMd5Sums() Md5Sums {
	if md5sums == nil {
		if fileExists(cfg.Md5File) {
			md5sums.restore()
		} else {
			md5sums = make(Md5Sums)
		}
	}
	return md5sums
}

func (sums *Md5Sums) restore() {
	*sums = make(Md5Sums)
	yml, _ := ioutil.ReadFile(cfg.Md5File)
	yaml.Unmarshal(yml, sums)
}

func (sums Md5Sums) store() {
	yml, _ := yaml.Marshal(sums)
	ioutil.WriteFile(cfg.Md5File, yml, 0644)
}

func fileMd5(filePath string) (string, error) {
	var result []byte
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Sprintf("%x", result), err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return fmt.Sprintf("%x", result), err
	}
	return fmt.Sprintf("%x", hash.Sum(result)), nil
}

func saveMd5Sums() {
	fmt.Println("calculating checksums")
	newMd5Sums := make(Md5Sums)

	withEachRegularFile(
		func(path string, mtime int64) {
			alreadyInMd5Sums := false
			if _, ok := md5sums[path]; ok {
				alreadyInMd5Sums = true
			}
			if isMtimeUnchanged(path, mtime) && alreadyInMd5Sums {
				newMd5Sums[path] = md5sums[path]
			} else {
				newMd5Sums[path], _ = fileMd5(path)
			}
		},
	)

	yml, _ := yaml.Marshal(newMd5Sums)
	ioutil.WriteFile(cfg.Md5File, yml, 0644)
}

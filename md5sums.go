package main

import (
	"crypto/md5"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
)

var (
	md5sums Md5Sums
)

type Md5Sums map[string][]byte

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

func fileMd5(filePath string) ([]byte, error) {
	var result []byte
	file, err := os.Open(filePath)
	if err != nil {
		return result, err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return result, err
	}
	return hash.Sum(result), nil
}

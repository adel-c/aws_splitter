package main

import (
	"os"
	"path/filepath"
)

type FilesList struct {
	filesMap       map[string]*os.File
	workDir        string
	truncateOnOpen bool
}

type FileHandler interface {
	getFile(s string) *os.File
}

func closeAllFiles(r FilesList) {
	for s := range r.filesMap {
		err := r.filesMap[s].Close()
		if err != nil {
			println(err)
		}
	}
}
func (r FilesList) getFile(s string) *os.File {
	var existingFile, ok = r.filesMap[s]
	// If the key exists
	if ok {
		return existingFile
	}

	var path = filepath.Join(r.workDir, s)
	var dir = filepath.Dir(path)
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		panic(err)
	}

	flag := os.O_WRONLY | os.O_CREATE | os.O_APPEND
	if r.truncateOnOpen {
		flag = os.O_TRUNC | os.O_CREATE | os.O_WRONLY
	}
	f, err := os.OpenFile(path, flag, 0600)
	if err != nil {
		panic(err)
	}
	r.filesMap[s] = f
	return f

}

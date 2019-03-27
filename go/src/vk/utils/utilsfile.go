package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func OpenFile(path string, fileFlags int, userPerms os.FileMode) (f *os.File, err error) {

	if err = FileDir(path); nil != err {
		return
	}

	f, err = os.OpenFile(path, fileFlags, userPerms)
	if err != nil {
		err = fmt.Errorf("error opening file: %v", err)
		return
	}

	return
}

func FileDir(full string) (err error) {

	permDir := os.FileMode(0700)

	dirpath := filepath.Dir(full)

	if err = os.MkdirAll(dirpath, permDir); nil != err {
		return
	}

	return
}

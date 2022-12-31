package helper

import (
	"os"
	"strings"
)

func CreateFile(path string) (*os.File, error) {
	if strings.ContainsRune(path, '/') {
		arr := strings.Split(path, "/")
		dir := strings.Join(arr[:len(arr)-1], "/")
		if err := EnsureDir(dir); err != nil {
			return nil, err
		}
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func EnsureDir(path string) error {
	if ok, err := PathExists(path); err != nil || !ok {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

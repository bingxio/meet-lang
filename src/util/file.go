package util

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func Suffix(name string) bool {
	suffix := path.Ext(name)

	return suffix == ".meet"
}

func Read(path string) ([]byte, error) {
	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if err != nil {
		panic("获取当前目录失败：" + err.Error())
	}

	return strings.Replace(dir, "\\", "/", -1)
}

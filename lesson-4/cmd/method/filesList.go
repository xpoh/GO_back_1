package main

import (
	"io/fs"
	"io/ioutil"
	"strings"
	"time"
)

type JsonFileInfo struct {
	Name    string      `json:"Name"`
	Size    int64       `json:"Size"`
	Mode    fs.FileMode `json:"Mode"`
	ModTime time.Time   `json:"ModTime"`
	IsDir   bool        `json:"IsDir"`
}

type FilesList struct {
	list []JsonFileInfo
	dir  string
}

func (fl FilesList) create(dir string) FilesList {
	return FilesList{
		list: nil,
		dir:  dir,
	}
}

func (fl *FilesList) refresh(ext string) error {
	files, err := ioutil.ReadDir(fl.dir)
	if err != nil {
		return err
	}

	fl.list = nil // TODO Возможна ли утечка памяти?
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ext) {
			fl.list = append(fl.list, JsonFileInfo{
				Name:    file.Name(),
				Size:    file.Size(),
				Mode:    file.Mode(),
				ModTime: file.ModTime(),
				IsDir:   file.IsDir(),
			})
		}
	}
	return err
}

package test

import (
	"os"
	"time"
)

type FileInfoMock struct {
	FName string
	Dir   bool
}

func (fi FileInfoMock) Name() string {
	return fi.FName
}

func (fi FileInfoMock) Size() int64 {
	return 0
}

func (fi FileInfoMock) Mode() os.FileMode {
	return os.FileMode(0)
}

func (fi FileInfoMock) ModTime() time.Time {
	return time.Now()
}

func (fi FileInfoMock) IsDir() bool {
	return fi.Dir
}

func (fi FileInfoMock) Sys() interface{} {
	return ""
}

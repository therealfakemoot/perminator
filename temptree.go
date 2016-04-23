package main

import (
	"io/ioutil"
)

func create_zeroed_tree(target string, depth int, width int, files int) string {
	tmpDir, _ := ioutil.TempDir(target, "massperms_")

	for i := 0; i <= files; i++ {
		ioutil.TempFile(tmpDir, "massperms_")
	}

	for i := 0; i <= width; i++ {
		tmpDir, _ := ioutil.TempDir(tmpDir, "massperms_")
		create_zeroed_tree(tmpDir, depth-1, width, files)
	}

	return tmpDir

}

func test() {
	create_zeroed_tree("", 5, 5, 2)
}

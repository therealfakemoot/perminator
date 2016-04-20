package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

func getCwd() (dir string) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func getPatterns() {

	if os.Getenv("MASSPERMS_CONF") {
		configPath := filepath.Join(os.Getenv("MASSPERMS_CONF"), ".massperms.rc")
	} else {
		configPath := filepath.Join(currentUser.HomeDir, ".massperms.rc")
		currentUser, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
	}

	conf, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(conf)
}

func main() {
	fmt.Println(getCwd())
}

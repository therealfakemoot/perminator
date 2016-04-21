package main

import (
	"flag"
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

func getConfigPath() (configPath string) {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	configPath = filepath.Join(currentUser.HomeDir, ".massperms.rc")
	return configPath

}

func LoadConfig(configPath string) (conf []byte) {

	conf, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	return conf
}

func buildConfigPath() {
}

func main() {

	configPathPtr := flag.String("config", getConfigPath(), "Path to your massperms patterns file.")
	targetDirPtr := flag.String("target", getCwd(), "Path to directory to apply patterns to.")

	flag.Parse()

	log.Print("CLI FLAG: config:", *configPathPtr)
	log.Print("CLI FLAG: target:", *targetDirPtr)

	var file_list []string
	file_list, err := filepath.Glob(*targetDirPtr)
	if err != nil {
		log.Fatal(err)
	}
	//log.Print("Matched files", file_list)
	for _, file := range file_list {
		fmt.Println("Matched File:", file)
	}
}

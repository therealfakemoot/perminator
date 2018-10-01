package main

import (
	"flag"
	"log"
	"os/user"
	"path"
	"path/filepath"

	p "github.com/therealfakemoot/perminator/src"
)

var (
	targetDir  string
	configPath string
)

func homeDir() string {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return currentUser.HomeDir
}

func main() {
	var err error
	log.Print("perminator start")

	flag.StringVar(&targetDir, "targetDir", homeDir(), "Target directory.")
	flag.StringVar(&configPath, "configPath", path.Join(homeDir(), ".perminator.conf"), "Config file location")

	flag.Parse()

	targetDir, err = filepath.Abs(targetDir)

	if err != nil {
		log.Fatalf("could not render absolute target directory: %s", err)
	}

	rs, err := p.LoadRules(configPath)

	if err != nil {
		log.Fatalf("error loading ruleset: %s", err)
	}

	log.Printf("loaded ruleset: %+v\n", rs)

	path, err := filepath.Abs(targetDir)
	if err != nil {
		log.Fatal(err)
	}

	err = filepath.Walk(path, p.Apply(rs, targetDir))
	if err != nil {
		log.Fatal(err)
	}

	log.Print("perminator exit")
}

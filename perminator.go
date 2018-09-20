package main

import (
	"flag"
	// "fmt"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

var (
	targetDir  string
	configPath string
	debugMode  bool
)

type Rule struct {
	Pattern string
	Type    string
	Perm    os.FileMode
}

type RuleSet []Rule

func Debug(v ...interface{}) {
	if debugMode {
		log.Print(v...)
	}
}

func Debugf(fmt string, v ...interface{}) {
	if debugMode {
		log.Printf(fmt, v...)
	}
}

func homeDir() string {
	currentUser, err := user.Current()
	if err != nil {
		log.Panic(err)
	}
	return currentUser.HomeDir
}

func loadRules(path string) RuleSet {
	var rs RuleSet

	return rs
}

func parseRule(s string) Rule {
	var r Rule

	return r
}

func match(p string, r Rule) bool {
	return false
}

func Apply(rules RuleSet) filepath.WalkFunc {
	f := func(path string, info os.FileInfo, err error) error {
		for _, r := range rules {
			if match(path, r) {
			}
		}
		return nil
	}

	return f
}

func main() {
	Debug("Perminator start.")

	flag.StringVar(&targetDir, "targetDir", homeDir(), "Target directory.")
	flag.StringVar(&configPath, "configPath", path.Join(homeDir(), ".perminator.rc"), "Config file location.")
	flag.BoolVar(&debugMode, "verbose", false, "Verbose logging mode.")

	flag.Parse()

	rs := loadRules(configPath)

	Debugf("Loaded ruleset: %+v\n", rs)

	err := filepath.Walk(targetDir, Apply(rs))
	if err != nil {
		log.Panic(err)
	}
	Debug("Perminator exit.")
}

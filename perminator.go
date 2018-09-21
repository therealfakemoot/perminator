package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	targetDir  string
	configPath string
	debugMode  bool
)

type Rule struct {
	Pattern string
	Type    string
	Mode    os.FileMode
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

	conf, err := ioutil.ReadFile(path)

	if err != nil {
		log.Panic(err)
	}

	for _, line := range strings.Split(string(conf), "\n") {
		rs = append(rs, parseRule(line))
	}

	return rs
}

func parseRule(s string) Rule {
	var (
		pattern string
		fstype  string
		rawMode string
	)

	fmt.Sscanf(s, "%s %1s%s", &pattern, &fstype, &rawMode)
	i, err := strconv.ParseUint(rawMode, 8, 32)

	if err != nil {
		log.Panic(err)
	}

	r := Rule{
		Pattern: pattern,
		Type:    fstype,
		Mode:    os.FileMode(uint32(i)),
	}

	return r
}

func match(p string, r Rule) bool {
	return false
}

func Apply(rules RuleSet) filepath.WalkFunc {
	f := func(path string, info os.FileInfo, err error) error {
		log.Println(path)
		for _, r := range rules {
			if match(path, r) {
				os.Chmod(path, r.Mode)
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

	path, err := filepath.Abs(targetDir)
	if err != nil {
		log.Panic(err)
	}

	err = filepath.Walk(path, Apply(rs))
	if err != nil {
		log.Panic(err)
	}
	Debug("Perminator exit.")
}

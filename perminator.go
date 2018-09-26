package main

import (
	"errors"
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
	ErrBadFileType = errors.New("invalid filetype in rule")
	ErrBadPerms    = errors.New("invalid file permissions in rule")
	ErrBadPattern  = errors.New("invalid file globbing pattern in rule")
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

func homeDir() string {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return currentUser.HomeDir
}

func loadRules(path string) (RuleSet, error) {
	var rs RuleSet

	conf, err := ioutil.ReadFile(path)

	if err != nil {
		log.Print("error loading config file")
		return rs, err
	}

	lines := strings.Split(string(conf), "\n")

	log.Printf("raw config: %q", string(conf))
	log.Printf("config lines: %+v", lines)
	for _, line := range strings.Split(string(conf), "\n") {
		if line == "" {
			continue
		}
		r, err := parseRule(line)
		if err != nil {
			log.Printf("bad rule encountered, %q: %s", line, err)
			continue
		}
		rs = append(rs, r)
	}

	return rs, nil
}

func parseRule(s string) (Rule, error) {
	var (
		pattern string
		fstype  string
		rawMode string
		r       Rule
	)

	_, err := fmt.Sscanf(s, "%s %1s%s", &pattern, &fstype, &rawMode)
	if err != nil {
		return r, err
	}

	switch fstype {
	case "d", "f", "a":
	default:
		return r, ErrBadFileType

	}

	i, err := strconv.ParseUint(rawMode, 8, 32)

	if err != nil {
		return r, ErrBadPerms
	}

	r = Rule{
		Pattern: pattern,
		Type:    fstype,
		Mode:    os.FileMode(uint32(i)),
	}

	return r, nil
}

func match(pattern, name string) (bool, error) {
	return filepath.Match(pattern, name)
}

func Apply(rules RuleSet) filepath.WalkFunc {
	f := func(fname string, info os.FileInfo, err error) error {
		log.Printf("walking over %s", fname)
		for _, r := range rules {
			pattern := path.Join(targetDir, r.Pattern)
			log.Printf("matching against pattern: %s", pattern)
			m, err := match(pattern, fname)
			if err != nil {
				return err
			}
			if m {
				log.Printf("updating permissions for %s: %s", fname, r.Mode)
				err := os.Chmod(fname, r.Mode)
				if err != nil {
					log.Printf("unable to modify %s: %s", fname, err)
					return err
				}
			}
		}
		return nil
	}

	return f
}

func main() {
	var err error
	log.Print("perminator start.")

	flag.StringVar(&targetDir, "targetDir", homeDir(), "Target directory.")
	flag.StringVar(&configPath, "configPath", path.Join(homeDir(), ".perminator.rc"), "Config file location.")
	flag.BoolVar(&debugMode, "verbose", false, "Verbose logging mode.")

	flag.Parse()

	targetDir, err = filepath.Abs(targetDir)

	if err != nil {
		log.Fatalf("could not render absolute target directory: %s", err)
	}

	rs, err := loadRules(configPath)

	if err != nil {
		log.Fatalf("error loading ruleset: %s", err)
	}

	log.Printf("loaded ruleset: %+v\n", rs)

	path, err := filepath.Abs(targetDir)
	if err != nil {
		log.Fatal(err)
	}

	err = filepath.Walk(path, Apply(rs))
	if err != nil {
		log.Fatal(err)
	}

	log.Print("perminator exit.")
}

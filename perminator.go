package main

import (
	"fmt"
	logrus "github.com/Sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

type rule struct {
	pattern     string
	fstype      string
	permissions os.FileMode
}

func (p rule) String() string {
	//out_strings := []string{p.pattern, p.fstype, string(p.permissions)}
	//return strings.Join(out_strings, ":")
	out_string := p.pattern + ":" + p.fstype + ":" + fmt.Sprintf("%d", p.permissions)
	return out_string
}

type ruleSet struct {
	rules []rule
}

func (p ruleSet) String() string {
	out_string := ""
	for _, rule := range p.rules {
		out_string += "|"
		out_string += rule.String()
	}
	out_string += "|"

	return out_string
}

type Perminator struct {
	rules     []rule
	targetDir string
}

func (p Perminator) Apply(path string, info os.FileInfo, err error) error {
	for i := len(p.rules) - 1; i >= 0; i-- {
		rule := p.rules[i]
		pattern := getCwd() + "/" + rule.pattern

		absPath, err := filepath.Abs(path)
		if err != nil {
			logger.Error(err)
		}

		logger.WithFields(logrus.Fields{
			"pattern":  pattern,
			"filePath": absPath,
		}).Debug("Matching.")
		match, _ := filepath.Match(pattern, absPath)

		if match == true {

			os.Chmod(path, os.FileMode(rule.permissions))
			logger.WithFields(logrus.Fields{
				"pattern":     pattern,
				"filePath":    absPath,
				"permissions": (int(rule.permissions.Perm())),
			}).Info("File modified.")

			break

		}
	}
	logger.Info("No rules applied.")
	return err
}

var logger = logrus.New()

func getCwd() (dir string) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Error(err)
	}
	return dir
}

func getConfigPath() (configPath string) {
	currentUser, err := user.Current()
	if err != nil {
		logger.Error(err)
	}
	configPath = filepath.Join(currentUser.HomeDir, ".perminator.rc")
	return configPath
}

func loadConfig(configPath string) (parsedConfig ruleSet, err error) {
	conf, err := ioutil.ReadFile(configPath)
	if err != nil {
		logger.Error(err)
	}

	rawConfLines := strings.Split(string(conf), "\n")
	confLines := rawConfLines[:len(rawConfLines)-1]
	rules := make([]rule, 0, len(confLines))

	for index, line := range confLines {
		var pattern string
		var fstype string
		var perms string
		var permissions os.FileMode

		fmt.Sscanf(line, "%s %1s%d", &pattern, &fstype, &perms)

		permsInt, _ := strconv.ParseUint(perms, 10, 32)
		logger.WithFields(logrus.Fields{
			"permsInt": permsInt,
		}).Info("Permissions integer.")
		if err != nil {
		}
		permissions = os.FileMode(int32(permsInt))

		switch fstype {
		case "a", "d", "f":
			_ = ""
		default:
			logger.WithFields(logrus.Fields{
				"raw":    line,
				"line":   index,
				"fstype": fstype,
			}).Info("Invalid filesystem type value.")
		}

		if err != nil || (permissions < 0 || permissions > 777) {
			logger.WithFields(logrus.Fields{
				"line":        index,
				"permsString": permissions,
			}).Info("Invalid file permissions value.")
		}
		rules = append(rules, rule{pattern, fstype, permissions})

		logger.WithFields(logrus.Fields{
			"rulePriority": index,
			"pattern":      pattern,
			"fstype":       fstype,
			"permissions":  permissions,
		}).Debug("Rule loaded.")
	}

	ruleSet := ruleSet{rules}

	logger.WithFields(logrus.Fields{
		"ruleCount":  len(ruleSet.rules),
		"configPath": configPath,
	}).Debug("Configuration file loaded.")

	logger.WithFields(logrus.Fields{
		"fullRules": ruleSet.String(),
	}).Debug("Full parsed ruleset.")

	return ruleSet, err

}

func main() {

	logger.Out = os.Stderr

	var (
		configPath = kingpin.Flag("config", "Configuration file path.").Short('c').Default(getConfigPath()).ExistingFile()
		targetDir  = kingpin.Flag("target", "Target directory.").Short('d').Default(getCwd()).ExistingDir()
		debugMode  = kingpin.Flag("debug", "Enable debugging output.").Bool()
	)

	kingpin.Parse()

	if *debugMode {
		logger.Level = logrus.DebugLevel
	} else {
		logger.Level = logrus.InfoLevel
	}

	logger.WithFields(logrus.Fields{
		"configPath": *configPath,
		"targetPath": *targetDir,
	}).Info("perminator begins")

	ruleSet, err := loadConfig(*configPath)
	if err != nil {
		logger.Info("Error parsing config file.")
	}
	P := Perminator{ruleSet.rules, *targetDir}

	filepath.Walk(*targetDir, P.Apply)

	logger.WithFields(logrus.Fields{
		"ruleSetLength": len(ruleSet.rules),
	}).Info("Ruleset processed.")
}

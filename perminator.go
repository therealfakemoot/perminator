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

type permsRule struct {
	pattern     string
	fstype      string
	permissions int
}

func (p permsRule) String() string {
	//out_strings := []string{p.pattern, p.fstype, string(p.permissions)}
	//return strings.Join(out_strings, ":")
	out_string := p.pattern + ":" + p.fstype + ":" + strconv.Itoa(p.permissions)
	return out_string
}

type permsRuleSet struct {
	rules []permsRule
}

func (p permsRuleSet) String() string {
	out_string := ""
	for _, rule := range p.rules {
		out_string += "|"
		out_string += rule.String()
	}
	out_string += "|"

	return out_string
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

func loadConfig(configPath string) permsRuleSet {
	conf, err := ioutil.ReadFile(configPath)
	if err != nil {
		logger.Error(err)
	}

	rawConfLines := strings.Split(string(conf), "\n")
	confLines := rawConfLines[:len(rawConfLines)-1]

	permsRules := make([]permsRule, 0, len(confLines))

	for index, line := range confLines {
		var pattern string
		var fstype string
		var permissions int

		fmt.Sscanf(line, "%s %1s%d", &pattern, &fstype, &permissions)

		permsRules = append(permsRules, permsRule{pattern, fstype, permissions})

		logger.WithFields(logrus.Fields{
			"rulePriority": index,
			"pattern":      pattern,
			"fstype":       fstype,
			"permissions":  permissions,
		}).Debug("Rule loaded.")
	}

	ruleSet := permsRuleSet{permsRules}

	logger.WithFields(logrus.Fields{
		"ruleCount":  len(ruleSet.rules),
		"configPath": configPath,
	}).Debug("Configuration file loaded.")

	logger.WithFields(logrus.Fields{
		"fullRules": ruleSet.String(),
	}).Debug("Full parsed ruleset.")

	return ruleSet

}

func buildWalkFunc(ruleSet permsRuleSet) func(string, os.FileInfo, error) error {
	return func(path string, info os.FileInfo, err error) error {
		for _, rule := range ruleSet.rules {
			match, _ := filepath.Match(rule.pattern, path)
			if match {
				if logger.Level != logrus.DebugLevel {
					os.Chmod(path, os.FileMode(rule.permissions))
				}
				logger.WithFields(logrus.Fields{
					"matchPath":   path,
					"rule":        rule,
					"permissions": strconv.Itoa(rule.permissions),
				}).Info("File modified.")
			}
		}
		return err
	}
}

func applyRules(targetDir string, ruleSet permsRuleSet) {
	walkfunc := buildWalkFunc(ruleSet)
	filepath.Walk(targetDir, walkfunc)
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

	ruleSet := loadConfig(*configPath)

	applyRules(*targetDir, ruleSet)

	logger.WithFields(logrus.Fields{
		"ruleSetLength": len(ruleSet.rules),
	}).Info("Ruleset processed.")
}

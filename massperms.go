package main

import (
	"flag"
	"fmt"
	logrus "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
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
	out_string := p.pattern + ":" + p.fstype + ":" + string(p.permissions)
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
	configPath = filepath.Join(currentUser.HomeDir, ".massperms.rc")
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

func globWalk(pattern string, dir string) []string {

	return []string{"whatever"}

}

func main() {

	logger.Out = os.Stderr
	//logger.SetFormatter(&logrus.TextFormatter{})

	configPathPtr := flag.String("config", getConfigPath(), "Path to your massperms patterns file.")
	targetDirPtr := flag.String("target", getCwd(), "Path to directory to apply patterns to.")
	debugLevelPtr := flag.Bool("debug", false, "Debug mode.")

	flag.Parse()

	if *debugLevelPtr {
		logger.Level = logrus.DebugLevel
	} else {
		logger.Level = logrus.InfoLevel
	}

	logger.WithFields(logrus.Fields{
		"configPath": *configPathPtr,
		"targetPath": *targetDirPtr,
	}).Info("massperms begins")

	ruleSet := loadConfig(*configPathPtr)
	logger.WithFields(logrus.Fields{
		"ruleSetLength": len(ruleSet.rules),
	}).Info("Ruleset processed.")
}

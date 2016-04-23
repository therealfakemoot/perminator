package main

import (
	"flag"
	logrus "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

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

func LoadConfig(configPath string) (conf []byte) {
	conf, err := ioutil.ReadFile(configPath)
	if err != nil {
		logger.Error(err)
	}

	return conf
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
	}).Debug("massperms begins")

	var file_list []string
	file_list, err := filepath.Glob(*targetDirPtr)
	if err != nil {
		logger.Error(err)
	}

	logger.WithFields(logrus.Fields{
		"fileCount": len(file_list),
	}).Info("Totle files affected.")
}

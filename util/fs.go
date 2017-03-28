package util

import (
	"errors"
	"os"
	"path"
)

func IsFileExists(path string) bool {
	stat, err := os.Stat(path)
	if err == nil {
		if stat.Mode() & os.ModeType == 0 {
			return true
		}
		return false
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func GetConfigPath(fileName string) (configFile string, err error) {
	configFile = fileName
	var binDir string

	if IsFileExists(configFile) {
		return
	}

	if binDir, err = os.Getwd(); err == nil {
		configFile = path.Join(binDir, fileName)
	}

	if IsFileExists(configFile) {
		return
	}

	if executable, err := os.Executable(); err == nil {
		binDir = path.Dir(executable)
		configFile = path.Join(binDir, fileName)
	}

	if IsFileExists(configFile) {
		return
	}

	return "", errors.New("File not found")
}


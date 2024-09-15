package setup

import (
	"fmt"
	"os"
	"path/filepath"
)

func createDirIfNotExists(dirPath string) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			errMsg := fmt.Sprintln("Can't create config dir at \"%s\". Output is: \n%s", dirPath, err.Error())
			panic(errMsg)
		}
	}
}
func createFileIfNotExists(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err := os.WriteFile(filePath, []byte(""), 0755)
		if err != nil {
			errMsg := fmt.Sprintln("Can't create config file at \"%s\". Output is: \n%s", filePath, err.Error())
			panic(errMsg)
		}
	}
}

type PathFromHome struct {
	ConfigFile string
	CssFile    string
}

func NewPathFromHome(configFile string, cssFile string) *PathFromHome {
	homeDir := getHomeDir()
	return &PathFromHome{
		ConfigFile: filepath.Join(homeDir, configFile),
		CssFile:    filepath.Join(homeDir, cssFile),
	}
}

func (f *PathFromHome) CreateFilesAndDirs() {
	configFileDir := filepath.Dir(f.ConfigFile)
	configCssDir := filepath.Dir(f.ConfigFile)

	createDirIfNotExists(configFileDir)
	createDirIfNotExists(configCssDir)

	createFileIfNotExists(f.ConfigFile)
	createFileIfNotExists(f.CssFile)
}

func getHomeDir() string {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		errMsg := fmt.Sprintln("Can't read env $HOME")
		panic(errMsg)
	}

	return homeDir
}

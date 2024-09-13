package setup

import (
	"fmt"
	"os"
)

func createIfNotExists(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err := os.WriteFile(filePath, []byte(""), 0755)
		if err != nil {
			errMsg := fmt.Sprintln("Can't create config file at \"%s\". Output is: \n%s", filePath, err.Error())
			panic(errMsg)
		}
	}
}

func CreateConfigFiles(filePaths ...string) {
	for _, filePath := range filePaths {
		createIfNotExists(filePath)
	}
}

func GetHomeDir() string {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		errMsg := fmt.Sprintln("Can't read env $HOME")
		panic(errMsg)
	}

    return homeDir
}

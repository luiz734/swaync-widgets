package setup

import (
	"fmt"
	"os"
	"path/filepath"
)

func createDirIfNotExists(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return fmt.Errorf("can't create config dir at %s: %w", dirPath, err)
		}
	}
	return nil
}
func createFileIfNotExists(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := os.WriteFile(filePath, []byte(""), 0755); err != nil {
			return fmt.Errorf("can't create file at %s: %w", filePath, err)
		}
	}
	return nil
}

type PathFromHome struct {
	ConfigFile string
	CssFile    string
}

func NewPathFromHome(configFile string, cssFile string) (*PathFromHome, error) {
	homeDir, err := getHomeDir()
	if err != nil {
        return nil, fmt.Errorf("can't get home dir: %w", err)
	}
	return &PathFromHome{
		ConfigFile: filepath.Join(homeDir, configFile),
		CssFile:    filepath.Join(homeDir, cssFile),
	}, nil
}

func (f *PathFromHome) CreateFilesAndDirs() error {
	configFileDir := filepath.Dir(f.ConfigFile)
	configCssDir := filepath.Dir(f.ConfigFile)

	if err := createDirIfNotExists(configFileDir); err != nil {
		return fmt.Errorf("can't create config dir: %w", err)
	}
	if err := createDirIfNotExists(configCssDir); err != nil {
		return fmt.Errorf("can't create css dir: %w", err)
	}
	if err := createFileIfNotExists(f.ConfigFile); err != nil {
		return fmt.Errorf("can't create config file: %w", err)
	}

	if err := createFileIfNotExists(f.CssFile); err != nil {
		return fmt.Errorf("can't create css file: %w", err)
	}

	return nil
}

func getHomeDir() (string, error) {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		return "", fmt.Errorf("can't read env $HOME")
	}
	return homeDir, nil
}

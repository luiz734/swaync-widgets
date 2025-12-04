package config

import (
	"os"
	"path/filepath"
	"strings"
)

func ExpandPath(path string) string {
	if path == "" {
		return path
	}

	if strings.HasPrefix(path, "~/") {
		if home := os.Getenv("HOME"); home != "" {
			path = filepath.Join(home, path[2:])
		}
	} else if path == "~" {
		if home := os.Getenv("HOME"); home != "" {
			path = home
		}
	}

	return os.ExpandEnv(path)
}

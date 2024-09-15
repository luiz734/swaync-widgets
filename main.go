package main

import (
	"fmt"
	"log"
	"os"
	"swaync-widgets/app"
	"swaync-widgets/cli"
	"swaync-widgets/config"
	"swaync-widgets/setup"
	"github.com/pelletier/go-toml/v2"
)

func main() {
	configFilePath := setup.NewPathFromHome(".config/swaync-widgets/config.toml", ".config/swaync/widgets.css")
	// configFilePath := setup.NewPathFromHome("config.toml", "widgets.css")
	configFilePath.CreateFilesAndDirs()

	// Try to read the config file
	file, err := os.ReadFile(configFilePath.ConfigFile)
	if err != nil {
		log.Fatalf("Can't read config file at \"%s\". Output is: \n%s", configFilePath.ConfigFile, err.Error())
	}

	// Config file should be avaliabe now
	var cfg config.Config
	err = toml.Unmarshal(file, &cfg)
	if err != nil {
		log.Fatalf("Error parsing config file at \"%s\". Output is: \n%s", configFilePath.ConfigFile, err.Error())
	}

	args := cli.ParseCliArgs()

	targetWidget, err := args.TargetWidget(&cfg)
	if err != nil {
        errMsg := fmt.Errorf("can't get target: %w\n", err)
		fmt.Fprintf(os.Stderr, errMsg.Error())
        os.Exit(1)
	}

	app.ToggleWidget(*targetWidget)
	app.UpdateConfigFiles(cfg)
	app.ReloadConfigFiles(cfg)
}

package main

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"log"
	"os"
	"swaync-widgets/app"
	"swaync-widgets/cli"
	"swaync-widgets/config"
	"swaync-widgets/setup"
)

func main() {
	configFilePath := setup.NewPathFromHome(".config/swaync-widgets/config.toml", ".config/swaync/widgets.css")
	// configFilePath := setup.NewPathFromHome("config.toml", "widgets.css")
	configFilePath.CreateFilesAndDirs()
	// Try to read the config file
	file, err := os.ReadFile(configFilePath.ConfigFile)
	if err != nil {
        errMsg := fmt.Errorf("can't read config file: %w", err)
		fmt.Fprintf(os.Stderr, errMsg.Error())
        os.Exit(1)
	}
	// Config file should be avaliabe now
	var cfg config.Config
	err = toml.Unmarshal(file, &cfg)
	if err != nil {
        errMsg := fmt.Errorf("can't parse config file: %w", err)
		fmt.Fprintf(os.Stderr, errMsg.Error())
        os.Exit(1)
	}
    // Parse arguments
	args, err := cli.ParseCliArgs()
    if err != nil {
        errMsg := fmt.Errorf("invalid option(s): %w\n", err)
		fmt.Fprintf(os.Stderr, errMsg.Error())
        os.Exit(1)
    }
    // Arguments are good. Find a target
	targetWidget, err := args.TargetWidget(&cfg)
	if err != nil {
		errMsg := fmt.Errorf("can't get target: %w\n", err)
		fmt.Fprintf(os.Stderr, errMsg.Error())
		os.Exit(1)
	}
    // targetWidget = nil ==> just a reload on the config files
	if targetWidget != nil {
		app.RunToggleWidget(*targetWidget)
	}
    // Update both config and widgets css files, them reload swaync
	app.WriteConfigAndCss(cfg)
	app.RunReloadConfigFiles(cfg)
}

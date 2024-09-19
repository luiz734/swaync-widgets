package main

import (
	"fmt"
	"os"
	"swaync-widgets/app"
	"swaync-widgets/cli"
	"swaync-widgets/config"
	"swaync-widgets/setup"

	"github.com/pelletier/go-toml/v2"
)

func main() {
	configFilePath, err := setup.NewPathFromHome(".config/swaync-widgets/config.toml", ".config/swaync/widgets.css")
	if err != nil {
		errMsg := fmt.Errorf("can't setup config paths: %w", err)
		fmt.Fprintf(os.Stderr, errMsg.Error())
		os.Exit(1)
	}
	// configFilePath := setup.NewPathFromHome("config.toml", "widgets.css")
	if err := configFilePath.CreateFilesAndDirs(); err != nil {
		errMsg := fmt.Errorf("can't create files and dirs: %w", err)
		fmt.Fprintf(os.Stderr, errMsg.Error())
		os.Exit(1)
	}
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
        errMsg := fmt.Errorf("can't parse config file: %w: ", err)
		fmt.Fprintf(os.Stderr, errMsg.Error())
		os.Exit(1)
	}
	// Validate structs
    if err := config.ValidateConfig(cfg); err != nil {
		errMsg := fmt.Errorf("can't validate config file: %w", err)
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
		if err := app.RunToggleWidget(*targetWidget); err != nil {
			errMsg := fmt.Errorf("can't get target: %w\n", err)
			fmt.Fprintf(os.Stderr, errMsg.Error())
			os.Exit(1)
		}
	}
	// Update both config and widgets css files, them reload swaync
	if err := app.WriteConfigAndCss(cfg); err != nil {
        errMsg := fmt.Errorf("can't write config and css: %w: ", err)
		fmt.Fprintf(os.Stderr, errMsg.Error())
		os.Exit(1)
	}
	app.RunReloadConfigFiles(cfg)
}

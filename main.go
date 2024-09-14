package main

import (
	"fmt"
	"log"
	"os"

	// "swaync-widgets/app"
	"swaync-widgets/app"
	"swaync-widgets/cli"
	"swaync-widgets/config"
	"swaync-widgets/setup"

	"github.com/pelletier/go-toml/v2"
)

func main() {
	args := cli.ParseCliArgs()
	_ = args
	homeDir := setup.GetHomeDir()
	configFile := homeDir + "/.config/swaync-widgets/config.toml"
	// configFile = "config.toml"
	cssFile := homeDir + "/.config/swaync/widgets.css"
	setup.CreateConfigFiles(configFile, cssFile)

	// Try to read or create a config file
	file, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Can't read config file at \"%s\". Output is: \n%s", configFile, err.Error())
	}

	// Config file should be avaliabe now
	var cfg config.Config
	err = toml.Unmarshal(file, &cfg)
	if err != nil {
		log.Fatalf("Error parsing config file at \"%s\". Output is: \n%s", configFile, err.Error())
	}

	var targetWidgetName = args.Widget

	if targetWidgetName != "" {
		var targetWidget *config.WidgetConfig
		var options []string

		for _, w := range cfg.Widgets {
			options = append(options, w.Desc)
			if w.Desc == targetWidgetName {
				targetWidget = &w
			}
		}

		if targetWidget == nil {
            fmt.Printf("Invalid widget \"%s\". Options are \"%s\"\n", targetWidgetName, options)
            os.Exit(1)
		}
        app.ToggleWidget(*targetWidget)
	}

	app.UpdateConfigFiles(cfg)
	app.ReloadConfigFiles(cfg)
}

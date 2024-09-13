package main

import (
	"github.com/pelletier/go-toml/v2"
	"log"
	"os"
	"swaync-widgets/app"
	"swaync-widgets/cli"
	"swaync-widgets/config"
	"swaync-widgets/setup"
)

func main() {
	args := cli.ParseCliArgs()
	homeDir := setup.GetHomeDir()
	configFile := homeDir + "/.config/swaync-widgets/config.toml"
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

	switch args.Widget {
	case "mute":
		app.ToggleWidget(cfg.WidgetMute)
	case "vpn":
		app.ToggleWidget(cfg.WidgetVpn)
	case "wifi":
		app.ToggleWidget(cfg.WidgetWifi)
	case "bluetooth":
		app.ToggleWidget(cfg.WidgetBluetooth)
	default:
	}

	app.UpdateConfigFiles(cfg)
	app.ReloadConfigFiles(cfg)
}

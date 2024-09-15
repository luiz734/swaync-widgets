package app

import (
	"log"
	"os/exec"
	"swaync-widgets/config"
)

func RunToggleWidget(widgetConfig config.WidgetConfig) {
	stateOn := RunGetWidgetState(widgetConfig.CheckStatusCommand)
	var command string
	if stateOn {
		command = widgetConfig.TurnOffCommand
	} else {
		command = widgetConfig.TurnOnCommand
	}
	cmd := exec.Command("bash", "-c", command)
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("Can't execute command \"%s\". Output is: \n%s\nCheck your config for widget \"%s\"", command, err.Error(), widgetConfig.Desc)
	}
}

func RunReloadConfigFiles(cfg config.Config) {
	cmd := exec.Command("bash", "-c", cfg.SwayncReloadCommand)
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("Can't reload config with command \"%s\". Output is: \n%s", cfg.SwayncReloadCommand, err.Error())
	}
}

func RunGetWidgetState(command string) bool {
	cmd := exec.Command("bash", "-c", command)
	out, err := cmd.Output()
	StateOn := err == nil && len(out) != 0
	// println(StateOn, " ", " cmd: "+command)
	return StateOn
}

package app

import (
	"fmt"
	"os/exec"
	"swaync-widgets/config"
)

func RunToggleWidget(widgetConfig config.WidgetConfig) error {
	stateOn := RunGetWidgetState(widgetConfig.CheckStatusCommand)
	var command string
	if stateOn {
		command = widgetConfig.TurnOffCommand
	} else {
		command = widgetConfig.TurnOnCommand
	}
	cmd := exec.Command("bash", "-c", command)
	if _, err := cmd.Output(); err != nil {
		return fmt.Errorf("error running command \"%s\" %w", command, err)
	}

    return nil
}

func RunReloadConfigFiles(cfg config.Config) error {
	cmd := exec.Command("bash", "-c", cfg.SwayncReloadCommand)
	if _, err := cmd.Output(); err != nil {
		return fmt.Errorf("error running command \"%s\" %w", cfg.SwayncReloadCommand, err)
	}

	return nil
}

// This function is weird. The goal is to check if the state is on/off
// based on the output. Because we are running the command inside a bash
// process, it's not possible to tell if the error was intentional (the
// state is really down), or the command simply doesn't exists. In both
// ways, the error would be an ExitError, and we can't tell the
// difference. The config file should be double checked in the option
// "check_status_command". If a state is always down, could be it.
func RunGetWidgetState(command string) bool {
	cmd := exec.Command("bash", "-c", command)
	out, err := cmd.Output()

    // This will not work as explained above
	// if err != nil {
	//        var exitErr *exec.ExitError
	// 	if errors.As(err, &exitErr) {
	//            return false, nil
	//        }
	//        return false, fmt.Errorf("error running command \"%s\": %w", command, err)
	// }

	StateOn := err == nil && len(out) != 0
	return StateOn
}

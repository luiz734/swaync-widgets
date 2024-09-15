package cli

import (
	// "errors"
	"fmt"
	"os"
	"swaync-widgets/config"
)

type CliArgs struct {
	Widget string
}

func ShowUsage(msg string) {
	usageMsg := "swaync-widgets [wifi|bluetooth|mute|vpn]"
	fmt.Sprintln("%s. Usage is: \n%s", msg, usageMsg)
	os.Exit(1)
}

func ParseCliArgs() CliArgs {
	widget := ""
	// action := ""
	if len(os.Args) > 1 {
		widget = os.Args[1]
	}
	return CliArgs{
		Widget: widget,
	}
}

func (a *CliArgs) TargetWidget(cfg *config.Config) (*config.WidgetConfig, error) {

	if a.Widget == "" {
		return nil, fmt.Errorf("missing required argument <widget>")
	}

	var targetWidget *config.WidgetConfig
	var options []string

	for _, w := range cfg.Widgets {
		options = append(options, w.Desc)
		if w.Desc == a.Widget {
			targetWidget = &w
		}
	}

    var ErrInvalidWidget = fmt.Errorf("options are %s", options)
	if targetWidget == nil {
        return nil, fmt.Errorf("invalid widget \"%s\": %w", a.Widget, ErrInvalidWidget)
	}

    return targetWidget, nil
}

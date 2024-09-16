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

func ParseCliArgs() (*CliArgs, error) {
	widget := ""

	if len(os.Args) == 2 {
		widget = os.Args[1]
	} else if len(os.Args) > 2 {
        return nil, fmt.Errorf("too many arguments")

    }
	return &CliArgs{
		Widget: widget,
	}, nil
}

// Try to get a target to work on based on the cli arg
func (a *CliArgs) TargetWidget(cfg *config.Config) (*config.WidgetConfig, error) {
    // There is no match, but not a problem. The default option (just reload) will be used
	if a.Widget == "" {
		return nil, nil //fmt.Errorf("missing required argument <widget>")
	}
    // Finds a matching widget in the config file
	var targetWidget *config.WidgetConfig
	var options []string
	for _, w := range cfg.Widgets {
		options = append(options, w.Desc)
		if w.Desc == a.Widget {
			targetWidget = &w
		}
	}
    // Nothing to lookup on the config file
    if len(options) == 0 {
        return nil, fmt.Errorf("no options avaliable in config file")
    }
    // There are stuff to lookup, but none of them match
    var ErrInvalidWidget = fmt.Errorf("options are %s", options)
	if targetWidget == nil {
        return nil, fmt.Errorf("invalid widget \"%s\": %w", a.Widget, ErrInvalidWidget)
	}

    return targetWidget, nil
}

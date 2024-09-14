package cli

import (
	"fmt"
	"os"
)

type CliArgs struct {
	Widget string
	Action string
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
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
	// if len(os.Args) > 2 {
	// 	widget = os.Args[2]
	// }

	// if !contains([]string{"mute", "vpn", "wifi", "bluetooth", ""}, widget) {
	// 	ShowUsage("Invalid option " + widget)
	// }

	// if !contains([]string{"on", "off", "toggle", ""}, action) {
	// 	ShowUsage("Invalid option " + action)
	// }

	return CliArgs{
		Widget: widget,
		Action: "toggle",
	}
}

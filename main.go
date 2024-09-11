package main

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Styles struct {
	CssButton      string `toml:"css_button"`
	CssButtonHover string `toml:"css_button_hover"`
	CssLabel       string `toml:"css_label"`
	CssLabelHover  string `toml:"css_label_hover"`
}

type WidgetConfig struct {
	Desc               string `toml:"desc"`
	Index              string `toml:"index"`
	OffLabel           string `toml:"off_label"`
	OnLabel            string `toml:"on_label"`
	TurnOnCommand      string `toml:"turn_on_command"`
	TurnOffCommand     string `toml:"turn_off_command"`
	CheckStatusCommand string `toml:"check_status_command"`
}

type Config struct {
	SwayncCssWidgets       string       `toml:"swaync_css_widgets"`
	SwayncConfigFile       string       `toml:"swaync_config_file"`
	SwayncReloadCommand    string       `toml:"swaync_reload_command"`
	CSSPrepend             string       `toml:"css_prepend"`
	CSSButtonSelector      string       `toml:"css_button_selector"`
	CSSButtonHoverSelector string       `toml:"css_button_hover_selector"`
	CSSLabelSelector       string       `toml:"css_label_selector"`
	CSSLabelHoverSelector  string       `toml:"css_label_hover_selector"`
	StylesOn               Styles       `toml:"styles_on"`
	StylesOff              Styles       `toml:"styles_off"`
	WidgetVpn              WidgetConfig `toml:"vpn"`
	WidgetMute             WidgetConfig `toml:"mute"`
	WidgetWifi             WidgetConfig `toml:"wifi"`
	WidgetBluetooth        WidgetConfig `toml:"bluetooth"`
}

func GetWidgetCss(cfg Config, widgetConfig WidgetConfig) string {
	stateOn := GetWidgetState(widgetConfig.CheckStatusCommand)
	if stateOn {
		return GetOnCss(cfg, widgetConfig.Index, widgetConfig.Desc)
	} else {
		return GetOffCss(cfg, widgetConfig.Index, widgetConfig.Desc)
	}
}

func GetWidgetState(command string) bool {
	cmd := exec.Command("bash", "-c", command)
	out, err := cmd.Output()
	StateOn := err == nil && len(out) != 0
	// println(StateOn, " ", " cmd: "+command)
	return StateOn
}

func UpdateConfigFiles(cfg Config) {
	outputCss := cfg.CSSPrepend

	outputCss += GetWidgetCss(cfg, cfg.WidgetVpn)
	sedConfigFile(cfg, cfg.WidgetVpn)

	outputCss += GetWidgetCss(cfg, cfg.WidgetMute)
	sedConfigFile(cfg, cfg.WidgetMute)

	outputCss += GetWidgetCss(cfg, cfg.WidgetWifi)
	sedConfigFile(cfg, cfg.WidgetWifi)

	outputCss += GetWidgetCss(cfg, cfg.WidgetBluetooth)
	sedConfigFile(cfg, cfg.WidgetBluetooth)

	err := os.WriteFile(cfg.SwayncCssWidgets, []byte(outputCss), 0755)
	if err != nil {
		log.Fatalf("Can't write file at \"%s\". Output is: \n%s", cfg.SwayncReloadCommand, err.Error())
	}

}

func GetOnCss(cfg Config, index string, comment string) string {
	output := "/* widget " + comment + " */\n"
	output += cfg.CSSButtonSelector
	output += "{" + cfg.StylesOn.CssButton + "}\n"
	output += cfg.CSSButtonHoverSelector
	output += "{" + cfg.StylesOn.CssButtonHover + "}\n"
	output += cfg.CSSLabelSelector
	output += "{" + cfg.StylesOn.CssLabel + "}\n"
	output += cfg.CSSLabelHoverSelector
	output += "{" + cfg.StylesOn.CssLabelHover + "}\n"
	output = strings.Replace(output, "?", index, -1)
	return output
}

func sedConfigFile(cfg Config, widgetConfig WidgetConfig) {
	stateOn := GetWidgetState(widgetConfig.CheckStatusCommand)
	firstPart := widgetConfig.OnLabel
	secondPart := widgetConfig.OffLabel
	var sedCommand string
	if stateOn {
		firstPart = widgetConfig.OffLabel
		secondPart = widgetConfig.OnLabel
	}
	sedCommand = fmt.Sprintf(
		"sed -i 's/\"label\": \"%s\"/\"label\": \"%s\"/' \"%s\"",
		firstPart,
		secondPart,
		cfg.SwayncConfigFile)

	cmd := exec.Command("bash", "-c", sedCommand)
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("Can't execute command \"%s\". Output is: \n%s\nCheck your config for widget \"%s\"", sedCommand, err.Error(), widgetConfig.Desc)
	}
}

func GetOffCss(cfg Config, index string, comment string) string {
	output := "/* widget " + comment + " */\n"
	output += cfg.CSSButtonSelector
	output += "{" + cfg.StylesOff.CssButton + "}\n"
	output += cfg.CSSButtonHoverSelector
	output += "{" + cfg.StylesOff.CssButtonHover + "}\n"
	output += cfg.CSSLabelSelector
	output += "{" + cfg.StylesOff.CssLabel + "}\n"
	output += cfg.CSSLabelHoverSelector
	output += "{" + cfg.StylesOff.CssLabelHover + "}\n"
	output = strings.Replace(output, "?", index, -1)
	return output
}

type CliArgs struct {
	widget string
	action string
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
	log.Fatalf("%s. Usage is: \n%s", msg, usageMsg)
}

func ParseCliArgs() CliArgs {
	widget := ""
	action := ""
	if len(os.Args) > 1 {
		widget = os.Args[1]
	}
	if len(os.Args) > 2 {
		widget = os.Args[2]
	}

	if !contains([]string{"mute", "vpn", "wifi", "bluetooth", ""}, widget) {
		ShowUsage("Invalid option " + widget)
	}

	if !contains([]string{"on", "off", "toggle", ""}, action) {
		ShowUsage("Invalid option " + action)
	}

	return CliArgs{
		widget: widget,
		action: action,
	}
}

func ToggleWidget(widgetConfig WidgetConfig) {
	stateOn := GetWidgetState(widgetConfig.CheckStatusCommand)
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

func ReloadConfigFiles(cfg Config) {
	cmd := exec.Command("bash", "-c", cfg.SwayncReloadCommand)
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("Can't reload config with command \"%s\". Output is: \n%s", cfg.SwayncReloadCommand, err.Error())
	}
}

func createConfigFile(filePath string) {
	err := os.WriteFile(filePath, []byte("brand new config file"), 0755)
	if err != nil {
		log.Fatalf("Can't create config file at \"%s\". Output is: \n%s", filePath, err.Error())
	}
}

func main() {
	args := ParseCliArgs()
    homeDir := os.Getenv("HOME")
    if homeDir == "" {
		log.Fatalf("Can't read env $HOME")
    }
	configFile := homeDir + "/.config/swaync-widgets/config.toml"

	// Try to read or create a config file
	file, err := os.ReadFile(configFile)
	if err != nil {
		createConfigFile(configFile)
		file, err = os.ReadFile(configFile)
		if err != nil {
			log.Fatalf("Can't read config file at \"%s\". Output is: \n%s", configFile, err.Error())
		}
	}

	// Config file should be avaliabe now
	var cfg Config
	err = toml.Unmarshal(file, &cfg)
	if err != nil {
		log.Fatalf("Error parsing config file at \"%s\". Output is: \n%s", configFile, err.Error())
	}

	switch args.widget {
	case "mute":
		ToggleWidget(cfg.WidgetMute)
	case "vpn":
		ToggleWidget(cfg.WidgetVpn)
	case "wifi":
		ToggleWidget(cfg.WidgetWifi)
	case "bluetooth":
		ToggleWidget(cfg.WidgetBluetooth)
	default:
	}

	UpdateConfigFiles(cfg)
	ReloadConfigFiles(cfg)
}

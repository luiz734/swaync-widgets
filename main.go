package main

import (
	// "fmt"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/pelletier/go-toml/v2"
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
	widgetsData := GetWidgetsJsonData(cfg)
	outputCss := cfg.CSSPrepend

	outputCss += GetWidgetCss(cfg, cfg.WidgetVpn)
	widgetsData = UpdateWidgetBasedOnState(cfg, cfg.WidgetVpn, widgetsData)

	outputCss += GetWidgetCss(cfg, cfg.WidgetMute)
	widgetsData = UpdateWidgetBasedOnState(cfg, cfg.WidgetMute, widgetsData)

	outputCss += GetWidgetCss(cfg, cfg.WidgetWifi)
	widgetsData = UpdateWidgetBasedOnState(cfg, cfg.WidgetWifi, widgetsData)

	outputCss += GetWidgetCss(cfg, cfg.WidgetBluetooth)
	widgetsData = UpdateWidgetBasedOnState(cfg, cfg.WidgetBluetooth, widgetsData)

	err := os.WriteFile(cfg.SwayncCssWidgets, []byte(outputCss), 0755)
	if err != nil {
		log.Fatalf("Can't write file at \"%s\". Output is: \n%s", cfg.SwayncCssWidgets, err.Error())
	}
	UpdateConfigFile(cfg, widgetsData)
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

type WidgetJsonData struct {
	Label   string
	Command string
}

func GetWidgetsJsonData(cfg Config) []WidgetJsonData {
	configFile := cfg.SwayncConfigFile
	file, err := os.ReadFile(cfg.SwayncConfigFile)
	if err != nil {
		log.Fatalf("Can't read config file at \"%s\". Output is: \n%s", configFile, err.Error())
	}

	widgetsJsonPath := []string{"widget-config", "buttons-grid", "actions"}
	rawBytes, _, _, err := jsonparser.Get(file, widgetsJsonPath...)

	if err != nil {
		log.Fatalf("Can't parse  file at \"%s\". Output is: \n%s", configFile, err.Error())
	}

	readJsonKey := func(data []byte, key string) string {
		value, _, _, err := jsonparser.Get(data, key)
		if err != nil {
			log.Panicf("Can't read json key \"%s\" file at \"%s\". Output is: \n%s", key, configFile, err.Error())
		}
		return string(value)
	}

	widgets := []WidgetJsonData{}
	jsonparser.ArrayEach(rawBytes, func(data []byte, dataType jsonparser.ValueType, offset int, err error) {

		widgets = append(widgets, WidgetJsonData{
			Label:   readJsonKey(data, "label"),
			Command: readJsonKey(data, "command"),
		})
	})

	return widgets
}

func UpdateWidgetBasedOnState(cfg Config, widgetConfig WidgetConfig, widgetJsonData []WidgetJsonData) []WidgetJsonData {
	stateOn := GetWidgetState(widgetConfig.CheckStatusCommand)
	index, err := strconv.Atoi(widgetConfig.Index)
	if err != nil {
		log.Panicf("Can't convert \"%s\" to integet. Check your config for widget \"%s\"", widgetConfig.Index, widgetConfig.Desc)
	}
	// index in config file starts in 1 because CSS
	index = index - 1
	// targetWidget := widgetJsonData[index]
	// fmt.Println(targetWidget, "\n", widgetJsonData[index])

	if stateOn {
		widgetJsonData[index].Label = widgetConfig.OnLabel
	} else {
		widgetJsonData[index].Label = widgetConfig.OffLabel
	}

	return widgetJsonData
}

func UpdateConfigFile(cfg Config, widgetsJsonData []WidgetJsonData) {
	configFile := cfg.SwayncConfigFile
	file, err := os.ReadFile(cfg.SwayncConfigFile)
	if err != nil {
		log.Fatalf("Can't read config file at \"%s\". Output is: \n%s", configFile, err.Error())
	}

	for i, value := range widgetsJsonData {
		jsonIndex := "[" + strconv.Itoa(i) + "]"
		// fmt.Println(jsonIndex)
		widgetsJsonPath := []string{"widget-config", "buttons-grid", "actions", jsonIndex, "label"}
        labelDoubleQuotes := fmt.Sprintf("\"%s\"", value.Label)
		file, err = jsonparser.Set(file, []byte(labelDoubleQuotes), widgetsJsonPath...)
		if err != nil {
			log.Fatalf("Can't parse  file at \"%s\". Output is: \n%s", configFile, err.Error())
		}
	}

	err = os.WriteFile(cfg.SwayncConfigFile, file, 0755)
	if err != nil {
		log.Fatalf("Can't write file at \"%s\". Output is: \n%s", cfg.SwayncConfigFile, err.Error())
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

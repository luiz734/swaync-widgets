package app

import (
	"fmt"
	"github.com/buger/jsonparser"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"swaync-widgets/config"
)

func GetWidgetCss(cfg config.Config, widgetConfig config.WidgetConfig) string {
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

func UpdateConfigFiles(cfg config.Config) {
	widgetsData := GetWidgetsJsonData(cfg)
	outputCss := cfg.CSSPrepend

	for _, w := range cfg.Widgets {
		outputCss += GetWidgetCss(cfg, w)
		widgetsData = UpdateWidgetBasedOnState(cfg, w, widgetsData)
	}

	err := os.WriteFile(cfg.SwayncCssWidgets, []byte(outputCss), 0755)
	if err != nil {
		log.Fatalf("Can't write file at \"%s\". Output is: \n%s", cfg.SwayncCssWidgets, err.Error())
	}
	UpdateConfigFile(cfg, widgetsData)
}

func GetOnCss(cfg config.Config, index string, comment string) string {
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

func GetWidgetsJsonData(cfg config.Config) []WidgetJsonData {
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

func UpdateWidgetBasedOnState(cfg config.Config, widgetConfig config.WidgetConfig, widgetJsonData []WidgetJsonData) []WidgetJsonData {
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

func UpdateConfigFile(cfg config.Config, widgetsJsonData []WidgetJsonData) {
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

func GetOffCss(cfg config.Config, index string, comment string) string {
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
func ToggleWidget(widgetConfig config.WidgetConfig) {
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

func ReloadConfigFiles(cfg config.Config) {
	cmd := exec.Command("bash", "-c", cfg.SwayncReloadCommand)
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("Can't reload config with command \"%s\". Output is: \n%s", cfg.SwayncReloadCommand, err.Error())
	}
}

package app

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"swaync-widgets/config"

	"github.com/buger/jsonparser"
)

type WidgetJsonData struct {
	Label   string
	Command string
}

func ReadWidgetsJsonData(cfg config.Config) []WidgetJsonData {
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


func WriteConfigFile(cfg config.Config, widgetsJsonData []WidgetJsonData) {
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


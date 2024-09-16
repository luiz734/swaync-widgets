package app

import (
	"fmt"
	"os"
	"strconv"
	"swaync-widgets/config"

	"github.com/buger/jsonparser"
)

type WidgetJsonData struct {
	Label   string
	Command string
}

func ReadWidgetsJsonData(cfg config.Config) ([]WidgetJsonData, error) {
	configFile := cfg.SwayncConfigFile
	file, err := os.ReadFile(cfg.SwayncConfigFile)
	if err != nil {
		return nil, fmt.Errorf("can't read file %s: %w", configFile, err)
	}

	widgetsJsonPath := []string{"widget-config", "buttons-grid", "actions"}
	rawBytes, _, _, err := jsonparser.Get(file, widgetsJsonPath...)
	if err != nil {
		return nil, fmt.Errorf("can't parse file %s: %w", configFile, err)
	}

	readJsonKey := func(data []byte, key string) (string, error) {
		value, _, _, err := jsonparser.Get(data, key)
		if err != nil {
			return "", fmt.Errorf("can't read json key %s: %w", key, err)
		}
		return string(value), nil
	}
	widgets := []WidgetJsonData{}
	jsonparser.ArrayEach(rawBytes, func(data []byte, dataType jsonparser.ValueType, offset int, err error) {
		var label string
		var command string
		if label, err = readJsonKey(data, "label"); err != nil {
			err = fmt.Errorf("can't read label key: %w", err)
			return
		}
		if command, err = readJsonKey(data, "command"); err != nil {
			err = fmt.Errorf("can't read command key: %w", err)
			return
		}
		widgets = append(widgets, WidgetJsonData{
			Label:   label,
			Command: command,
		})
	})

	if err != nil {
		return nil, fmt.Errorf("error processing widgets array: %w", err)
	}

	return widgets, nil
}

func WriteConfigFile(cfg config.Config, widgetsJsonData []WidgetJsonData) error {
	file, err := os.ReadFile(cfg.SwayncConfigFile)
	if err != nil {
		return fmt.Errorf("can't read file %s: %w", cfg.SwayncConfigFile, err)
	}

	for i, value := range widgetsJsonData {
		jsonIndex := "[" + strconv.Itoa(i) + "]"
		widgetsJsonPath := []string{"widget-config", "buttons-grid", "actions", jsonIndex, "label"}
		labelDoubleQuotes := fmt.Sprintf("\"%s\"", value.Label)
		file, err = jsonparser.Set(file, []byte(labelDoubleQuotes), widgetsJsonPath...)
		if err != nil {
			return fmt.Errorf("can't parse file %s: %w", cfg.SwayncConfigFile, err)
		}
	}

	if err := os.WriteFile(cfg.SwayncConfigFile, file, 0755); err != nil {
		return fmt.Errorf("can't write file %s: %w", cfg.SwayncConfigFile, err)
	}
	return nil
}

package app

import (
	"fmt"
	"strconv"
	"swaync-widgets/config"
)

func UpdateWidgetBasedOnState(cfg config.Config, widgetConfig config.WidgetConfig, widgetJsonData []WidgetJsonData) ([]WidgetJsonData, error) {
    stateOn := RunGetWidgetState(widgetConfig.CheckStatusCommand)
	index, err := strconv.Atoi(widgetConfig.Index)
	if err != nil {
        return nil, fmt.Errorf("can't convert %s to integer: %w :", err)
	}
	// index in config file starts in 1 because CSS
	index = index - 1
    if index >= len(widgetJsonData) {
        return nil, fmt.Errorf("index %s out of bounds: ", widgetConfig.Index)
    }

	if stateOn {
		widgetJsonData[index].Label = widgetConfig.OnLabel
	} else {
		widgetJsonData[index].Label = widgetConfig.OffLabel
	}

	return widgetJsonData, nil
}


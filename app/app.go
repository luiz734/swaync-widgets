package app

import (
	"log"
	"strconv"
	"swaync-widgets/config"
)

func UpdateWidgetBasedOnState(cfg config.Config, widgetConfig config.WidgetConfig, widgetJsonData []WidgetJsonData) []WidgetJsonData {
	stateOn := RunGetWidgetState(widgetConfig.CheckStatusCommand)
	index, err := strconv.Atoi(widgetConfig.Index)
	if err != nil {
		log.Panicf("Can't convert \"%s\" to integet. Check your config for widget \"%s\"", widgetConfig.Index, widgetConfig.Desc)
	}
	// index in config file starts in 1 because CSS
	index = index - 1

	if stateOn {
		widgetJsonData[index].Label = widgetConfig.OnLabel
	} else {
		widgetJsonData[index].Label = widgetConfig.OffLabel
	}

	return widgetJsonData
}


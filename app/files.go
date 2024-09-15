package app

import (
	"log"
	"os"
	"swaync-widgets/config"
)

func WriteConfigAndCss(cfg config.Config) {
	widgetsData := ReadWidgetsJsonData(cfg)
	outputCss := cfg.CSSPrepend

	for _, w := range cfg.Widgets {
		outputCss += GenerateWidgetCss(cfg, w)
		widgetsData = UpdateWidgetBasedOnState(cfg, w, widgetsData)
	}

	err := os.WriteFile(cfg.SwayncCssWidgets, []byte(outputCss), 0755)
	if err != nil {
		log.Fatalf("Can't write file at \"%s\". Output is: \n%s", cfg.SwayncCssWidgets, err.Error())
	}
	WriteConfigFile(cfg, widgetsData)
}



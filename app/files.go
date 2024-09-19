package app

import (
	"fmt"
	"os"
	"swaync-widgets/config"
)

func WriteConfigAndCss(cfg config.Config) error {
	widgetsData, err := ReadWidgetsJsonData(cfg)
	if err != nil {
		return fmt.Errorf("error getting widgets data: %w", err)
	}
	outputCss := cfg.CSSPrepend

	for _, w := range cfg.Widgets {
		widgetCss, err := GenerateWidgetCss(cfg, w)
		if err != nil {
			return fmt.Errorf("can't generate css: %w", err)
		}
		outputCss += widgetCss
		widgetsData, err = UpdateWidgetBasedOnState(cfg, w, widgetsData)
		if err != nil {
			return fmt.Errorf("error retrieving data for widget %s: %w", w.Desc, err)
		}
	}

	err = os.WriteFile(cfg.SwayncCssWidgets, []byte(outputCss), 0755)
	if err != nil {
		return fmt.Errorf("can't write css file: %w", err)
	}
	if err := WriteConfigFile(cfg, widgetsData); err != nil {
		return fmt.Errorf("can't write config file: %w", err)
	}
	return nil
}

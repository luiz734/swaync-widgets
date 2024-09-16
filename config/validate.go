package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

func ValidateConfig(cfg Config) error {
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return fmt.Errorf("validation failed: %w: ", err)
	}
	err := validateWidgetsConfig(cfg.Widgets)
	return err
}

func validateWidgetsConfig(cfgWidgets []WidgetConfig) error {
	validate := validator.New()
	for i, widget := range cfgWidgets {
		if err := validate.Struct(widget); err != nil {
			return fmt.Errorf("validation failed for widget #%d: %w", i+1, err)
		}
	}
	return nil
}

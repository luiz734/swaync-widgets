package config

type Config struct {
	SwayncCssWidgets       string         `toml:"swaync_css_widgets" validate:"required"`
	SwayncConfigFile       string         `toml:"swaync_config_file" validate:"required"`
	SwayncReloadCommand    string         `toml:"swaync_reload_command" validate:"required"`
	CSSPrepend             string         `toml:"css_prepend" validate:"required"`
	CSSButtonSelector      string         `toml:"css_button_selector" validate:"required"`
	CSSButtonHoverSelector string         `toml:"css_button_hover_selector" validate:"required"`
	CSSLabelSelector       string         `toml:"css_label_selector" validate:"required"`
	CSSLabelHoverSelector  string         `toml:"css_label_hover_selector" validate:"required"`
	StylesOn               Styles         `toml:"styles_on" validate:"required"`
	StylesOff              Styles         `toml:"styles_off" validate:"required"`
	Widgets                []WidgetConfig `toml:"widget" validate:"required"`
}

type WidgetConfig struct {
	Desc               string `toml:"desc" validate:"required"`
	Index              string `toml:"index" validate:"required"`
	OffLabel           string `toml:"off_label" validate:"required"`
	OnLabel            string `toml:"on_label" validate:"required"`
	TurnOnCommand      string `toml:"turn_on_command" validate:"required"`
	TurnOffCommand     string `toml:"turn_off_command" validate:"required"`
	CheckStatusCommand string `toml:"check_status_command" validate:"required"`
}

type Styles struct {
	CssButton      string `toml:"css_button" validate:"required"`
	CssButtonHover string `toml:"css_button_hover" validate:"required"`
	CssLabel       string `toml:"css_label" validate:"required"`
	CssLabelHover  string `toml:"css_label_hover" validate:"required"`
}

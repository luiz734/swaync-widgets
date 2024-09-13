package config

type Config struct {
	Foo                    string       `toml:"foo"`
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

type WidgetConfig struct {
	Desc               string `toml:"desc"`
	Index              string `toml:"index"`
	OffLabel           string `toml:"off_label"`
	OnLabel            string `toml:"on_label"`
	TurnOnCommand      string `toml:"turn_on_command"`
	TurnOffCommand     string `toml:"turn_off_command"`
	CheckStatusCommand string `toml:"check_status_command"`
}

type Styles struct {
	CssButton      string `toml:"css_button"`
	CssButtonHover string `toml:"css_button_hover"`
	CssLabel       string `toml:"css_label"`
	CssLabelHover  string `toml:"css_label_hover"`
}

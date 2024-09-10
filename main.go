package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

type StylesOn struct {
    CssButton string `toml:"css_button"`
    CssButtonHover string `toml:"css_button_hover"`
    CssLabel string `toml:"css_label"`
    CssLabelHover string `toml:"css_label_hover"`
}

type StylesOff struct {
    CssButton string `toml:"css_button"`
    CssButtonHover string `toml:"css_button_hover"`
    CssLabel string `toml:"css_label"`
    CssLabelHover string `toml:"css_label_hover"`
}

type WidgetConfig struct {
    Desc string `toml:"desc"`
    Index string `toml:"index"`
    OffLabel string `toml:"off_label"`
    OnLabel string `toml:"on_label"`
    TurnOnCommand string `toml:"turn_on_command"`
    TurnOffCommand string `toml:"turn_off_command"`
    CheckStatusCommand string `toml:"check_status_command"`
}


type Config struct {
    SwayncCssWidgets string `toml:"swaync_css_widgets"`
    SwayncConfigFile string `toml:"swaync_config_file"`
    SwayncReloadCommand string `toml:"swaync_reload_command"`
    CSSPrepend string `toml:"css_prepend"`
    CSSButtonSelector string `toml:"css_button_selector"`
    CSSButtonHoverSelector string `toml:"css_button_hover_selector"`
    CSSLabelSelector string `toml:"css_label_selector"`
    CSSLabelHoverSelector string `toml:"css_label_hover_selector"`
    StylesOn StylesOn `toml:"styles_on"`
    StylesOff StylesOff `toml:"styles_off"`
    WidgetVpn WidgetConfig `toml:"vpn"`
    WidgetMute WidgetConfig `toml:"mute"`
}

func get_widget_css(cfg Config, widgetConfig WidgetConfig) string {
    state_on := get_widget_state(widgetConfig.CheckStatusCommand)
    if state_on {
        return get_on_css(cfg, widgetConfig.Index, widgetConfig.Desc)
    } else {
        return get_off_css(cfg, widgetConfig.Index, widgetConfig.Desc)
    }
}

func get_widget_state(command string) bool {
    cmd := exec.Command("bash", "-c", command) 
    out, err := cmd.Output()
    state_on := err == nil && len(out) != 0
    println(state_on, " ", " cmd: " + command)
    return state_on
}

func UpdateConfigFiles(cfg Config) {
    outputCss := cfg.CSSPrepend

    outputCss += get_widget_css(cfg, cfg.WidgetVpn)
    sedConfigFile(cfg, cfg.WidgetVpn)

    outputCss += get_widget_css(cfg, cfg.WidgetMute)
    sedConfigFile(cfg, cfg.WidgetMute)

    err := os.WriteFile(cfg.SwayncCssWidgets, []byte(outputCss), 0755)
    if err != nil {
        panic(err.Error())
    }

}

func get_on_css(cfg Config, index string, comment string) string {
    output := "/* widget " +  comment + " */\n"
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

func sedConfigFile(cfg Config, widgetConfig WidgetConfig) {
    state_on := get_widget_state(widgetConfig.CheckStatusCommand)
    firstPart := widgetConfig.OnLabel
    secondPart := widgetConfig.OffLabel
    var sedCommand string
    if state_on {
        firstPart = widgetConfig.OffLabel
        secondPart = widgetConfig.OnLabel
    }
    sedCommand = fmt.Sprintf(
        "sed -i 's/\"label\": \"%s\"/\"label\": \"%s\"/' \"%s\"",
        firstPart,
        secondPart,
        cfg.SwayncConfigFile)

    cmd := exec.Command("bash", "-c", sedCommand) 
    _, err := cmd.Output()
    if err != nil {
        panic(err.Error)
    }
}

func get_off_css(cfg Config, index string, comment string) string { 
    output := "/* widget " +  comment + " */\n"
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

type CliArgs struct {
    widget string
    action string
}

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func parse_cli() CliArgs {
    widget := ""
    action := ""
    if len(os.Args) > 1 {
        widget = os.Args[1]
    }
    if len(os.Args) > 2 {
        widget = os.Args[2]
    }

    if !contains([]string{"mute", "vpn", ""}, widget) {
        panic("Invalid option")
    }

    if !contains([]string{"on", "off", "toggle", ""}, action) {
        panic("Invalid action")
    }

    return CliArgs {
        widget: widget,
        action: action,
    }
}


func toggle_widget(widgetConfig WidgetConfig) {
    state_on := get_widget_state(widgetConfig.CheckStatusCommand)
    var command string
    if state_on {
        command = widgetConfig.TurnOffCommand
    } else {
        command = widgetConfig.TurnOnCommand
    }
    cmd := exec.Command("bash", "-c", command) 
    _, err := cmd.Output()
    if err != nil {
        panic(err.Error)
    }
}

func ReloadConfigFiles(cfg Config) {
    cmd := exec.Command("bash", "-c", cfg.SwayncReloadCommand) 
    _, err := cmd.Output()
    if err != nil {
        panic(err.Error)
    }
}

func main() {
    args := parse_cli()
    
    file, err := os.ReadFile("config.toml")
    if err != nil {
        panic(err)
    }
    var cfg Config
    err = toml.Unmarshal(file, &cfg)


    switch args.widget {
    case "mute":
        toggle_widget(cfg.WidgetMute)
    case "vpn":
        toggle_widget(cfg.WidgetVpn)
    case "":
    }

    UpdateConfigFiles(cfg)
    ReloadConfigFiles(cfg)
}

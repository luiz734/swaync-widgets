package app

import (
	// "fmt"
	// "github.com/buger/jsonparser"
	// "log"
	// "os"
	// "os/exec"
	// "strconv"

	"strings"
	"swaync-widgets/config"
)

func GenerateWidgetCss(cfg config.Config, widgetConfig config.WidgetConfig) (string, error) {
    if stateOn := RunGetWidgetState(widgetConfig.CheckStatusCommand); stateOn {
		return GenerateOnCss(cfg, widgetConfig.Index, widgetConfig.Desc), nil
	}
	return GenerateOffCss(cfg, widgetConfig.Index, widgetConfig.Desc), nil
}

func GenerateOnCss(cfg config.Config, index string, comment string) string {
	output := "/* widget " + comment + " */\n"
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

func GenerateOffCss(cfg config.Config, index string, comment string) string {
	output := "/* widget " + comment + " */\n"
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

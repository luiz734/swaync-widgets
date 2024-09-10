OUTPUT="swaync-widgets"
[ -d "build" ] || mkdir build
go build
mv "$OUTPUT" build
cp "build/$OUTPUT" "/home/$USER/.local/bin"

CONFIG_DIR="$HOME/.config/swaync-widgets"
CONFIG_FILE="$CONFIG_DIR/config.toml"
[ -d "$CONFIG_DIR" ] || mkdir "$CONFIG_DIR"
[ -f "$CONFIG_FILE" ] || cp "config.toml" "$CONFIG_FILE"



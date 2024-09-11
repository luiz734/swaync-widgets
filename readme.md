# swaync-widgets
This tool allows you to dynamically update the `buttons-grid` in your [swaync](https://github.com/ErikReider/SwayNotificationCenter) `config.json` and `styles.css`.

**Disclaimer:** I've maid this tool for **me**. It's working as I intended by now. I may add features or do bugfixes, but it's always better to clone the repo if you need stuff to be done now. The source code is very compact, so adding new features should not be complicated.

## How it works
This is a snipped from my config file at `~/.config/swaync/config.json`
```json
{
    "widget-config": {
        "buttons-grid": {
            "actions": [
                {
                    "label": "   Connected",
                    "command": "swaync-widgets wifi"
                },
                {
                    "label": "󰕾  Mute",
                    "command": "swaync-widgets mute"
                }
                // more stuff
            ]
    }
}
```

For each widget in your `swaync-widgets` config file, you can define commands to:
- Define custom css based on states (on/off)
- Check the status (state) of a widget
- Enable a widget
- Disable a widget

For example, consider the default configuration for the `mute` widget
```toml
[mute]
desc="mute"
index="2"
off_label="󰕾  Mute"
on_label="󰕾  Muted"
turn_on_command="pactl set-sink-mute @DEFAULT_SINK@ on"
turn_off_command="pactl set-sink-mute @DEFAULT_SINK@ off"
check_status_command="pactl get-sink-mute @DEFAULT_SINK@ | grep \"yes\""
```

Lets do a quick tour for each option:
- desc: a description to make `widgets.css` more legible
- index: used on css selector like `div:nth-child(X)`
- off_label: the label when state is off
- on_label: the label when state is on
- turn_on_command: the command that will run when state is off and you call `swaync-widgets widget`
- turn_off_command: the command that will run when state is on and you call `swaync-widgets widget`
- check_status_command: if ouput is empty or return code is not `0`, state is considered off

## CSS selectors
Currently, it supports 4 selectors:
- The widget button
- The widget button when hover
- The widget label
- The widget label when hover

For each selector, you define css properties for states `on` and `off`.
- It will apply the properties based on the widget state.
- You can't define custom css per widget, only per state

## Using the tool
When you call `swaync-widgets` without any args, it updates:
- The swaync `config.json` file by using the `sed` command and the labels
    -  Example: `sed -i 's/"label": "󰝟  Muted"/"label": "󰕾  Mute"/' "/home/$USER/.config/swaync/config.json"`.
    - The `sed` command is generated based on the labels and the current state of the widget: if the widget is active, it will try to replace the `on_label` to `off_label`. If it is not active, it will do the oposite.
- A new css file at `~/.config/swaync/widgets.css`. This file is generated by `swaync-widgets`. You just need to import it to your `styles.css` file.

### Toggle widgets
The command on the snippet is `"command": "swaync-widgets wifi"`. Internally, it checks the current state by using the `check_status_command` for the widget, and run the apropriate command in response. Then it reloads both your `config.json` and `styles.css` config files.

> A full explanation of each entry is on swaync-widgets default config.

## Setup
- Clone the repo and run the install script
- By default, it installs on `~/.local/bin`, so no root access is needed
- I encourage you to read the install script before running it

### ~/.config/swaync/config.json
- Because it uses `sed`, you initial `swaync` configuration must match the `swaync-widgets` config or nothing will happens to the label when you run the command.
- For example, if your swaync config has the entry `"label": "   Connected"`, your `swaync-widget` must match `   Connected` like this `on_label=   Connected`.
- If only you CSS is changed, probably there is a typo and they don't match.


### ~/.config/swaync/style.css
- You need to manually include the follow import to your current config `@import url("widgets.css");`
- `swaync-widget` doesn't touch your css config file.
- If like me, you use colors, you can add them in the `prepend_css` field on your config: `css_prepend="@define-color pink rgb(245, 194, 231);`
- Everything there will be put before any css rules

### ~/.config/swaync/widgets.css
This is an example of a generated `widgets.css` file. It's content is made entirely by the content of `~/.config/swaync-widgets/config.toml`.
> This file is always changing. Don't source control it

> You will never need to touch this file, but if something is broken, investigating it may help debug
```css
@define-color pink rgb(245, 194, 231); @define-color crust rgb(17, 17, 27); @define-color surface0 #313244; @define-color text rgb(205, 214, 244);
/* widget vpn */
.widget-buttons-grid>flowbox>flowboxchild:nth-child(4)>button{background: @transparent; border: 2px solid @surface0}
.widget-buttons-grid>flowbox>flowboxchild:nth-child(4)>button:hover{background: @transparent; border: 2px solid @pink}
.widget-buttons-grid>flowbox>flowboxchild:nth-child(4)>button>label{color: @text;}
.widget-buttons-grid>flowbox>flowboxchild:nth-child(4)>button:hover>label{color: @pink;}
/* widget mute */
.widget-buttons-grid>flowbox>flowboxchild:nth-child(2)>button{background: @transparent; border: 2px solid @surface0}
.widget-buttons-grid>flowbox>flowboxchild:nth-child(2)>button:hover{background: @transparent; border: 2px solid @pink}
.widget-buttons-grid>flowbox>flowboxchild:nth-child(2)>button>label{color: @text;}
.widget-buttons-grid>flowbox>flowboxchild:nth-child(2)>button:hover>label{color: @pink;}
/* widget wifi */
.widget-buttons-grid>flowbox>flowboxchild:nth-child(1)>button{background: @pink; border: 2px solid @pink}
.widget-buttons-grid>flowbox>flowboxchild:nth-child(1)>button:hover{background: @pink; border: 2px solid @pink}
.widget-buttons-grid>flowbox>flowboxchild:nth-child(1)>button>label{color: @crust;}
.widget-buttons-grid>flowbox>flowboxchild:nth-child(1)>button:hover>label{color: @crust;}
/* widget bluetooth */
.widget-buttons-grid>flowbox>flowboxchild:nth-child(3)>button{background: @transparent; border: 2px solid @surface0}
.widget-buttons-grid>flowbox>flowboxchild:nth-child(3)>button:hover{background: @transparent; border: 2px solid @pink}
.widget-buttons-grid>flowbox>flowboxchild:nth-child(3)>button>label{color: @text;}
.widget-buttons-grid>flowbox>flowboxchild:nth-child(3)>button:hover>label{color: @pink;}
```

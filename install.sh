#!/bin/env bash
OUTPUT="swaync-widgets"
[ -d "build" ] || mkdir build
go build
mv "$OUTPUT" build
cp "build/$OUTPUT" "/home/$USER/.local/bin"

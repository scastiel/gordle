#!/bin/sh

PACKAGE=gordle
FILENAME=src/resources.go

fyne bundle assets/AppIcon.svg > "$FILENAME"
fyne bundle -append assets/AppIcon.png >> "$FILENAME"
fyne bundle -append assets/example.png >> "$FILENAME"
fyne bundle -append assets/about_part1.md >> "$FILENAME"
fyne bundle -append assets/about_part2.md >> "$FILENAME"

sed -i "" -e "s/package main/package $PACKAGE/" "$FILENAME"
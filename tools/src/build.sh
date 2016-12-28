#!/bin/sh
set -e
DIRECTORY="`dirname \"${0}\"`"
TARGET_DIRECTORY="${DIRECTORY}\..\bin"

GOOS=windows GOARH=amd64 go build -o "${TARGET_DIRECTORY}/update_hugo.exe" "${DIRECTORY}/update_hugo.go"
GOOS=linux GOARH=amd64 go build -o "${TARGET_DIRECTORY}/update_hugo" "${DIRECTORY}/update_hugo.go"

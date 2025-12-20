#!/usr/bin/bash

LOCAL_REPO_PATH="/home/kyle/repos/rpi-rgb-screen"

date

echo "Changing directory to: $LOCAL_REPO_PATH"
cd $LOCAL_REPO_PATH

echo "Fetching latest changes from repo..."
git pull

if [ $? -ne 0 ]; then
    echo "Failed to fetch latest changes from remote repository."
    exit 1
fi

echo "Git pull completed successfully."

echo "Clearing the build cache..."
go clean -cache

echo "Building the app..."
go build -o rpi-rgb-screen main.go
if [ $? -ne 0 ]; then
    echo "Build failed."
    exit 1
fi
echo "Build completed successfully."

echo "Starting up the app..."
sudo ./rpi-rgb-screen --led-slowdown-gpio=4

#!/bin/bash

# Paths and variables
XCODE_PROJECT_PATH="./area-mobile"                  # Relative path to Xcode folder
SCHEME_NAME="area"                                 # Scheme name (confirmed)
APP_NAME="area"                                    # App name (confirmed)
BUNDLE_ID="dydy.area"                              # Bundle ID confirmed
SIMULATOR_NAME="iPhone 17"                         # Updated simulator (iOS 26.0)
SIMULATOR_OS="26.0"                                # Updated iOS version

echo "Starting Docker Compose (backend, frontend, DB)..."
docker-compose down
docker-compose up --build -d

echo "Waiting for backend to be ready..."
until curl -s http://localhost:8080/health > /dev/null; do
  sleep 1
done
echo "Backend is ready!"

echo "What would you like to do?"
echo "1. Launch iOS app in simulator ($SIMULATOR_NAME)"
echo "2. Generate .ipa for connected device (requires device connected via USB)"
echo "3. Exit"
read -p "Enter your choice (1, 2, or 3): " choice

if [ "$choice" = "1" ]; then
  echo "Compiling and launching in simulator $SIMULATOR_NAME..."
  cd "$XCODE_PROJECT_PATH" || { echo "Error: Xcode path not found"; exit 1; }
  xcodebuild build -scheme "$SCHEME_NAME" -destination "platform=iOS Simulator,name=$SIMULATOR_NAME,OS=$SIMULATOR_OS" || { echo "Error: Build failed"; exit 1; }
  DERIVED_DATA=$(xcodebuild -project ./area.xcodeproj -scheme "$SCHEME_NAME" -destination "platform=iOS Simulator,name=$SIMULATOR_NAME,OS=$SIMULATOR_OS" -showBuildSettings | grep " BUILT_PRODUCTS_DIR =" | awk '{print $3}')
  APP_PATH="$DERIVED_DATA/$APP_NAME.app"
  if [ ! -d "$APP_PATH" ]; then
    echo "Error: File $APP_PATH not found."
    echo "Check the path in ~/Library/Developer/Xcode/DerivedData/area-<hash>/Build/Products/Debug-iphonesimulator/"
    exit 1
  fi
  echo "App found at: $APP_PATH"
  open -a Simulator
  SIMULATOR_STATUS=$(xcrun simctl list devices | grep "$SIMULATOR_NAME" | grep "(Booted)")
  if [ -n "$SIMULATOR_STATUS" ]; then
    echo "Simulator already started, attempting restart to avoid conflicts..."
    xcrun simctl shutdown "$SIMULATOR_NAME" 2>/dev/null
    sleep 2 
  fi
  xcrun simctl boot "$SIMULATOR_NAME" 2>/dev/null || { echo "Warning: Simulator boot failed, but it might already be started."; }
  xcrun simctl install "$SIMULATOR_NAME" "$APP_PATH" || { echo "Error: Installation failed"; exit 1; }
  xcrun simctl launch "$SIMULATOR_NAME" "$BUNDLE_ID" || { echo "Error: Launch failed. Check BUNDLE_ID ($BUNDLE_ID)"; exit 1; }
  echo "App launched in simulator!"

elif [ "$choice" = "2" ]; then
  echo "Generating .ipa for connected device..."
  echo "Make sure your iPhone/iPad is connected via USB and unlocked."
  cd "$XCODE_PROJECT_PATH" || { echo "Error: Xcode path not found"; exit 1; }
  DEVICE_UDID=$(xcrun xctrace list devices | grep -E "iPhone|iPad" | grep -v Simulator | head -n 1 | awk '{print $NF}' | tr -d '()')
  if [ -z "$DEVICE_UDID" ]; then
    echo "Error: No connected device detected. Connect an iPhone/iPad via USB."
    exit 1
  fi
  echo "Device detected: UDID $DEVICE_UDID"
  xcodebuild build -scheme "$SCHEME_NAME" -destination "platform=iOS,id=$DEVICE_UDID" || { echo "Error: Build failed"; exit 1; }
  xcodebuild archive -scheme "$SCHEME_NAME" -sdk iphoneos -archivePath "./build/$APP_NAME.xcarchive" || { echo "Error: Archive failed"; exit 1; }
  xcodebuild -exportArchive -archivePath "./build/$APP_NAME.xcarchive" -exportPath "./build/$APP_NAME" -exportOptionsPlist <(echo '<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>method</key>
    <string>development</string>
</dict>
</plist>') || { echo "Error: Export failed"; exit 1; }
  echo ".ipa generated in $XCODE_PROJECT_PATH/build/$APP_NAME !"
  echo "To install, open Xcode, go to Window > Devices and Simulators, and drag the .ipa onto the device."

elif [ "$choice" = "3" ]; then
  echo "Stopping script."
  exit 0
else
  echo "Invalid choice."
  exit 1
fi
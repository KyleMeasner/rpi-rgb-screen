package main

import (
	"os"
	"os/signal"
	"rpi-rgb-screen/constants"
	"rpi-rgb-screen/data"
	"rpi-rgb-screen/fonts"
	"rpi-rgb-screen/manager"
	"syscall"

	rgbmatrix "github.com/KyleMeasner/go-rpi-rgb-led-matrix"
)

func main() {
	config := &rgbmatrix.DefaultConfig
	config.Rows = constants.SCREEN_HEIGHT
	config.Cols = constants.SCREEN_WIDTH
	config.Brightness = 100
	config.HardwareMapping = "adafruit-hat-pwm"
	config.ShowRefreshRate = true

	matrix, err := rgbmatrix.NewRGBLedMatrix(config)
	if err != nil {
		panic(err)
	}
	go clearScreenOnExit(matrix)

	toolKit := rgbmatrix.NewToolKit(matrix)
	defer toolKit.Close()

	fontCache := fonts.LoadFonts()

	dataManager := data.NewDataManager()
	screenManager := manager.NewScreenManager(fontCache, toolKit, dataManager)
	screenManager.Run()
}

func clearScreenOnExit(matrix rgbmatrix.Matrix) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)
	<-signalChan

	// Cleanup actions
	matrix.Close()

	os.Exit(0)
}

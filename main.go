package main

import (
	"os"
	"os/signal"
	"rpi-rgb-screen/constants"
	"rpi-rgb-screen/fonts"
	"rpi-rgb-screen/screen"
	"rpi-rgb-screen/transition"
	"syscall"
	"time"

	rgbmatrix "github.com/KyleMeasner/go-rpi-rgb-led-matrix"
)

func main() {
	config := &rgbmatrix.DefaultConfig
	config.Rows = constants.SCREEN_HEIGHT
	config.Cols = constants.SCREEN_WIDTH
	config.Brightness = 100
	config.HardwareMapping = "adafruit-hat"
	config.ShowRefreshRate = true

	matrix, err := rgbmatrix.NewRGBLedMatrix(config)
	if err != nil {
		panic(err)
	}

	toolKit := rgbmatrix.NewToolKit(matrix)
	defer toolKit.Close()
	go clearScreenOnExit(toolKit.Canvas)

	fontCache := fonts.LoadFonts()

	screen1 := screen.NewDummyScreen(fontCache)
	for {
		screen2 := screen.NewDummyScreen(fontCache)

		err = toolKit.PlayAnimation(transition.NewSlideIn(screen1, screen2))
		if err != nil {
			panic(err)
		}
		time.Sleep(5 * time.Second)
		screen1 = screen2
	}
}

func clearScreenOnExit(canvas *rgbmatrix.Canvas) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)
	<-signalChan

	// Cleanup actions
	canvas.Close()

	os.Exit(0)
}

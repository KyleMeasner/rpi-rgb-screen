package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"rpi-rgb-screen/config"
	"rpi-rgb-screen/constants"
	"rpi-rgb-screen/data"
	"rpi-rgb-screen/fonts"
	"rpi-rgb-screen/manager"
	"syscall"
	"time"

	rgbmatrix "github.com/KyleMeasner/go-rpi-rgb-led-matrix"
)

func main() {
	// Setup logging
	err := os.MkdirAll("./logs", os.ModePerm)
	if err != nil {
		log.Fatalf("error creating logs directory: %v", err)
	}
	deleteOldLogs()

	logFileName := fmt.Sprintf("./logs/%s.txt", time.Now().Format(time.RFC3339))
	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Load config
	err = config.LoadConfig()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	// Initialize RGB matrix
	matrixConfig := &rgbmatrix.DefaultConfig
	matrixConfig.Rows = constants.SCREEN_HEIGHT
	matrixConfig.Cols = constants.SCREEN_WIDTH
	matrixConfig.Brightness = 100
	matrixConfig.HardwareMapping = "adafruit-hat-pwm"
	matrixConfig.ShowRefreshRate = false

	matrix, err := rgbmatrix.NewRGBLedMatrix(matrixConfig)
	if err != nil {
		log.Fatalf("error initializing RGB matrix: %v", err)
	}
	go clearScreenOnExit(matrix)

	canvas := rgbmatrix.NewCanvas(matrix)
	defer canvas.Close()

	// Load fonts
	fontCache := fonts.LoadFonts()

	// Initialize data and screen managers
	dataManager := data.NewDataManager()
	screenManager := manager.NewScreenManager(fontCache, canvas, dataManager)
	screenManager.Initialize()
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

func deleteOldLogs() {
	files, err := os.ReadDir("./logs")
	if err != nil {
		log.Fatalf("error reading logs directory: %v", err)
		return
	}

	cutoffTime := time.Now().AddDate(0, 0, -7) // 1 week of logs

	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			log.Printf("error getting file info for log file %s: %v", file.Name(), err)
			continue
		}

		if info.ModTime().Before(cutoffTime) {
			err := os.Remove("./logs/" + file.Name())
			if err != nil {
				log.Printf("error deleting old log file %s: %v", file.Name(), err)
			} else {
				log.Printf("deleted old log file: %s", file.Name())
			}
		}
	}
}

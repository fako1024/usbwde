package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/fako1024/usbwde"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

// Health denotes the result of a health check
type Health struct {
	OK      bool
	Details string
}

const maxDataAge = 2 * time.Minute

// Simple global variables to hold configuration / data
var (
	devicePath     string
	serverEndpoint string

	currentData *usbwde.DataPoint
	health      *Health
)

func main() {

	// Parse command line parameters
	readFlags()

	// Start the echo server
	go startServer()

	// Continuously try looping / extracting data (wrapped in additional loop since
	// the device occasionally loses connection)
	for {
		readLoop()

		// Back off for ten seconds to allow device to (re-)settle
		time.Sleep(10 * time.Second)
	}
}

// readLoop continuously reads lines from the device
func readLoop() {

	// Recover from potential panic when reading from device
	defer func() {
		if r := recover(); r != nil {
			logrus.StandardLogger().Errorf("Panic recovered in readLoop(): %s", r)
			health = &Health{
				OK:      false,
				Details: fmt.Sprintf("Panic recovered in readLoop(): %s", r),
			}
		}
	}()

	// Initialize a new USBWDE sensor / station
	sensor, err := usbwde.New(devicePath)
	if err != nil {
		logrus.StandardLogger().Errorf("Error opening %s: %s", devicePath, err)
		health = &Health{
			OK:      false,
			Details: fmt.Sprintf("Error opening %s: %s", devicePath, err),
		}

		return
	}
	defer sensor.Close()

	// Continuously throw a log message upon reception of updated data
	for {

		// Read single data point
		dataPoint, err := sensor.Read()
		if err != nil {
			logrus.StandardLogger().Errorf("Error reading data from %s: %s", devicePath, err)
			health = &Health{
				OK:      false,
				Details: fmt.Sprintf("Error reading data from %s: %s", devicePath, err),
			}
		}

		// Assign newly read data to current data
		currentData = dataPoint
		health = &Health{
			OK: true,
		}
	}
}

// readFlags parses command line parameters
func readFlags() {
	flag.StringVar(&devicePath, "d", "/dev/ttyUSB0", "Device / socket path to connect to")
	flag.StringVar(&serverEndpoint, "s", "0.0.0.0:8000", "Server endpoint to listen on")

	flag.Parse()
}

// startServer launches an echo middleware to listen for data requests
func startServer() {

	// Create echo server instance
	e := echo.New()

	// Routes
	e.GET("/", returnData)
	e.GET("/health", returnHealth)

	// Start server
	logrus.StandardLogger().Fatal(e.Start(serverEndpoint))
}

// Data return handler
func returnData(c echo.Context) error {

	// If there is no data (yet), signify via HTTP error
	if currentData == nil {
		return c.String(http.StatusNoContent, "No data yet")
	}

	return c.JSONPretty(http.StatusOK, currentData, "  ")
}

// Health handler
func returnHealth(c echo.Context) error {

	// If there is no health (yet), signify via HTTP error
	if health == nil {
		return c.String(http.StatusNoContent, "No health data yet")
	}

	// If the current health is not ok, return it
	if health.OK == false {
		return c.JSONPretty(http.StatusOK, health, "  ")
	}

	// If the data is too old, return a warning
	if currentData != nil && currentData.TimeStamp.Before(time.Now().Add(-2*time.Minute)) {
		return c.JSONPretty(http.StatusOK, &Health{
			OK:      false,
			Details: fmt.Sprintf("Data is older than %v", maxDataAge),
		}, "  ")
	}

	return c.JSONPretty(http.StatusOK, health, "  ")
}

package main

import (
	"flag"

	"github.com/fako1024/usbwde"
	"github.com/sirupsen/logrus"
)

var (
	devicePath string
)

func main() {

	// Parse command line parameters
	readFlags()

	// Initialize a new USBWDE sensor / station
	sensor, err := usbwde.New(devicePath)
	if err != nil {
		logrus.StandardLogger().Fatalf("Error opening %s: %s", devicePath, err)
	}
	defer sensor.Close()

	// Continuously throw a log message upon reception of updated data
	for {

		// Read single data point
		dataPoint, err := sensor.Read()
		if err != nil {
			logrus.StandardLogger().Errorf("Error reading data from %s: %s", devicePath, err)
		}

		// Log data
		logrus.StandardLogger().Infof("Read data from %s: %s", devicePath, dataPoint)
	}
}

// readFlags parses command line parameters
func readFlags() {
	flag.StringVar(&devicePath, "d", "/dev/ttyUSB0", "Device / socket path to connect to")

	flag.Parse()
}

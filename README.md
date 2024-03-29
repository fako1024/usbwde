# Simple go package to read data from USB-WDE weather data logger (ELV)

[![Github Release](https://img.shields.io/github/release/fako1024/usbwde.svg)](https://github.com/fako1024/usbwde/releases)
[![GoDoc](https://godoc.org/github.com/fako1024/usbwde?status.svg)](https://godoc.org/github.com/fako1024/usbwde/)
[![Go Report Card](https://goreportcard.com/badge/github.com/fako1024/usbwde)](https://goreportcard.com/report/github.com/fako1024/usbwde)
[![Build/Test Status](https://github.com/fako1024/usbwde/workflows/Go/badge.svg)](https://github.com/fako1024/usbwde/actions?query=workflow%3AGo)

This package allows to extract structured data from a USB-WDE weather station device (available from [here](https://www.elv.de/usb-wetterdaten-empfaenger-usb-wde1-komplettbausatz-1.html)). Usage is fairly trivial (see examples directory for a simple console logger implementation).

## Features
- Extraction of USBWDE RF sensor data
	- Up to 8 temperature / humidity sensors
	- Hybrid sensor providing temperature / humidity, wind speed and precipitation data

## Installation
```bash
go get -u github.com/fako1024/usbwde
```

## Example
```go
// Initialize a new USBWDE sensor / station
sensor, err := usbwde.New("/dev/ttyUSB0")
if err != nil {
    logrus.StandardLogger().Fatalf("Error opening /dev/ttyUSB0: %s", err)
}

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
```

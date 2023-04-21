package usbwde

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/jacobsa/go-serial/serial"
)

// USBWDE denotes an ELV USBWDE endpoint
type USBWDE struct {
	socket string
	port   io.ReadWriteCloser
}

// New creates a new USBWDE object
func New(socket string) (*USBWDE, error) {

	// Define default options for USBWDE device
	defaultOptions := serial.OpenOptions{
		PortName:        socket,
		BaudRate:        9600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 1,
	}

	// Open the port
	port, err := serial.Open(defaultOptions)
	if err != nil {
		return nil, err
	}

	// Create and return new object
	return &USBWDE{
		socket: socket,
		port:   port,
	}, nil
}

// Close closes the connection to the device
func (s *USBWDE) Close() error {
	return s.port.Close()
}

// Read extracts a single data point from the sensor / device
func (s *USBWDE) Read() (*DataPoint, error) {

	// Extract raw data line / bytes
	rawData, err := s.readRawData()
	if err != nil {
		return nil, err
	}

	return s.parse(rawData)
}

////////////////////////////////////////////////////////////////////////////////

// parse parses / normalizes the raw data received from the socket
func (s *USBWDE) parse(rawData []byte) (*DataPoint, error) {

	// Split result string by sepearator character and perform sanity check
	splitString := strings.Split(string(rawData), ";")
	if len(splitString) != 25 {
		return nil, fmt.Errorf("invalid data received: %s", string(rawData))
	}

	// Extract & format temperature values
	temperatures, err := normalizeTemperature(splitString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse temperature values from %s: %w", splitString, err)
	}

	// Extract & format humidity values
	humidities, err := normalizeHumidity(splitString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse humidity values from %s: %w", splitString, err)
	}

	hybridSensor, err := normalizeHybridSensor(splitString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse hybrid sensor values from %s: %w", splitString, err)
	}

	// Create & return a data point
	return &DataPoint{
		TimeStamp:    time.Now(),
		Temperature:  temperatures,
		Humidity:     humidities,
		HybridSensor: hybridSensor,
	}, nil
}

// readRawData extracts data from the port
func (s *USBWDE) readRawData() ([]byte, error) {

	// Wrap reader around port
	reader := bufio.NewReader(s.port)

	// Read full data line until termination signal is received
	reply, err := reader.ReadBytes('\x0a')
	if err != nil {
		return nil, err
	}

	// Return the raw data received
	return reply, nil
}

// normalize converts the (German) floating point format string to a well-defined
// float64 representation
func normalize(in string) (float64, error) {
	if in == "" {
		return 0.0, nil
	}
	return strconv.ParseFloat(strings.Replace(in, ",", ".", -1), 64)
}

// normalizeTemperature converts the temperature strings to well-defined values
func normalizeTemperature(in []string) (result [8]float64, err error) {
	for i, val := range in[3:11] {
		if result[i], err = normalize(val); err != nil {
			return
		}
	}
	return
}

// normalizeTemperature converts the humidity strings to well-defined values
func normalizeHumidity(in []string) (result [8]float64, err error) {
	for i, val := range in[11:19] {
		if result[i], err = normalize(val); err != nil {
			return
		}
	}
	return
}

// normalizeHybridSensor converts the hybrid sensor data to well-defined values
func normalizeHybridSensor(in []string) (result HybridSensor, err error) {

	if result.Temperature, err = normalize(in[19]); err != nil {
		return
	}
	if result.Humidity, err = normalize(in[20]); err != nil {
		return
	}
	if result.WindSpeed, err = normalize(in[21]); err != nil {
		return
	}
	if result.Precipitation, err = strconv.Atoi(in[22]); err != nil {
		return
	}

	result.IsRaining, err = strconv.ParseBool(in[23])

	return
}

package usbwde

import (
	"errors"
	"fmt"
	"math"
	"time"
)

const ValueDelta = 0.000001

// HybridSensor denotes the data from the hybrid sensor (if present)
type HybridSensor struct {
	Temperature   float64
	Humidity      float64
	WindSpeed     float64
	Precipitation int
	IsRaining     bool
}

// DataPoint denotes a set of data taken at a specific point in time
type DataPoint struct {
	TimeStamp   time.Time
	Temperature [8]float64
	Humidity    [8]float64

	HybridSensor HybridSensor
}

// String returns a well-formatted string for the data point, fulfilling the Stringer interface
func (p *DataPoint) String() string {
	return fmt.Sprintf("%s: (%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f,%.1f)°C (%g,%g,%g,%g,%g,%g,%g,%g)%% - Hybrid Sensor: %.1f°C, %g%%, %.1fkm/h, isRaining: %v (%d)", p.TimeStamp.Format(time.RFC1123),
		p.Temperature[0],
		p.Temperature[1],
		p.Temperature[2],
		p.Temperature[3],
		p.Temperature[4],
		p.Temperature[5],
		p.Temperature[6],
		p.Temperature[7],
		p.Humidity[0],
		p.Humidity[1],
		p.Humidity[2],
		p.Humidity[3],
		p.Humidity[4],
		p.Humidity[5],
		p.Humidity[6],
		p.Humidity[7],
		p.HybridSensor.Temperature,
		p.HybridSensor.Humidity,
		p.HybridSensor.WindSpeed,
		p.HybridSensor.IsRaining,
		p.HybridSensor.Precipitation)
}

// IsComplete returns if the data point has all fields set (and at least one of each temperature / Humidity is non-zero)
func (p *DataPoint) IsComplete(includeHybridSensor bool) (bool, error) {
	for i := 0; i < 8; i++ {
		if math.Abs(p.Temperature[i]) < ValueDelta && p.Humidity[i] < ValueDelta {
			return false, fmt.Errorf("missing temperature and humidity data for index %d", i)
		}
	}
	if includeHybridSensor {
		if math.Abs(p.HybridSensor.Temperature) < ValueDelta && p.HybridSensor.Humidity < ValueDelta &&
			math.Abs(p.HybridSensor.WindSpeed) < ValueDelta && !p.HybridSensor.IsRaining &&
			p.HybridSensor.Precipitation == 0 {
			return false, errors.New("missing hybrid sensor data")
		}
	}
	return true, nil
}

package usbwde

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

var testString = "$1;1;;21,1;21,2;20,8;18,4;21,6;21,2;20,8;;47;46;46;44;50;49;46;;3,1;30;8,0;455;1;0"

func TestParse(t *testing.T) {

	mock := &USBWDE{}
	data, err := mock.parse([]byte(testString))

	if err != nil {
		t.Fatalf("parsing failed: %s", err)
	}

	if data.String() != fmt.Sprintf("%s: (21.1,21.2,20.8,18.4,21.6,21.2,20.8,0.0)°C (47,46,46,44,50,49,46,0)%% - Hybrid Sensor: 3.1°C, 30%%, 8.0km/h, isRaining: true (455)", data.TimeStamp.Format(time.RFC1123)) {
		t.Fatalf("unexpected String() method result")
	}

	expected := DataPoint{
		TimeStamp: data.TimeStamp,
		Temperature: [8]float64{
			21.1, 21.2, 20.8, 18.4, 21.6, 21.2, 20.8,
		},
		Humidity: [8]float64{
			47, 46, 46, 44, 50, 49, 46,
		},
		HybridSensor: HybridSensor{
			Temperature:   3.1,
			Humidity:      30,
			WindSpeed:     8,
			Precipitation: 455,
			IsRaining:     true,
		},
	}

	if !reflect.DeepEqual(*data, expected) {
		t.Fatalf("unexpected data received, want %+v, have %+v", expected, *data)
	}

}

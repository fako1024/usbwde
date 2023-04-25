package usbwde

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

type testCase struct {
	input             string
	expected          string
	expectedDatapoint DataPoint
}

// var (
// 	testString  =
// 	testString1 = "$1;1;;;22,0;22,3;21,8;20,9;17,0;21,4;21,8;;50;50;47;48;60;52;47;;;;;;0\r\n"
// )

func TestParse(t *testing.T) {

	mock := &USBWDE{}

	for _, cs := range []testCase{
		{
			input:    "$1;1;;21,1;21,2;20,8;18,4;21,6;21,2;20,8;;47;46;46;44;50;49;46;;3,1;30;8,0;455;1;0\r\n",
			expected: "%s: (21.1,21.2,20.8,18.4,21.6,21.2,20.8,0.0)°C (47,46,46,44,50,49,46,0)%% - Hybrid Sensor: 3.1°C, 30%%, 8.0km/h, isRaining: true (455)",
			expectedDatapoint: DataPoint{
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
			},
		},
		{
			input:    "$1;1;;;22,0;22,3;21,8;20,9;17,0;21,4;21,8;;50;50;47;48;60;52;47;;;;;;0\r\n",
			expected: "%s: (0.0,22.0,22.3,21.8,20.9,17.0,21.4,21.8)°C (0,50,50,47,48,60,52,47)%% - Hybrid Sensor: 0.0°C, 0%%, 0.0km/h, isRaining: false (0)",
			expectedDatapoint: DataPoint{
				Temperature: [8]float64{
					0.0, 22.0, 22.3, 21.8, 20.9, 17.0, 21.4, 21.8,
				},
				Humidity: [8]float64{
					0, 50, 50, 47, 48, 60, 52, 47,
				},
				HybridSensor: HybridSensor{
					Temperature:   0.0,
					Humidity:      0,
					WindSpeed:     0,
					Precipitation: 0,
					IsRaining:     false,
				},
			},
		},
	} {
		data, err := mock.parse([]byte(cs.input))

		if err != nil {
			t.Fatalf("parsing failed: %s", err)
		}

		if data.String() != fmt.Sprintf(cs.expected, data.TimeStamp.Format(time.RFC1123)) {
			t.Fatalf("unexpected String() method result, want %s, have %s", data.String(),
				fmt.Sprintf(cs.expected, data.TimeStamp.Format(time.RFC1123)))
		}

		cs.expectedDatapoint.TimeStamp = data.TimeStamp
		if !reflect.DeepEqual(*data, cs.expectedDatapoint) {
			t.Fatalf("unexpected data received, want %+v, have %+v", cs.expectedDatapoint, *data)
		}
	}
}

package oura

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var sleepTestCases = []struct {
	name        string
	start       string
	end         string
	expectedURL string
	mock        string
}{
	{
		name:        "get sleep without specific dates",
		start:       "",
		end:         "",
		expectedURL: "/v1/sleep",
		mock: `{
			"sleep": [{
				"summary_date": "2017-11-05",
				"period_id": 0,
				"is_longest": 1,
				"timezone": 120,
				"bedtime_start": "2017-11-06T02:13:19+02:00",
				"bedtime_end": "2017-11-06T08:12:19+02:00",
				"score": 70,
				"score_total": 57,
				"score_disturbances": 83,
				"score_efficiency": 99,
				"score_latency": 88,
				"score_rem": 97,
				"score_deep": 59,
				"score_alignment": 31,
				"total": 20310,
				"duration": 21540,
				"awake": 1230,
				"light": 10260,
				"rem": 7140,
				"deep": 2910,
				"onset_latency": 480,
				"restless": 39,
				"efficiency": 94,
				"midpoint_time": 11010,
				"hr_lowest": 49,
				"hr_average": 56.375,
				"rmssd": 54,
				"breath_average": 13,
				"temperature_delta": -0.06,
				"hypnogram_5min": "443432222211222333321112222222222111133333322221112233333333332232222334",
				"hr_5min": [0, 53, 51, 0, 50, 50, 49, 49, 50, 50, 51, 52, 52, 51, 53, 58, 60, 60, 59, 58, 58, 58, 58, 55, 55, 55, 55, 56, 56, 55, 53, 53, 53, 53, 53, 53, 57, 58, 60, 60, 59, 57, 59, 58, 56, 56, 56, 56, 55, 55, 56, 56, 57, 58, 55, 56, 57, 60, 58, 58, 59, 57, 54, 54, 53, 52, 52, 55, 53, 54, 56, 0],
				"rmssd_5min": [0, 0, 62, 0, 75, 52, 56, 56, 64, 57, 55, 78, 77, 83, 70, 35, 21, 25, 49, 44, 48, 48, 62, 69, 66, 64, 79, 59, 67, 66, 70, 63, 53, 57, 53, 57, 38, 26, 18, 24, 30, 35, 36, 46, 53, 59, 50, 50, 53, 53, 57, 52, 41, 37, 49, 47, 48, 35, 32, 34, 52, 57, 62, 57, 70, 81, 81, 65, 69, 72, 64, 0]
			}]
		}`,
	},
	{
		name:        "get sleep with only start date",
		start:       "2020-01-20",
		end:         "",
		expectedURL: "/v1/sleep?start=2020-01-20",
		mock: `{
			"sleep": [
				{
					"summary_date": "2020-01-20",
					"duration": 21540
				},
				{
					"summary_date": "2020-01-21",
					"duration": 21541
				},
				{
					"summary_date": "2020-01-22",
					"duration": 21541
				},
				{
					"summary_date": "2020-01-23",
					"duration": 21541
				}
			]
		}`,
	},
	{
		name:        "get sleep with start and end dates",
		start:       "2020-01-20",
		end:         "2020-01-22",
		expectedURL: "/v1/sleep?end=2020-01-22&start=2020-01-20",
		mock: `{
			"sleep": [
				{
					"summary_date": "2020-01-20",
					"duration": 21540
				},
				{
					"summary_date": "2020-01-21",
					"duration": 21541
				},
				{
					"summary_date": "2020-01-22",
					"duration": 21541
				}
			]
		}`,
	},
}

func TestGetSleep(t *testing.T) {
	for _, tc := range sleepTestCases {
		t.Run(tc.name, func(st *testing.T) {
			testGetSleep(st, tc.start, tc.end, tc.expectedURL, tc.mock)
		})
	}
}

func testGetSleep(t *testing.T, start, end, expectedURL, mock string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v1/sleep", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, expectedURL, r.URL.String())
		fmt.Fprint(w, mock)
	})

	got, _, err := client.GetSleep(context.Background(), start, end)
	assert.NoError(t, err, "should not return an error")

	want := &Sleeps{}
	json.Unmarshal([]byte(mock), want) //nolint:errcheck

	assert.ObjectsAreEqual(want, got)
}

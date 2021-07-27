package types

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type FullWeatherReport struct {
	CityName string  `json:"city_name"`
	Data     []*Stat `json:"data"`
	Count    int     `json:"count"`
}

type Stat struct {
	PartOfDay        string  `json:"pod"`
	CityName         string  `json:"city_name"`
	DateTime         string  `json:"datetime"`
	WindDirection    string  `json:"wind_cdir"`
	SunriseTime      string  `json:"sunrise"`
	SunsetTime       string  `json:"sunset"`
	RelativeHumidity float64 `json:"rh"`
	WindSpeedMs      float64 `json:"wind_spd"`
	IndexUV          float64 `json:"uv"`
	Precipitation    float64 `json:"precip"`
	PressureMb       float64 `json:"pres"` //pressure in millibar
	Temperature      float64 `json:"temp"`
	FeelsLikeTemp    float64 `json:"app_temp"`
	HighTemp         float64 `json:"high_temp"`
	LowTemp          float64 `json:"low_temp"`
	Weather          `json:"weather"`
	CloudCoverage    int `json:"clouds"`
	Snow             int `json:"snow"`
	IndexAirQuality  int `json:"aqi"`
}

type Weather struct {
	Description string `json:"description"`
	Code        int    `json:"code"`
}

// ParseWeather parses JSON into FullWeatherReport,	which then can be used to retrieve weather information.
func ParseWeather(weather []byte) (*FullWeatherReport, error) {
	data := FullWeatherReport{}
	err := json.Unmarshal(weather, &data)
	if err != nil {
		return &FullWeatherReport{}, errors.Wrap(err, "cannot unmarshal")
	}

	return &data, nil
}

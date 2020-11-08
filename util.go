package main

import (
	"encoding/json"
	"fmt"
	es "github.com/pkg/errors"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
)

type FullWeatherReport struct {
	Data     []Stat `json:"data"`
	Count    int    `json:"count"`
	CityName string `json:"city_name"`
}

type Stat struct {
	RelativeHumidity float64 `json:"rh"`
	PartOfDay        string  `json:"pod"`
	PressureMb       float64 `json:"pres"` //pressure in millibar
	CloudCoverage    int     `json:"clouds"`
	CityName         string  `json:"city_name"`
	DateTime         string  `json:"datetime"`
	WindSpeedMs      float64 `json:"wind_spd"`
	WindDirection    string  `json:"wind_cdir"`
	SunsetTime       string  `json:"sunset"`
	Snow             int     `json:"snow"`
	IndexUV          float64 `json:"uv"`
	Precipitation    float64 `json:"precip"`
	SunriseTime      string  `json:"sunrise"`
	IndexAirQuality  int     `json:"aqi"`
	Weather          `json:"weather"`
	Temperature      float64 `json:"temp"`
	FeelsLikeTemp    float64 `json:"app_temp"`
	HighTemp         float64 `json:"high_temp"`
	LowTemp          float64 `json:"low_temp"`
}

type Weather struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

func formatCity(city string) string {
	/*Formats City into a string
	NB: only to be used inside formatNow
	*/
	fmtCity := fmt.Sprintf(
		"Weather in %v:\n", city,
	)
	return fmtCity
}
func AddNewLine(data Stat) string {
	/*Adds a new line
	 */
	_ = data
	return "\n"
}
func formatAqi(data Stat) string {
	/*Formats AQI into a string,
	adding an emoji and a comment
	*/
	fmtAqi := ""
	prefix := "Air Quality Index: "
	comment := "Unknown"
	emoji := emojis["Question"]
	var aqi = data.IndexAirQuality
	switch {
	case aqi <= 50:
		comment = "Good"
		emoji = emojis["Check"]
	case aqi > 50 && aqi <= 100:
		comment = "Moderate"
		emoji = emojis["Check"]
	case aqi > 101 && aqi <= 150:
		comment = "Unhealthy for Sensitive Groups"
		emoji = emojis["Smoking"]
	case aqi > 151 && aqi <= 200:
		comment = "Unhealthy"
		emoji = emojis["Warning"]
	case aqi > 201 && aqi <= 300:
		comment = "Very Unhealthy"
		emoji = emojis["Warning"]
	case aqi > 301 && aqi <= 500:
		comment = "Hazardous"
		emoji = emojis["Cross"]
	}
	fmtAqi = fmt.Sprintf(
		"%v\n%c %v %v", prefix, emoji, aqi, comment,
	)
	return fmtAqi
}
func formatHumidity(data Stat) string {
	/*Formats humidity into a string,
	adding an emoji and a comment
	*/
	fmtHumidity := ""
	humidity := data.RelativeHumidity
	emoji := emojis["Check"]
	comment := "Comfortable"
	switch {
	case humidity < 30:
		emoji = emojis["Dry"]
		comment = "Dry"
	case humidity > 50:
		emoji = emojis["Wet"]
		comment = "Wet"
	}
	fmtHumidity = fmt.Sprintf(
		"Humidity:\n%c %v (%v%%)", emoji, comment, math.Round(humidity),
	)
	return fmtHumidity
}
func formatPressure(data Stat) string {
	/*Formats pressure into a string,
	adding an emoji and a sign
	*/
	fmtPressure := ""
	pressure := data.PressureMb
	norm := 1013.25
	emoji := emojis["Check"]
	switch {
	case pressure < norm:
		emoji = emojis["Down"]
	case pressure > norm:
		emoji = emojis["Up"]
	}
	fmtPressure = fmt.Sprintf(
		"Pressure:\n%c %v mb (%v mmHg)", emoji, math.Round(pressure), millibarsToMmHg(math.Round(pressure)),
	)
	return fmtPressure
}
func formatSun(data Stat) string {
	/*Formats sun related data into a string
	 */
	fmtSun := fmt.Sprintf(
		"Sunrise: %v\nSunset: %v", data.SunriseTime, data.SunsetTime,
	)
	return fmtSun
}
func formatTempAndFeels(data Stat) string {
	/*Formats actual temperature together with apparent one into a string,
	adding an emoji and a sign
	*/
	fmtTempAndFeels := ""
	temperature := data.Temperature
	tempString := fmt.Sprint(temperature)
	if temperature > 0 {
		tempString = fmt.Sprintf("+%v", temperature)
	}
	feels := data.FeelsLikeTemp
	feelsString := fmt.Sprint(feels)
	if feels > 0 {
		feelsString = fmt.Sprintf("+%v", feels)
	}
	fmtTempAndFeels = fmt.Sprintf(
		"%c %v (%v) ", emojis["Thermometer"], tempString, feelsString,
	)
	return fmtTempAndFeels
}
func formatTempLowHigh(data Stat) string {
	/*Formats mean low and mean high day temperature into a string,
	adding emojis and signs
	*/
	fmtTempLowHigh := ""
	tempLow := data.LowTemp
	tempLowString := fmt.Sprint(tempLow)
	if tempLow > 0 {
		tempLowString = fmt.Sprintf("+%v", tempLow)
	}
	tempHigh := data.HighTemp
	tempHighString := fmt.Sprint(tempHigh)
	if tempHigh > 0 {
		tempHighString = fmt.Sprintf("+%v", tempHigh)
	}
	fmtTempLowHigh = fmt.Sprintf(
		"%c %v %c %v ", emojis["High"], tempHighString, emojis["Low"], tempLowString,
	)
	return fmtTempLowHigh
}
func formatTempSmall(data Stat) string {
	/*Formats Temperature into a string,
	adding a sign
	*/
	fmtTemperature := ""
	temperature := data.Temperature
	tempString := fmt.Sprint(temperature)
	if temperature > 0 {
		tempString = fmt.Sprintf("+%v", temperature)
	}
	fmtTemperature = fmt.Sprintf(
		"%v ", tempString,
	)
	return fmtTemperature
}
func formatDate(data Stat) string {
	/*Formats date into a string
	 */
	fmtDate := ""
	monthCode := data.DateTime[5:7]
	month := monthsEn[monthCode]
	date := data.DateTime[8:10]
	if date[0] == '0' {
		date = string(date[1])
	}
	fmtDate = fmt.Sprintf(
		"%v %v:", date, month,
	)
	return fmtDate
}
func formatTime(data Stat) string {
	/*Formats time into a string
	 */
	fmtTime := ""
	time := data.DateTime[11:]
	fmtTime = fmt.Sprintf(
		"%vh ", time,
	)
	return fmtTime
}
func formatUv(data Stat) string {
	/*Formats UV into a string,
	adding an emoji and a piece of advice
	*/
	if data.PartOfDay != "d" {
		return ""
	}
	fmtUv := ""
	prefix := "UV Index:"
	uv := data.IndexUV
	switch {
	case uv <= 2:
		fmtUv = fmt.Sprintf(
			"%v\n%v %v", prefix, commentsEn["no"], commentsEn["hat"],
		)
	case uv > 2 && uv <= 5:
		fmtUv = fmt.Sprintf(
			"%v\n%v %v %v %v 40 min.", prefix, commentsEn["little"], commentsEn["hat"], commentsEn["spf15"],
			commentsEn["sunburn"],
		)
	case uv > 5 && uv <= 7:
		fmtUv = fmt.Sprintf(
			"%v\n%v %v %v %v %v 30 min.", prefix, commentsEn["high"], commentsEn["hat"], commentsEn["spf30"],
			commentsEn["cover"], commentsEn["sunburn"],
		)
	case uv > 7 && uv <= 10:
		fmtUv = fmt.Sprintf(
			"%v\n%v %v %v %v %v 20 min.", prefix, commentsEn["vhigh"], commentsEn["hat"], commentsEn["spf50"],
			commentsEn["cover"], commentsEn["sunburn"],
		)
	case uv > 10:
		fmtUv = fmt.Sprintf(
			"%v\n%v %v 20 min.", prefix, commentsEn["extreme"], commentsEn["sunburn"],
		)
	}
	return fmtUv
}
func formatWeatherCode(data Stat) string {
	/*Formats WeatherCode into a string,
	adding an emoji and a description
	*/
	fmtWeatherCode := ""
	emoji := emojis["Question"]
	code := data.Weather.Code
	switch code {
	case 200, 201, 202:
		emoji = emojis["ThunderRain"]
	case 230, 231, 232, 233:
		emoji = emojis["Thunder"]
	case 300, 301, 302:
		emoji = emojis["Umbrella"]
	case 500, 501:
		emoji = emojis["UmbrellaRain"]
	case 502, 511, 520, 521, 522:
		emoji = emojis["CloudRain"]
	case 600, 601, 602, 610, 621, 622, 623:
		emoji = emojis["CloudSnow"]
	case 611, 612:
		emoji = emojis["Saturn"]
	case 700, 711, 721, 731, 741, 751:
		emoji = emojis["Fog"]
	case 800:
		emoji = emojis["Sun"]
	case 801, 802:
		emoji = emojis["SunSCloud"]
	case 803:
		emoji = emojis["SunMCloud"]
	case 804:
		emoji = emojis["Cloud"]
	case 900:
		emoji = emojis["Comet"]
	}
	fmtWeatherCode = fmt.Sprintf(
		"%c %v ", emoji, data.Weather.Description,
	)
	return fmtWeatherCode
}
func formatWind(data Stat) string {
	/*Formats wind direction data into a string,
	adding emojis and speed information
	*/
	fmtWind := ""
	direction := data.WindDirection
	speed := data.WindSpeedMs
	emoji := emojis["Balloon"]
	code, ok := emojis[direction]
	if ok {
		emoji = code
	}
	fmtWind = fmt.Sprintf(
		"%c  Wind:\n%c %v %v m/sec", emojis["Wind"], emoji, direction, math.Round(speed),
	)
	return fmtWind
}

var lineFormatterNow = []func(data Stat) string{
	formatDate, AddNewLine,
	formatTempAndFeels, AddNewLine,
	formatWeatherCode, AddNewLine,
	formatWind, AddNewLine,
	formatHumidity, AddNewLine,
	formatPressure, AddNewLine,
	formatSun, AddNewLine,
	formatUv, AddNewLine,
	formatAqi, AddNewLine,
}

var lineFormatterHours = []func(data Stat) string{
	formatTime, formatTempSmall, formatWeatherCode, AddNewLine,
}

var lineFormatterDays = []func(data Stat) string{
	formatDate, AddNewLine,
	formatTempLowHigh, formatWeatherCode, AddNewLine,
}

func (wr *FullWeatherReport) formatNow() string {
	fmt.Println("EXECUTING: formatNow")
	day := 0
	var res string
	res += formatCity(wr.Data[0].CityName)
	for _, formatter := range lineFormatterNow {
		res += formatter(wr.Data[day])
	}
	return res
}

func (wr *FullWeatherReport) formatHours(hours int) string {
	fmt.Println("EXECUTING: formatHours")
	var res string
	res += formatCity(wr.CityName)
	prevWeather := 0
	prevDate := ""
	currentDate := ""
	for i := 0; i < hours; i++ {
		if prevWeather == wr.Data[i].Code {
			continue
		}
		currentDate = formatDate(wr.Data[i])
		if prevDate != currentDate {
			res += currentDate + "\n"
		}
		for _, formatter := range lineFormatterHours {
			res += formatter(wr.Data[i])
		}
		prevWeather = wr.Data[i].Code
		prevDate = currentDate
	}
	return res
}

func (wr *FullWeatherReport) formatDays(days int) string {
	fmt.Println("EXECUTING: formatDays")
	var res string
	res += formatCity(wr.CityName)
	for i := 0; i < days; i++ {
		for _, formatter := range lineFormatterDays {
			res += formatter(wr.Data[i])
		}
	}
	return res
}

func getForecast(loc *tgbotapi.Location, period string) ([]byte, error) {
	/*
	 */
	fmt.Println("EXECUTING: getForecast")
	cfg := parseConfig()
	lat := fmt.Sprint(loc.Latitude)
	long := fmt.Sprint(loc.Longitude)
	baseWeatherURL := "https://weatherbit-v1-mashape.p.rapidapi.com/"
	url := fmt.Sprintf(
		"%v%vlang=%v&lat=%v&lon=%v", baseWeatherURL, forecasts[period], cfg.Language, lat, long,
	)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-rapidapi-host", "weatherbit-v1-mashape.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", cfg.WeatherAPI)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		err = es.Wrap(err, "failed to perform request")
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = es.Wrap(err, "failed to read from response body")
		return nil, err
	}
	fmt.Println("BODY:", string(body))
	return body, nil
}

func parseWeather(weather []byte) (FullWeatherReport, error) {
	/*Parses JSON into FullWeatherReport,
	which then can be used to retrieve weather information
	*/
	fmt.Println("EXECUTING: parseWeather")
	data := FullWeatherReport{}
	err := json.Unmarshal(weather, &data)
	if err != nil {
		err = es.Wrap(err, "failed to unmarshal")
		return data, err
	}
	return data, nil
}

func millibarsToMmHg(millibars float64) float64 {
	/*Converting pressure in millibars into millimeters of mercury
	 */
	return millibars * 0.75
}

func pickASaying(sayings []string) string {
	/*Picking a random saying
	 */
	randomIndex := rand.Intn(len(sayings))
	return sayings[randomIndex]
}

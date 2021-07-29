package service

import (
	"fmt"
	"math"
	"weather-or-not-bot/internal/types"
)

const (
	// Constants unlikely to be changed.
	AQIModerateThreshold              = 50
	AQIUnhealthyForSensitiveThreshold = 100
	AQIUnhealthyThreshold             = 150
	AQIVeryUnhealthyThreshold         = 200
	AQIHazardousThreshold             = 300
	AQIUpperBound                     = 500

	HumidityDryThreshold = 30
	HumidityWetThreshold = 50

	PressureNormalAtmospheric float64 = 1013.25

	Zero = 0

	Eastern           = "E"
	EastNorthEastern  = "ENE"
	EastSouthEastern  = "ESE"
	Northern          = "N"
	NorthEastern      = "NE"
	NorthNorthEastern = "NNE"
	NorthNorthWestern = "NNW"
	NorthWestern      = "NW"
	Southern          = "S"
	SouthEastern      = "SE"
	SouthSouthEastern = "SSE"
	SouthSouthWestern = "SSW"
	SouthWestern      = "SW"
	Western           = "W"
	WestNorthWestern  = "WNW"
	WestSouthWestern  = "WSW"

	NewLine = "\n"

	// Language specific strings, subject to localization.
	// TODO localize to other languages and move out to a repo
	AQIPrefix                = "Air Quality Index"
	AQIUnknown               = "Unknown"
	AQIGood                  = "Good"
	AQIModerate              = "Moderate"
	AQIUnhealthyForSensitive = "Unhealthy for Sensitive Groups"
	AQIUnhealthy             = "Unhealthy"
	AQIVeryUnhealthy         = "Very unhealthy"
	AQIHazardous             = "Hazardous"

	HumidityPrefix      = "Humidity"
	HumidityComfortable = "Comfortable"
	HumidityDry         = "Dry"
	HumidityWet         = "Wet"

	PressurePrefix                   = "Pressure"
	PressureMillibars                = "mb"
	PressureMillimetersOfQuicksilver = "mmHg"

	SunSunrise = "Sunrise"
	SunSunset  = "Sunset"

	WeatherCityPrefix = "Weather in"

	Day      = "d"
	UVPrefix = "UV Index"
)

type StatFormatter func(s *types.Stat) string

// addNewLine adds a new line
func addNewLine(s *types.Stat) string {
	_ = s
	return NewLine
}

// formatAQI formats AQI into a string, adding an emoji and a comment.
func formatAQI(s *types.Stat) string {
	var (
		emoji   types.Emoji
		comment string
	)

	switch {
	case s.IndexAirQuality <= AQIModerateThreshold:
		comment = AQIGood
		emoji = types.CheckMarkEmoji
	case s.IndexAirQuality <= AQIUnhealthyForSensitiveThreshold:
		comment = AQIModerate
		emoji = types.CheckMarkEmoji
	case s.IndexAirQuality <= AQIUnhealthyThreshold:
		comment = AQIUnhealthyForSensitive
		emoji = types.SmokingEmoji
	case s.IndexAirQuality <= AQIVeryUnhealthyThreshold:
		comment = AQIUnhealthy
		emoji = types.WarningEmoji
	case s.IndexAirQuality <= AQIHazardousThreshold:
		comment = AQIVeryUnhealthy
		emoji = types.WarningEmoji
	case s.IndexAirQuality <= AQIUpperBound:
		comment = AQIHazardous
		emoji = types.CrossEmoji
	default:
		comment = AQIUnknown
		emoji = types.QuestionMarkEmoji
	}

	return fmt.Sprintf("%s:\n%c %v %s", AQIPrefix, emoji, s.IndexAirQuality, comment)
}

// formatHumidity formats humidity into a string, adding an emoji and a comment.
func formatHumidity(s *types.Stat) string {
	var (
		emoji   types.Emoji
		comment string
	)

	switch {
	case s.RelativeHumidity < HumidityDryThreshold:
		emoji = types.DryEmoji
		comment = HumidityDry
	case s.RelativeHumidity > HumidityWetThreshold:
		emoji = types.WaterDropEmoji
		comment = HumidityWet
	default:
		emoji = types.CheckMarkEmoji
		comment = HumidityComfortable
	}

	return fmt.Sprintf("%s:\n%c %s (%v%%)", HumidityPrefix, emoji, comment, math.Round(s.RelativeHumidity))
}

// formatPressure formats pressure into a string, adding an emoji and a sign.
func formatPressure(s *types.Stat) string {
	var emoji types.Emoji

	switch {
	case s.PressureMb < PressureNormalAtmospheric:
		emoji = types.DownwardTriangleEmoji
	case s.PressureMb > PressureNormalAtmospheric:
		emoji = types.UpwardTriangleEmoji
	default:
		emoji = types.CheckMarkEmoji
	}

	return fmt.Sprintf("%s:\n%c %v %s (%v %s)",
		PressurePrefix, emoji, math.Round(s.PressureMb), PressureMillibars, millibarsToMmHg(math.Round(s.PressureMb)), PressureMillimetersOfQuicksilver)
}

// formatSun formats sun related data into a string.
func formatSun(s *types.Stat) string {
	return fmt.Sprintf("%s: %v\n%s: %v", SunSunrise, s.SunriseTime, SunSunset, s.SunsetTime)
}

// formatTemperatureAndFeels formats actual temperature together with apparent one into a string, adding an emoji and a sign.
func formatTemperatureAndFeels(s *types.Stat) string {
	temp := fmt.Sprint(s.Temperature)
	if s.Temperature > Zero {
		temp = fmt.Sprintf("+%v", s.Temperature)
	}

	feels := fmt.Sprint(s.FeelsLikeTemp)
	if s.FeelsLikeTemp > Zero {
		feels = fmt.Sprintf("+%v", s.FeelsLikeTemp)
	}

	return fmt.Sprintf("%c %s (%s) ", types.ThermometerEmoji, temp, feels)
}

// formatTemperatureLowHigh formats mean low and mean high day temperature into a string, adding emojis and signs.
func formatTemperatureLowHigh(s *types.Stat) string {
	tempLow := fmt.Sprint(s.LowTemp)
	if s.LowTemp > 0 {
		tempLow = fmt.Sprintf("+%v", s.LowTemp)
	}

	tempHigh := fmt.Sprint(s.HighTemp)
	if s.HighTemp > 0 {
		tempHigh = fmt.Sprintf("+%v", s.HighTemp)
	}

	return fmt.Sprintf("%c %v %c %v ", types.GoingUpEmoji, tempHigh, types.GoingDownEmoji, tempLow)
}

// formatTemperatureSmall formats Temperature into a string, adding a sign.
func formatTemperatureSmall(s *types.Stat) string {
	if s.Temperature > 0 {
		return fmt.Sprintf("+%v ", s.Temperature)
	}

	return fmt.Sprintf("%v ", s.Temperature)
}

// formatDate formats date into a string.
func formatDate(s *types.Stat) string {
	// TODO format date properly
	monthCode := s.DateTime[5:7]
	month := monthsEn[monthCode]
	date := s.DateTime[8:10]
	if date[0] == '0' {
		date = string(date[1])
	}
	return fmt.Sprintf("%v %v:", date, month)
}

// formatTime formats time into a string.
func formatTime(s *types.Stat) string {
	// TODO format time properly
	return fmt.Sprintf("%vh ", s.DateTime[11:])
}

// formatUv formats UV into a string, adding an emoji and a piece of advice.
func formatUv(s *types.Stat) string {
	if s.PartOfDay != Day {
		return ""
	}

	switch {
	case s.IndexUV <= 2:
		return fmt.Sprintf(
			"%s:\n%v %v", UVPrefix, commentsEn["no"], commentsEn["hat"],
		)
	case s.IndexUV > 2 && s.IndexUV <= 5:
		return fmt.Sprintf(
			"%v\n%v %v %v %v 40 min.", UVPrefix, commentsEn["little"], commentsEn["hat"], commentsEn["spf15"],
			commentsEn["sunburn"],
		)
	case s.IndexUV > 5 && s.IndexUV <= 7:
		return fmt.Sprintf(
			"%v\n%v %v %v %v %v 30 min.", UVPrefix, commentsEn["high"], commentsEn["hat"], commentsEn["spf30"],
			commentsEn["cover"], commentsEn["sunburn"],
		)
	case s.IndexUV > 7 && s.IndexUV <= 10:
		return fmt.Sprintf(
			"%v\n%v %v %v %v %v 20 min.", UVPrefix, commentsEn["vhigh"], commentsEn["hat"], commentsEn["spf50"],
			commentsEn["cover"], commentsEn["sunburn"],
		)
	case s.IndexUV > 10:
		return fmt.Sprintf(
			"%v\n%v %v 20 min.", UVPrefix, commentsEn["extreme"], commentsEn["sunburn"],
		)
	default:
		return ""
	}
}

// formatWeatherCode formats WeatherCode into a string, adding an emoji and a description.
func formatWeatherCode(s *types.Stat) string {
	var emoji types.Emoji

	//TODO replace codes with constants
	switch s.Weather.Code {
	case 200, 201, 202:
		emoji = types.ThunderRainEmoji
	case 230, 231, 232, 233:
		emoji = types.ThunderEmoji
	case 300, 301, 302:
		emoji = types.UmbrellaEmoji
	case 500, 501:
		emoji = types.UmbrellaRainEmoji
	case 502, 511, 520, 521, 522:
		emoji = types.CloudRainEmoji
	case 600, 601, 602, 610, 621, 622, 623:
		emoji = types.CloudSnowEmoji
	case 611, 612:
		emoji = types.SaturnEmoji
	case 700, 711, 721, 731, 741, 751:
		emoji = types.FogEmoji
	case 800:
		emoji = types.SunEmoji
	case 801, 802:
		emoji = types.SunWithSmallCloudEmoji
	case 803:
		emoji = types.SunWithMediumCloudEmoji
	case 804:
		emoji = types.CloudEmoji
	case 900:
		emoji = types.CometEmoji
	default:
		emoji = types.QuestionMarkEmoji
	}

	return fmt.Sprintf("%c %s ", emoji, s.Weather.Description)
}

// formatWind formats wind direction data into a string, adding emojis and speed information.
func formatWind(s *types.Stat) string {
	var emoji types.Emoji

	switch s.WindDirection {
	case Northern:
		emoji = types.NorthEmoji
	case NorthEastern, NorthNorthEastern, EastNorthEastern:
		emoji = types.NorthEastEmoji
	case NorthWestern, NorthNorthWestern, WestNorthWestern:
		emoji = types.NorthWestEmoji
	case Eastern:
		emoji = types.EastEmoji
	case Western:
		emoji = types.WestEmoji
	case Southern:
		emoji = types.SouthEmoji
	case SouthEastern, SouthSouthEastern, EastSouthEastern:
		emoji = types.SouthEastEmoji
	case SouthWestern, SouthSouthWestern, WestSouthWestern:
		emoji = types.SouthWestEmoji
	default:
		emoji = types.BalloonEmoji
	}

	return fmt.Sprintf("%c  Wind:\n%c %v %v m/sec", types.WindEmoji, emoji, s.WindDirection, math.Round(s.WindSpeedMs))
}

// formatCity formats a city into a string.
// NB: only to be used inside FormatNow
func formatCity(city string) string {
	return fmt.Sprintf("%s %s:\n", WeatherCityPrefix, city)
}

// millibarsToMmHg converts pressure in millibars into millimeters of mercury.
func millibarsToMmHg(millibars float64) float64 {
	return millibars * 0.75
}

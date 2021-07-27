package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"weather-or-not-bot/internal/types"
)

var (
	lineFormatterNow = []StatFormatter{
		formatDate, addNewLine,
		formatTemperatureAndFeels, addNewLine,
		formatWeatherCode, addNewLine,
		formatWind, addNewLine,
		formatHumidity, addNewLine,
		formatPressure, addNewLine,
		formatSun, addNewLine,
		formatUv, addNewLine,
		formatAQI, addNewLine,
	}

	lineFormatterHours = []StatFormatter{
		formatTime, formatTemperatureSmall, formatWeatherCode, addNewLine,
	}

	lineFormatterDays = []StatFormatter{
		formatDate, addNewLine,
		formatTemperatureLowHigh, formatWeatherCode, addNewLine,
	}
)

type Formatter struct {
}

func NewFormatter() *Formatter {
	return &Formatter{}
}

// FormatNow formats current weather state data.
func (f *Formatter) FormatNow(ctx context.Context, report *types.FullWeatherReport) string {
	ctxlogrus.Extract(ctx).Debug("Running format now")

	var buf bytes.Buffer

	buf.WriteString(formatCity(report.Data[0].CityName))
	for i := range lineFormatterNow {
		buf.WriteString(lineFormatterNow[i](report.Data[0]))
	}

	return buf.String()
}

// FormatHours formats hours.
func (f *Formatter) FormatHours(ctx context.Context, report *types.FullWeatherReport, hours int) string {
	ctxlogrus.Extract(ctx).Debug("Running format hours")

	var (
		buf             bytes.Buffer
		lastWeatherCode int
		prevDate        string
		currentDate     string
	)

	buf.WriteString(formatCity(report.CityName))

	for i := 0; i < hours; i++ {
		if lastWeatherCode == report.Data[i].Code {
			continue
		}

		currentDate = formatDate(report.Data[i])
		if prevDate != currentDate {
			buf.WriteString(fmt.Sprintln(currentDate))
		}

		for i := range lineFormatterHours {
			buf.WriteString(formatHour(report.Data[i]))
		}

		lastWeatherCode = report.Data[i].Code
		prevDate = currentDate
	}

	return buf.String()
}

// FormatDays formats days.
func (f *Formatter) FormatDays(ctx context.Context, report *types.FullWeatherReport, days int) string {
	ctxlogrus.Extract(ctx).Debug("Running format days")

	var buf bytes.Buffer
	buf.WriteString(formatCity(report.CityName))
	for i := 0; i < days; i++ {
		buf.WriteString(formatDay(report.Data[i]))
	}

	return buf.String()
}

func formatHour(s *types.Stat) string {
	var buf bytes.Buffer
	for i := range lineFormatterDays {
		buf.WriteString(lineFormatterHours[i](s))
	}

	return buf.String()
}

func formatDay(s *types.Stat) string {
	var buf bytes.Buffer
	for i := range lineFormatterDays {
		buf.WriteString(lineFormatterDays[i](s))
	}

	return buf.String()
}

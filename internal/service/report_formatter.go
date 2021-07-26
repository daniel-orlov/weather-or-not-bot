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

type ReportFormatter struct {
	Report *types.FullWeatherReport
}

// FormatNow formats current weather state data.
func (f *ReportFormatter) FormatNow(ctx context.Context) string {
	ctxlogrus.Extract(ctx).Debug("Running format now")

	var buf bytes.Buffer

	buf.WriteString(formatCity(f.Report.Data[0].CityName))
	for i := range lineFormatterNow {
		buf.WriteString(lineFormatterNow[i](f.Report.Data[0]))
	}

	return buf.String()
}

// FormatHours formats hours.
func (f *ReportFormatter) FormatHours(ctx context.Context, hours int) string {
	ctxlogrus.Extract(ctx).Debug("Running format hours")

	var (
		buf             bytes.Buffer
		lastWeatherCode int
		prevDate        string
		currentDate     string
	)

	buf.WriteString(formatCity(f.Report.CityName))

	for i := 0; i < hours; i++ {
		if lastWeatherCode == f.Report.Data[i].Code {
			continue
		}

		currentDate = formatDate(f.Report.Data[i])
		if prevDate != currentDate {
			buf.WriteString(fmt.Sprintln(currentDate))
		}

		for i := range lineFormatterHours {
			buf.WriteString(formatHour(f.Report.Data[i]))
		}

		lastWeatherCode = f.Report.Data[i].Code
		prevDate = currentDate
	}

	return buf.String()
}

// FormatDays formats days.
func (f *ReportFormatter) FormatDays(ctx context.Context, days int) string {
	ctxlogrus.Extract(ctx).Debug("Running format days")

	var buf bytes.Buffer
	buf.WriteString(formatCity(f.Report.CityName))
	for i := 0; i < days; i++ {
		buf.WriteString(formatDay(f.Report.Data[i]))
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

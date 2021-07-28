package service

import (
	"context"
	"math/rand"
	"strconv"
	"strings"
	"weather-or-not-bot/internal/types"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	bot "gopkg.in/telegram-bot-api.v4"
)

type MessageService struct {
	botCmd   *bot.BotAPI
	usrRepo  UserDataRepo
	botRepo  BotUIRepo
	locRepo  LocationRepo
	forecast ForecastClient
	format   ReportFormatter
}

func NewMessageService(botCmd *bot.BotAPI, usrRepo UserDataRepo, botRepo BotUIRepo, locRepo LocationRepo, forecast ForecastClient, format ReportFormatter) *MessageService {
	return &MessageService{botCmd: botCmd, usrRepo: usrRepo, botRepo: botRepo, locRepo: locRepo, forecast: forecast, format: format}
}

const (
	Start            = "/start"
	BackToMainMenu   = "< Back"
	Back             = "<< Back"
	ByHours          = "By Hours"
	ByDays           = "By Days"
	CurrentWeather   = "Now"
	EmptyMessage     = ""
	WeatherHere      = "Weather at my location"
	WeatherElsewhere = "Weather elsewhere"
	Stop             = "/stop"

	ThreeDays          = "3 days"
	FiveDays           = "5 days"
	SevenDays          = "7 days"
	TenDays            = "10 days"
	SixteenDays        = "16 days"
	TwentyFourHours    = "24 hours"
	FortyEightHours    = "48 hours"
	SeventyTwoHours    = "72 hours"
	NinetySixHours     = "96 hours"
	HundredTwentyHours = "120 hours"

	HOURLY = true
	DAILY  = false
)

func (s *MessageService) HandleNewMessage(ctx context.Context, upd *bot.Update) error {
	log := ctxlogrus.Extract(ctx).WithFields(logrus.Fields{
		"username":     upd.Message.From.UserName,
		"user_id":      upd.Message.From.ID,
		"message_text": upd.Message.Text,
	})
	log.Info("Handling a new message")

	var err error
	switch upd.Message.Text {
	case Start:
		err = s.handleStart(ctx, upd.Message)
	case BackToMainMenu:
		err = s.handleBackToMainMenu(ctx, upd.Message)
	case Back:
		err = s.handleBack(ctx, upd.Message)
	case ByHours:
		err = s.handleByHours(ctx, upd.Message)
	case ByDays:
		err = s.handleByDays(ctx, upd.Message)
	case CurrentWeather:
		err = s.handleNow(ctx, upd.Message)
	case ThreeDays, FiveDays, SevenDays, TenDays, SixteenDays:
		err = s.handlePeriod(ctx, upd.Message, HOURLY)
	case TwentyFourHours, FortyEightHours, SeventyTwoHours, NinetySixHours, HundredTwentyHours:
		err = s.handlePeriod(ctx, upd.Message, DAILY)
	case WeatherHere:
		err = s.handleLocationByCoordinates(ctx, upd.Message)
	case WeatherElsewhere:
		err = s.handleWeatherElsewhere(ctx, upd.Message)
	case EmptyMessage:
		err = s.handleEmptyMessage(ctx, upd.Message)
	case Stop:
		err = s.handleStop(ctx, upd.Message)
	default:
		err = s.handleUnknown(ctx, upd.Message)
	}

	if err != nil {
		return errors.Wrap(err, "cannot handle a new message")
	}

	return nil
}

func (s *MessageService) handleStop(ctx context.Context, req *bot.Message) error {
	ctxlogrus.Extract(ctx).Infof("Handling '%s'", req.Text)

	resp := bot.NewMessage(req.Chat.ID, commentsEn["End"])
	resp.ReplyMarkup = bot.ReplyKeyboardHide{HideKeyboard: true}

	_, err := s.botCmd.Send(resp)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	return nil
}

func (s *MessageService) handleStart(ctx context.Context, req *bot.Message) error {
	ctxlogrus.Extract(ctx).Infof("Handling '%s'", req.Text)

	err := s.usrRepo.AddUserIfNotExists(ctx, req.From)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	// TODO create func greet()
	resp := bot.NewMessage(req.Chat.ID, commentsEn["DefaultMessage"]+"\n"+pickASaying(sayingsEn))
	resp.ReplyMarkup = s.botRepo.GetMainMenuKeyboard()

	_, err = s.botCmd.Send(resp)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	return nil
}

func (s *MessageService) handleBackToMainMenu(ctx context.Context, req *bot.Message) error {
	ctxlogrus.Extract(ctx).Infof("Handling '%s'", req.Text)

	resp := bot.NewMessage(req.Chat.ID, commentsEn["ChooseLocation"])
	resp.ReplyMarkup = s.botRepo.GetMainMenuKeyboard()

	_, err := s.botCmd.Send(resp)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	return nil
}

func (s *MessageService) handleBack(ctx context.Context, req *bot.Message) error {
	ctxlogrus.Extract(ctx).Infof("Handling '%s'", req.Text)

	resp := bot.NewMessage(req.Chat.ID, commentsEn["ChoosePeriodType"])
	resp.ReplyMarkup = s.botRepo.GetDaysOrHoursKeyboard()

	_, err := s.botCmd.Send(resp)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	return nil
}

func (s *MessageService) handleByHours(ctx context.Context, req *bot.Message) error {
	ctxlogrus.Extract(ctx).Infof("Handling '%s'", req.Text)

	resp := bot.NewMessage(req.Chat.ID, commentsEn["ChoosePeriod"])
	resp.ReplyMarkup = s.botRepo.GetHoursKeyboard()

	_, err := s.botCmd.Send(resp)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	return nil
}

func (s *MessageService) handleByDays(ctx context.Context, req *bot.Message) error {
	ctxlogrus.Extract(ctx).Infof("Handling '%s'", req.Text)

	resp := bot.NewMessage(req.Chat.ID, commentsEn["ChoosePeriod"])
	resp.ReplyMarkup = s.botRepo.GetDaysKeyboard()

	_, err := s.botCmd.Send(resp)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	return nil
}

func (s *MessageService) handleNow(ctx context.Context, req *bot.Message) error {
	log := ctxlogrus.Extract(ctx)
	log.Infof("Handling '%s'", req.Text)

	loc, err := s.locRepo.GetUserRecentLocation(ctx, req.From.ID)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	forecast, err := s.forecast.GetForecast(ctx, loc, req.Text)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	wr, err := types.ParseWeather(forecast)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	resp := bot.NewMessage(req.Chat.ID, s.format.FormatNow(ctx, wr))
	resp.ReplyMarkup = s.botRepo.GetDaysOrHoursKeyboard()

	_, err = s.botCmd.Send(resp)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	err = s.locRepo.SaveLocationName(ctx, req.From.ID, wr.Data[0].CityName)
	if err != nil {
		log.WithError(err).Warnf("cannot save location '%s'", wr.Data[0].CityName)
	}

	return nil
}

func (s *MessageService) handlePeriod(ctx context.Context, req *bot.Message, byHours bool) error {
	log := ctxlogrus.Extract(ctx)
	log.Infof("Handling '%s'", req.Text)

	loc, err := s.locRepo.GetUserRecentLocation(ctx, req.From.ID)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	forecast, err := s.forecast.GetForecast(ctx, loc, req.Text)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	wr, err := types.ParseWeather(forecast)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	var resp bot.MessageConfig
	if byHours {
		resp = bot.NewMessage(req.Chat.ID, s.format.FormatHours(ctx, wr, extractNumerals(req.Text)))
		resp.ReplyMarkup = s.botRepo.GetHoursKeyboard()
	} else {
		resp = bot.NewMessage(req.Chat.ID, s.format.FormatDays(ctx, wr, extractNumerals(req.Text)))
		resp.ReplyMarkup = s.botRepo.GetDaysKeyboard()
	}

	_, err = s.botCmd.Send(resp)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	err = s.locRepo.SaveLocationName(ctx, req.From.ID, wr.Data[0].CityName)
	if err != nil {
		log.WithError(err).Warnf("cannot save location '%s'", wr.Data[0].CityName)
	}

	return nil
}

func (s *MessageService) handleLocationByCoordinates(ctx context.Context, req *bot.Message) error {
	log := ctxlogrus.Extract(ctx)
	log.Debug("Handling location by coordinates")

	err := s.locRepo.AddLocationByCoordinates(ctx, req.From.ID, req.Location)
	if err != nil {
		return errors.Wrap(err, types.ErrHandlingLocByCoords)
	}

	resp := bot.NewMessage(req.Chat.ID, commentsEn["CoordsAccepted"])
	resp.ReplyMarkup = s.botRepo.GetDaysOrHoursKeyboard()

	_, err = s.botCmd.Send(resp)
	if err != nil {
		return errors.Wrap(err, types.ErrHandlingLocByCoords)
	}

	return nil
}

func (s *MessageService) handleWeatherElsewhere(ctx context.Context, req *bot.Message) error {
	log := ctxlogrus.Extract(ctx)
	log.Infof("Handling '%s'", req.Text)

	resp := bot.NewMessage(req.Chat.ID, commentsEn["DiffPlaceAccepted"])
	resp.ReplyMarkup = bot.ReplyKeyboardHide{HideKeyboard: true}

	_, err := s.botCmd.Send(resp)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	return nil
}

func (s *MessageService) handleLocationByText(ctx context.Context, req *bot.Message) error {
	log := ctxlogrus.Extract(ctx)
	log.Infof("Handling location '%s' by text ", req.Text)

	loc, err := s.locRepo.GetCoordinatesByCityName(ctx, req.Text)
	if err != nil {
		return errors.Wrapf(err, types.ErrHandlingLocByText, req.Text)
	}

	var resp bot.MessageConfig
	if loc.Latitude == 0 && loc.Longitude == 0 {
		resp = bot.NewMessage(req.Chat.ID, commentsEn["TryAgain"])
		resp.ReplyMarkup = s.botRepo.GetBackToMainMenuKeyboard()
	} else {
		resp = bot.NewMessage(req.Chat.ID, commentsEn["CoordsAccepted"])
		resp.ReplyMarkup = s.botRepo.GetDaysOrHoursKeyboard()
	}

	err = s.locRepo.AddLocationByCoordinates(ctx, req.From.ID, &loc)
	if err != nil {
		return errors.Wrapf(err, types.ErrHandlingLocByText, req.Text)
	}

	_, err = s.botCmd.Send(resp)
	if err != nil {
		return errors.Wrapf(err, types.ErrHandlingLocByText, req.Text)
	}

	return nil
}

func (s *MessageService) handleUnknown(ctx context.Context, req *bot.Message) error {
	ctxlogrus.Extract(ctx).Infof("Handling '%s'", req.Text)

	resp := bot.NewMessage(req.Chat.ID, commentsEn["Unknown"])
	resp.ReplyMarkup = s.botRepo.GetMainMenuKeyboard()

	_, err := s.botCmd.Send(resp)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	return nil
}

func (s *MessageService) handleEmptyMessage(ctx context.Context, req *bot.Message) error {
	ctxlogrus.Extract(ctx).Infof("Handling empty message")

	if req.Location != nil {
		return s.handleLocationByCoordinates(ctx, req)
	}

	return s.handleUnknown(ctx, req)
}

// Picks a random saying.
func pickASaying(sayings []string) string {
	return sayings[rand.Intn(len(sayings))]
}

func extractNumerals(s string) int {
	num, _ := strconv.Atoi(strings.Split(s, " ")[0])
	return num
}

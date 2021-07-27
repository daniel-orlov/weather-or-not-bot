package service

import (
	"context"
	"fmt"
	"math/rand"
	"weather-or-not-bot/internal/types"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	bot "gopkg.in/telegram-bot-api.v4"
)

type MessageService struct {
	botClient *bot.BotAPI
	userRepo  UserDataRepo
	botUIRepo BotUIRepo
	locaRepo  LocationRepo
	forecasts ForecastClient
	formatter Formatter
}

// TODO add constructor here

const (
	Start          = "/start"
	BackToMainMenu = "< Back"
	Back           = "<< Back"
	ByHours        = "By Hours"
	ByDays         = "By Days"
	CurrentWeather = "Now"
	Stop           = "/stop"

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
	DAILY = false
)

var handlersEn = map[string]func(um *UserMessage){
	//"/stop":                  handleStop,
	//"/start":                 handleStart,
	"":                       handleEmpty,
	"Weather at my location": handleLocationByCoords,
	"Weather elsewhere":      handleWeatherElsewhere,
	//"< Back":                 handleBackToMainMenu,
	//"<< Back":                handleBack,
	//"By Days":                handleByDays,
	//"By Hours":               handleByHours,
	//"Now":                    handleNow,
	//"3 days":    handleDays,
	//"5 days":    handleDays,
	//"7 days":    handleDays,
	//"10 days":   handleDays,
	//"16 days":   handleDays,
	//"24 hours":  handleHours,
	//"48 hours":  handleHours,
	//"72 hours":  handleHours,
	//"96 hours":  handleHours,
	//"120 hours": handleHours,
}

func (s *MessageService) HandleNewMessage(ctx context.Context, upd bot.Update) error {
	log := ctxlogrus.Extract(ctx).WithFields(logrus.Fields{
		"username":     upd.Message.From.UserName,
		"user_id":      upd.Message.From.ID,
		"message_text": upd.Message.Text,
	})
	log.Debugf("Handling a new message")

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
	case Stop:
		err = s.handleStop(ctx, upd.Message)
	}

	if err != nil {
		return errors.Wrap(err, "cannot handle a new message")
	}

	return nil
}

func (s *MessageService) handleStop(ctx context.Context, req *bot.Message) error {
	ctxlogrus.Extract(ctx).Debugf("Handling '%s'", req.Text)

	resp := bot.NewMessage(req.Chat.ID, commentsEn["End"])
	resp.ReplyMarkup = bot.ReplyKeyboardHide{HideKeyboard: true}

	_, err := s.botClient.Send(resp)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	return nil
}

func (s *MessageService) handleStart(ctx context.Context, req *bot.Message) error {
	ctxlogrus.Extract(ctx).Debugf("Handling '%s'", req.Text)

	err := s.userRepo.AddUserIfNotExists(ctx, req.From)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	// TODO create func greet()
	resp := bot.NewMessage(req.Chat.ID, commentsEn["DefaultMessage"]+"\n"+pickASaying(sayingsEn))
	resp.ReplyMarkup = s.botUIRepo.GetMainMenuKeyboard()

	_, err = s.botClient.Send(resp)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	return nil
}

func (s *MessageService) handleBackToMainMenu(ctx context.Context, req *bot.Message) error {
	ctxlogrus.Extract(ctx).Debugf("Handling '%s'", req.Text)

	resp := bot.NewMessage(req.Chat.ID, commentsEn["ChooseLocation"])
	resp.ReplyMarkup = s.botUIRepo.GetMainMenuKeyboard()

	_, err := s.botClient.Send(resp)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	return nil
}

func (s *MessageService) handleBack(ctx context.Context, req *bot.Message) error {
	ctxlogrus.Extract(ctx).Debugf("Handling '%s'", req.Text)

	resp := bot.NewMessage(req.Chat.ID, commentsEn["ChoosePeriodType"])
	resp.ReplyMarkup = s.botUIRepo.GetDaysOrHoursKeyboard()

	_, err := s.botClient.Send(resp)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	return nil
}

func (s *MessageService) handleByHours(ctx context.Context, req *bot.Message) error {
	ctxlogrus.Extract(ctx).Debugf("Handling '%s'", req.Text)

	resp := bot.NewMessage(req.Chat.ID, commentsEn["ChoosePeriod"])
	resp.ReplyMarkup = s.botUIRepo.GetHoursKeyboard()

	_, err := s.botClient.Send(resp)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	return nil
}

func (s *MessageService) handleByDays(ctx context.Context, req *bot.Message) error {
	ctxlogrus.Extract(ctx).Debugf("Handling '%s'", req.Text)

	resp := bot.NewMessage(req.Chat.ID, commentsEn["ChoosePeriod"])
	resp.ReplyMarkup = s.botUIRepo.GetDaysKeyboard()

	_, err := s.botClient.Send(resp)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	return nil
}

func (s *MessageService) handleNow(ctx context.Context, req *bot.Message) error {
	log := ctxlogrus.Extract(ctx)
	log.Debugf("Handling '%s'", req.Text)

	loc, err := s.locaRepo.GetUserRecentLocation(ctx, req.From.ID)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	forecast, err := s.forecasts.GetForecast(ctx, loc, req.Text)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	wr, err := types.ParseWeather(forecast)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	resp := bot.NewMessage(req.Chat.ID, s.formatter.FormatNow(ctx))
	resp.ReplyMarkup = s.botUIRepo.GetDaysOrHoursKeyboard()

	_, err = s.botClient.Send(resp)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	err = s.locaRepo.SaveLocationName(ctx, req.From.ID, wr.Data[0].CityName)
	if err != nil {
		log.WithError(err).Warnf("cannot save location '%s'", wr.Data[0].CityName)
	}

	return nil
}

func (s *MessageService) handlePeriod(ctx context.Context, req *bot.Message, byHours bool) error {
	log := ctxlogrus.Extract(ctx)
	log.Debugf("Handling '%s'", req.Text)

	loc, err := s.locaRepo.GetUserRecentLocation(ctx, req.From.ID)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	forecast, err := s.forecasts.GetForecast(ctx, loc, req.Text)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	wr, err := types.ParseWeather(forecast)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	var resp bot.MessageConfig
	if byHours {
		resp = bot.NewMessage(req.Chat.ID, s.formatter.FormatHours(ctx, timePeriodsEn[req.Text]))
		resp.ReplyMarkup = s.botUIRepo.GetHoursKeyboard()
	} else {
		resp = bot.NewMessage(req.Chat.ID, s.formatter.FormatDays(ctx, timePeriodsEn[req.Text]))
		resp.ReplyMarkup = s.botUIRepo.GetDaysKeyboard()
	}

	_, err = s.botClient.Send(resp)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	err = s.locaRepo.SaveLocationName(ctx, req.From.ID, wr.Data[0].CityName)
	if err != nil {
		log.WithError(err).Warnf("cannot save location '%s'", wr.Data[0].CityName)
	}

	return nil
}

func (s *MessageService) handleLocationByCoordinates(ctx context.Context, req *bot.Message) error {
	log := ctxlogrus.Extract(ctx)
	log.Debugf("Handling '%s'", req.Text)

	loc, err := s.locaRepo.GetUserRecentLocation(ctx, req.From.ID)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	forecast, err := s.forecasts.GetForecast(ctx, loc, req.Text)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	wr, err := types.ParseWeather(forecast)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	resp := bot.NewMessage(req.Chat.ID, s.formatter.FormatDays(ctx, timePeriodsEn[req.Text]))
	resp.ReplyMarkup = s.botUIRepo.GetDaysKeyboard()


	_, err = s.botClient.Send(resp)
	if err != nil {
		return errors.Wrapf(err, types.ErrOnHandling, req.Text)
	}

	err = s.locaRepo.SaveLocationName(ctx, req.From.ID, wr.Data[0].CityName)
	if err != nil {
		log.WithError(err).Warnf("cannot save location '%s'", wr.Data[0].CityName)
	}

	return nil
}


func handleLocationByCoords(um *UserMessage) {
	fmt.Println("EXECUTING: handleLocationByCoords")
	addLocationByCoords(um.connection, um.update.Message)
	msg := bot.NewMessage(um.update.Message.Chat.ID, commentsEn["CoordsAccepted"])
	msg.ReplyMarkup = keyboards["period"]
	_, err := um.bot.Send(msg)
	if err != nil {
		err = errors.Wrap(err, "Unable to Send message")
		fmt.Println(err)
	}
}

func handleWeatherElsewhere(um *UserMessage) {
	fmt.Println("EXECUTING: handleWeatherElsewhere")
	msg := bot.NewMessage(um.update.Message.Chat.ID, commentsEn["DiffPlaceAccepted"])
	msg.ReplyMarkup = bot.ReplyKeyboardHide{HideKeyboard: true}
	_, err := um.bot.Send(msg)
	if err != nil {
		err = errors.Wrap(err, "Unable to Send message")
		fmt.Println(err)
	}
}

func handleLocationByText(um *UserMessage) {
	fmt.Println("EXECUTING: handleLocationByText")
	loc := retrieveCoordinates(um.connection, um.update.Message.Text)
	msg := bot.NewMessage(um.update.Message.Chat.ID, commentsEn["TryAgain"])
	msg.ReplyMarkup = keyboards["back"]
	if loc.Longitude != 0 && loc.Latitude != 0 { //this one makes it impossible to use bot from one place in Ghana
		um.update.Message.Location = &loc
		addLocationByCoords(um.connection, um.update.Message)
		msg = bot.NewMessage(um.update.Message.Chat.ID, commentsEn["CoordsAccepted"])
		msg.ReplyMarkup = keyboards["period"]
	}
	_, err := um.bot.Send(msg)
	if err != nil {
		err = errors.Wrap(err, "Unable to Send message")
		fmt.Println(err)
	}
}

func handleUnknown(um *UserMessage) {
	fmt.Println("EXECUTING: handleUnknown")
	msg := bot.NewMessage(um.update.Message.Chat.ID, commentsEn["Unknown"])
	msg.ReplyMarkup = keyboards["main"]
	_, err := um.bot.Send(msg)
	if err != nil {
		err = errors.Wrap(err, "Unable to Send message")
		fmt.Println(err)
	}
}
func handleEmpty(um *UserMessage) {
	fmt.Println("EXECUTING: handleEmpty")
	if um.update.Message.Location != nil {
		handleLocationByCoords(um)
	} else {
		handleUnknown(um)
	}
}

// Picks a random saying.
func pickASaying(sayings []string) string {
	return sayings[rand.Intn(len(sayings))]
}

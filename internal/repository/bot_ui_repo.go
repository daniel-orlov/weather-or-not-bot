package repository

import bot "gopkg.in/telegram-bot-api.v4"

var buttonsEN = map[string]string{
	"AtMyLocation": "Weather at my location",
	"AtADiffPlace": "Weather elsewhere",
	"Back0":        "< Back",
	"Back1":        "<< Back",
	"ByDays":       "By Days",
	"ByHours":      "By Hours",
	"Now":          "Now",
	"3Days":        "3 days",
	"5Days":        "5 days",
	"7Days":        "7 days",
	"10Days":       "10 days",
	"16Days":       "16 days",
	"24Hours":      "24 hours",
	"48Hours":      "48 hours",
	"72Hours":      "72 hours",
	"96Hours":      "96 hours",
	"120Hours":     "120 hours",
}

type BotUIRepo struct {
}

func NewBotUIRepo() *BotUIRepo {
	return &BotUIRepo{}
}

func (r *BotUIRepo) GetMainMenuKeyboard() bot.ReplyKeyboardMarkup {
	return bot.NewReplyKeyboard(
		bot.NewKeyboardButtonRow(bot.NewKeyboardButtonLocation(buttonsEN["AtMyLocation"])),
		bot.NewKeyboardButtonRow(bot.NewKeyboardButton(buttonsEN["AtADiffPlace"])),
	)
}

func (r *BotUIRepo) GetBackToMainMenuKeyboard() bot.ReplyKeyboardMarkup {
	return bot.NewReplyKeyboard(bot.NewKeyboardButtonRow(bot.NewKeyboardButton(buttonsEN["Back0"])))
}

func (r *BotUIRepo) GetDaysOrHoursKeyboard() bot.ReplyKeyboardMarkup {
	return bot.NewReplyKeyboard(
		bot.NewKeyboardButtonRow(bot.NewKeyboardButton(buttonsEN["ByHours"]), bot.NewKeyboardButton(buttonsEN["ByDays"])),
		bot.NewKeyboardButtonRow(bot.NewKeyboardButton(buttonsEN["Now"]), bot.NewKeyboardButton(buttonsEN["Back0"])),
	)
}

func (r *BotUIRepo) GetDaysKeyboard() bot.ReplyKeyboardMarkup {
	return bot.NewReplyKeyboard(
		bot.NewKeyboardButtonRow(
			bot.NewKeyboardButton(buttonsEN["3Days"]),
			bot.NewKeyboardButton(buttonsEN["5Days"]),
			bot.NewKeyboardButton(buttonsEN["7Days"]),
		),
		bot.NewKeyboardButtonRow(
			bot.NewKeyboardButton(buttonsEN["10Days"]),
			bot.NewKeyboardButton(buttonsEN["16Days"]),
			bot.NewKeyboardButton(buttonsEN["Back1"]),
		),
	)
}

func (r *BotUIRepo) GetHoursKeyboard() bot.ReplyKeyboardMarkup {
	return bot.NewReplyKeyboard(
		bot.NewKeyboardButtonRow(
			bot.NewKeyboardButton(buttonsEN["24Hours"]),
			bot.NewKeyboardButton(buttonsEN["48Hours"]),
			bot.NewKeyboardButton(buttonsEN["72Hours"]),
		),
		bot.NewKeyboardButtonRow(
			bot.NewKeyboardButton(buttonsEN["96Hours"]),
			bot.NewKeyboardButton(buttonsEN["120Hours"]),
			bot.NewKeyboardButton(buttonsEN["Back1"]),
		),
	)
}

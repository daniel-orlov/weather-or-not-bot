package repository

import bot "gopkg.in/telegram-bot-api.v4"

type BotDataRepo struct {
}

func (r *BotDataRepo) GetLocationKeyboard() bot.ReplyKeyboardMarkup {
	return bot.NewReplyKeyboard(
		bot.NewKeyboardButtonRow(bot.NewKeyboardButtonLocation(commentsEn["AtMyLocation"])),
		bot.NewKeyboardButtonRow(bot.NewKeyboardButton(commentsEn["AtADiffPlace"])),
	)
}

func (r *BotDataRepo) GetBackToMainMenuKeyboard() bot.ReplyKeyboardMarkup {
	return bot.NewReplyKeyboard(bot.NewKeyboardButtonRow(bot.NewKeyboardButton(commentsEn["Back0"])))
}

func (r *BotDataRepo) GetDaysOrHoursKeyboard() bot.ReplyKeyboardMarkup {
	return bot.NewReplyKeyboard(
		bot.NewKeyboardButtonRow(bot.NewKeyboardButton(commentsEn["ByHours"]), bot.NewKeyboardButton(commentsEn["ByDays"])),
		bot.NewKeyboardButtonRow(bot.NewKeyboardButton(commentsEn["Now"]), bot.NewKeyboardButton(commentsEn["Back0"])),
	)
}

func (r *BotDataRepo) GetDaysKeyboard() bot.ReplyKeyboardMarkup {
	return bot.NewReplyKeyboard(
		bot.NewKeyboardButtonRow(
			bot.NewKeyboardButton(commentsEn["3Days"]),
			bot.NewKeyboardButton(commentsEn["5Days"]),
			bot.NewKeyboardButton(commentsEn["7Days"]),
		),
		bot.NewKeyboardButtonRow(
			bot.NewKeyboardButton(commentsEn["10Days"]),
			bot.NewKeyboardButton(commentsEn["16Days"]),
			bot.NewKeyboardButton(commentsEn["Back1"]),
		),
	)
}

func (r *BotDataRepo) GetDaysKeyboard() bot.ReplyKeyboardMarkup {
	return bot.NewReplyKeyboard(
		bot.NewKeyboardButtonRow(
			bot.NewKeyboardButton(commentsEn["24Hours"]),
			bot.NewKeyboardButton(commentsEn["48Hours"]),
			bot.NewKeyboardButton(commentsEn["72Hours"]),
		),
		bot.NewKeyboardButtonRow(
			bot.NewKeyboardButton(commentsEn["96Hours"]),
			bot.NewKeyboardButton(commentsEn["120Hours"]),
			bot.NewKeyboardButton(commentsEn["Back1"]),
		),
	)
}

/*
var keyboards = map[string]bot.ReplyKeyboardMarkup{
	"main":   locationKeyboard,
	"period": daysOrHoursKeyboard,
	"days":   daysKeyboard,
	"hours":  hoursKeyboard,
	"back":   backToMainMenuKeyboard,
}
*/

//func fetchEmojis(backup map[string]int) map[string]int {
//	fmt.Println("EXECUTING: fetchEmojis")
//	cfg := parseConfig()
//	//establishing connection to database
//	conn, err := pgx.Connect(context.Background(), cfg.DbUrl)
//	if err != nil {
//		err = errors.Wrap(err, "Unable to connect to database")
//		fmt.Println(err)
//	}
//	defer conn.Close(context.Background())
//
//	var emojis = make(map[string]int)
//	sqlQuery := `SELECT name, code FROM emojis`
//	fmt.Println(sqlQuery)
//	rows, err := conn.Query(context.Background(), sqlQuery)
//	if err != nil {
//		err = errors.Wrap(err, "FAILED: Query when fetching Emojis")
//		fmt.Println(err)
//	}
//
//	var name string
//	var code int
//	for rows.Next() {
//		err = rows.Scan(&name, &code)
//		if err != nil {
//			err = errors.Wrap(err, "FAILED: Scanning a Row while fetching Emojis")
//			fmt.Println(err)
//		}
//		emojis[name] = code
//	}
//	err = rows.Err()
//	if err != nil {
//		err = errors.Wrap(err, "FAILED: Scan/Next a Row while fetching Emojis")
//		fmt.Println(err)
//	}
//	if len(emojis) == 0 {
//		return backup
//	}
//	return emojis
//}

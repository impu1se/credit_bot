package buttons

import tgbotapi "github.com/Syfaro/telegram-bot-api"

var FirstBtn = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Получить займ 💸"),
	),
)
var SecondBtn = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("От 100 грн.💰"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("От 200-400 грн.💰"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("От 500-1000 грн.💰"),
	),
)

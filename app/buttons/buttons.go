package buttons

import tgbotapi "github.com/Syfaro/telegram-bot-api"

var FirstBtn = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Получить займ 💸"),
	),
)
var SecondBtn = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("До 15.000р 💰"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("До 30.000р 💰"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("До 50.000р 💰"),
	),
)

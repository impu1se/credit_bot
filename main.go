package main

import (
	"fmt"
	"log"
	"net/http"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

const (
	apiToken = "843667644:AAEB7-te7PfsX2depO8nkeU3ZvNbEyDVpIk"
	credit30 = `
Ниже представлены компании подходящие под ваш запрос.

Чтобы получить займ до 30 тыс. руб. необходимо перейти по одной из ссылок ниже и заполнить анкету на сайте. (В течение 5 минут деньги придут вам на карту):

⬇️Компании в которых необходимо заполнить анкету:

Е-Капуста - первый займ до 30 000 руб. без переплаты
➡️ https://bit.ly/2NxFAAz

Робот Займер - автоматическое одобрение до 30 000 руб
➡️ https://bit.ly/2KZuVfW (нажмите на ссылку)

Е-заем - Акция: до 15 000 руб под 0% на 30 дней
➡️ https://bit.ly/2FTXs2e (нажмите на ссылку)

Metrokredit - первый заём под 0% на 15 дней
➡️ https://bit.ly/2LDWOJV

Совет: Чтобы увеличить вероятность и скорость одобрения займа, оставьте анкеты сразу в нескольких компаниях!✊
`
	url = "https://mysterious-woodland-23829.herokuapp.com/"
)

var firstBtn = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Получить займ 💸"),
	),
)
var secondBtn = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("До 30.000р 💰"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("До 50.000р 💰"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("До 100.000р 💰"),
	),
)

func main() {

	fmt.Println("Running bot...")
	bot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(url))
	if err != nil {
		log.Fatal(err)
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	updates := bot.ListenForWebhook("/" )
	go http.ListenAndServe(":80", nil)
	fmt.Println("Start serve")
	for update := range updates {

		if update.Message == nil {
			continue
		}
		chatID := update.Message.Chat.ID
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg := tgbotapi.NewMessage(chatID, "Hello!")
				msg.ReplyMarkup = firstBtn
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			}
		}
		if update.Message.Text != "" {
			switch update.Message.Text {
			case "Получить займ 💸":
				msg := tgbotapi.NewMessage(chatID, "Ок ща нахватаешь кредиков")
				msg.ReplyMarkup = secondBtn
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}

			case "До 30.000р 💰":
				msg := tgbotapi.NewMessage(chatID, credit30)
				msg.ReplyMarkup = firstBtn
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			case "До 50.000р 💰":
				msg := tgbotapi.NewMessage(chatID, "Получай 50 кесиков")
				msg.ReplyMarkup = firstBtn
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			case "До 100.000р 💰":
				msg := tgbotapi.NewMessage(chatID, "Получай 100 кесиков")
				msg.ReplyMarkup = firstBtn
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			}
		}
	}
}

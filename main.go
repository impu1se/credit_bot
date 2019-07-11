package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

const (
	url      = "https://credit-bot-for-yan.herokuapp.com/"
	apiToken = "843667644:AAEB7-te7PfsX2depO8nkeU3ZvNbEyDVpIk"
	credit30 = `
Ниже представлены компании подходящие под ваш запрос.

Чтобы получить займ до 30 тыс. руб. необходимо перейти по одной из ссылок ниже и заполнить анкету на сайте. (В течение 5 минут деньги придут вам на карту):

⬇️Компании в которых необходимо заполнить анкету:

Е-Капуста - первый займ до 30 000 руб. без переплаты
➡️ https://bit.ly/2YEuzyi (нажмите на ссылку)

Робот Займер - автоматическое одобрение до 30 000 руб
➡️ https://bit.ly/2JhC5dd (нажмите на ссылку)

Onzaem - первый займ от 2 000 до 30 000 рублей. Ставка 1.5%
➡️ https://bit.ly/32eNRg4 (нажмите на ссылку)

Деньги Сразу - первый займ до 30 000 руб. без переплаты
➡️ https://bit.ly/2XSvoGF (нажмите на ссылку)

Moneza - первый займ до 30 000 руб. без переплаты
➡️ https://bit.ly/2XtEHZp (нажмите на ссылку)

Совет: Чтобы увеличить вероятность и скорость одобрения займа, оставьте анкеты сразу в нескольких компаниях!
`
	credit15 = `
Чтобы получить займ до 15 тыс. руб. необходимо перейти по одной из ссылок ниже и заполнить анкету на сайте. (В течение 5 минут деньги придут вам на карту):

⬇️Компании в которых необходимо заполнить анкету:

Kredito24 - первый займ до 15 000 руб. Ставка 1,5%
➡️ https://bit.ly/2LKfnMJ (нажмите на ссылку)

E-zaem - первый займ до 15 000 руб. без переплат
➡️ https://bit.ly/2YAhdTD (нажмите на ссылку)

Metrokredit - первый займ до 15 000 рублей. без переплат
➡️ https://bit.ly/30nAkB7 (нажмите на ссылку)

SmartCredit - первый займ до 15 000 руб. 
Запустили программу лояльности - ставка 0.95%
➡️ https://bit.ly/2JfDW2a (нажмите на ссылку)

CreditPlus - первый займ до 15 000 руб. без переплаты
➡️ https://bit.ly/2Xx1atp (нажмите на ссылку)

Совет: Чтобы увеличить вероятность и скорость одобрения займа, оставьте анкеты сразу в нескольких компаниях!`

	credit50 = `
Чтобы получить займ до 50 тыс. руб. необходимо перейти по одной из ссылок ниже и заполнить анкету на сайте. (В течение 5 минут деньги придут вам на карту):

⬇️Компании в которых необходимо заполнить анкету:

Турбозайм - первый займ от 10 000 до 50 000. Ставка 0,8%
 ➡️ https://bit.ly/2S5aVt2 (нажмите на ссылку)

GreenMoney- первый займ до 40 000 руб. Ставка 0,95%
➡️ https://bit.ly/2RWrLKs (нажмите на ссылку)

Быстроденьги - первый займ до 50 000 рублей. 
➡️ https://bit.ly/2Xw5QQ8 (нажмите на ссылку)

Совет: Чтобы увеличить вероятность и скорость одобрения займа, оставьте анкеты сразу в нескольких компаниях!`
	welcome = `
Бот сотрудничает только с проверенными компаниями!

Перевод средств осуществляется любым удобным для вас способом.

P.S. Мы заинтересованы в том, чтобы вы получили займ! Внимательно следуйте инструкциям и вы получите перевод нужной вам суммы уже через 5 минут после оформления заявки!
Какая сумма вам нужна?

⬇️⬇️⬇️`

	afterStart = `
Здравствуйте, %v!
Проанализировав ваш профиль предлагаем займ у нашего партнера Е - капуста.
Для моментального, автоматического получения до 30.000 ₽ под 0 %% (сколько взяли столько и отдаете) до 30 дней оставьте заявку здесь: https://bit.ly/2YEuzyi (нажмите на ссылку)

💬 Или начните подбор других займов.`
	timerText = `
💳⁣ ЗАЙМ БЕЗ ПРОЦЕНТОВ 📢
Да-да, сколько взяли, столько отдали. Процентов - НЕТ❗️

📌 Шанс одобрения 98% 
📌 Нету процента по переплатам
📌 Самые проверенные предложения на рынке

Оформите заявку за 1 минуту прямо сейчас 👇

E-zaem - первый займ до 15 000 руб. без переплат
➡️ https://bit.ly/2YAhdTD (нажмите на ссылку)

Е-Капуста - первый займ до 30 000 руб. без переплаты
➡️ https://bit.ly/2YEuzyi (нажмите на ссылку)

CreditPlus - первый займ до 15 000 руб. без переплаты
➡️ https://bit.ly/2Xx1atp (нажмите на ссылку)
`
)

var firstBtn = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Получить займ 💸"),
	),
)
var secondBtn = tgbotapi.NewReplyKeyboard(
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

	var counter int64
	lastUpdate := make(map[int64]time.Time)
	ticker := time.NewTicker(1 * time.Hour)

	updates := bot.ListenForWebhook("/")
	go http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	fmt.Println("Start serve")
	for {
		select {
		case update := <-updates:
			if update.Message == nil {
				return
			}
			chatID := update.Message.Chat.ID
			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "start":
					counter++
					lastUpdate[chatID] = time.Now()
					msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(afterStart, update.Message.Chat.UserName))
					msg.ReplyMarkup = firstBtn
					if _, err := bot.Send(msg); err != nil {
						panic(err)
					}
				case "stat":
					msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%v", counter))
					if _, err := bot.Send(msg); err != nil {
						panic(err)
					}
				}
			}
			if update.Message.Text != "" {
				switch update.Message.Text {
				case "Получить займ 💸":
					msg := tgbotapi.NewMessage(chatID, welcome)
					msg.ReplyMarkup = secondBtn
					if _, err := bot.Send(msg); err != nil {
						panic(err)
					}

				case "До 15.000р 💰":
					msg := tgbotapi.NewMessage(chatID, credit15)
					msg.ReplyMarkup = firstBtn
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
					msg := tgbotapi.NewMessage(chatID, credit50)
					msg.ReplyMarkup = firstBtn
					if _, err := bot.Send(msg); err != nil {
						panic(err)
					}
				}
			}
		case _ := <-ticker.C:
			wakeUp(lastUpdate, bot)

		}
	}
}

func wakeUp(lastUpdate map[int64]time.Time, bot *tgbotapi.BotAPI) {
	timeNow := time.Now()
	for chatID, lastTime := range lastUpdate {
		diff := timeNow.Sub(lastTime)
		if diff > time.Duration(4*time.Hour) {
			msg := tgbotapi.NewMessage(chatID, timerText)
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
		}
		lastUpdate[chatID] = timeNow.Add(time.Duration(24) * time.Hour)

	}

}

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/go-redis/redis"
)

var textMsg = map[string]string{
	"credit30": `
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
`,
	"credit15": `
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

Совет: Чтобы увеличить вероятность и скорость одобрения займа, оставьте анкеты сразу в нескольких компаниях!`,

	"credit50": `
Чтобы получить займ до 50 тыс. руб. необходимо перейти по одной из ссылок ниже и заполнить анкету на сайте. (В течение 5 минут деньги придут вам на карту):

⬇️Компании в которых необходимо заполнить анкету:

Турбозайм - первый займ от 10 000 до 50 000. Ставка 0,8%
 ➡️ https://bit.ly/2S5aVt2 (нажмите на ссылку)

GreenMoney- первый займ до 40 000 руб. Ставка 0,95%
➡️ https://bit.ly/2RWrLKs (нажмите на ссылку)

Быстроденьги - первый займ до 50 000 рублей. 
➡️ https://bit.ly/2Xw5QQ8 (нажмите на ссылку)

Совет: Чтобы увеличить вероятность и скорость одобрения займа, оставьте анкеты сразу в нескольких компаниях!`,
	"welcome": `
Бот сотрудничает только с проверенными компаниями!

Перевод средств осуществляется любым удобным для вас способом.

P.S. Мы заинтересованы в том, чтобы вы получили займ! Внимательно следуйте инструкциям и вы получите перевод нужной вам суммы уже через 5 минут после оформления заявки!
Какая сумма вам нужна?

⬇️⬇️⬇️`,

	"afterStart": `
Здравствуйте, %v!
Проанализировав ваш профиль предлагаем займ у нашего партнера Е - капуста.
Для моментального, автоматического получения до 30.000 ₽ под 0 %% (сколько взяли столько и отдаете) до 30 дней оставьте заявку здесь: https://bit.ly/2YEuzyi (нажмите на ссылку)

💬 Или начните подбор других займов.`,
	"timerText": `
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
`,
}

const (
	layout = "2006-01-02T15:04:05"
)

func NewClient(addr string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       db, // use default DB
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}
	return client, nil
}

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

type Config struct {
	ApiToken string
	Port     string
	Address  string
	Debug    bool
	Tls      bool
}

func newConfig() *Config {

	var apiToken, port, addr string

	flag.StringVar(&apiToken, "token", "", "Telegram Bot Token")
	flag.StringVar(&port, "port", "80", "Port for server")
	flag.StringVar(&addr, "addr", "localhost", "Address for server")
	debug := flag.Bool("debug", false, "Debug true/false")
	tls := flag.Bool("tls", false, "TLS true/false")
	flag.Parse()

	if apiToken == "" {
		log.Print("-token is required")
		os.Exit(1)
	}

	return &Config{
		ApiToken: apiToken,
		Port:     port,
		Address:  addr,
		Debug:    *debug,
		Tls:      *tls,
	}
}

func main() {

	var mutex = &sync.Mutex{}

	client, err := NewClient("localhost:6379", 0)
	if err != nil {
		log.Fatal("not init redis client")
	}
	for typeText, text := range textMsg {
		err := setValue(client, typeText, text)
		if err != nil {
			log.Fatal(err)
		}
	}

	value, err := getValue(client, "counter")
	if value == "" {
		value = "0"
	}
	if err := setValue(client, "counter", value); err != nil {
		log.Fatal("can't set counter")
	}

	config := newConfig()

	fmt.Println("Running bot...")
	bot, err := tgbotapi.NewBotAPI(config.ApiToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = config.Debug
	log.Printf("Authorized on account %s", bot.Self.UserName)
	if config.Tls {
		_, err = bot.SetWebhook(tgbotapi.NewWebhookWithCert(config.Address+"/"+config.ApiToken, os.Getenv("CREDIT_CERT")))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		_, err = bot.SetWebhook(tgbotapi.NewWebhook(config.Address + "/" + config.ApiToken))
		if err != nil {
			log.Fatal(err)
		}
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	ticker := time.NewTicker(10 * time.Minute)

	updates := bot.ListenForWebhook("/" + bot.Token)
	if config.Tls {
		go http.ListenAndServeTLS(":"+config.Port, os.Getenv("CREDIT_CERT"), os.Getenv("CREDIT_KEY"), nil)
	} else {
		go http.ListenAndServe(":"+config.Port, nil)
	}

	fmt.Printf("Start server on %v:%v ", config.Address, config.Port)
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
					client.LPush("chatIds", chatID)
					mutex.Lock()
					counter, err := getValue(client, "counter")
					if err != nil {
						log.Fatal(err)
					}
					count, err := strconv.Atoi(counter)
					if err != nil {
						log.Fatal(err)
					}
					err = setValue(client, "counter", count+1)
					if err != nil {
						log.Fatal(err)
					}
					mutex.Unlock()
					updateTime := time.Now()

					if err := setValue(client, fmt.Sprintf("%v", chatID), updateTime.Format(layout)); err != nil {
						log.Fatal(err)
					}
					afterStart, _ := getValue(client, "afterStart")
					msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(afterStart, update.Message.Chat.UserName))
					msg.ReplyMarkup = firstBtn
					if _, err := bot.Send(msg); err != nil {
						log.Fatal(err)
					}
				case "stat":
					counter, _ := getValue(client, "counter")
					msg := tgbotapi.NewMessage(chatID, counter)
					if _, err := bot.Send(msg); err != nil {
						log.Fatal(err)
					}
				}
			}
			if update.Message.Text != "" {
				handlingText(update.Message.Text, chatID, bot, client)
			}
		case _ = <-ticker.C:
			wakeUp(bot, client)

		}
	}
}

func getValue(client *redis.Client, key string) (string, error) {
	val, err := client.Get(key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return val, nil
}

func setValue(client *redis.Client, key string, value interface{}) error {
	err := client.Set(key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func handlingText(text string, chatID int64, bot *tgbotapi.BotAPI, client *redis.Client) {
	switch text {
	case "Получить займ 💸":
		welcome, err := getValue(client, "welcome")
		if err != nil {
			log.Fatal("Not get value from welcome key")
		}
		msg := tgbotapi.NewMessage(chatID, welcome)
		msg.ReplyMarkup = secondBtn
		if _, err := bot.Send(msg); err != nil {
			log.Fatal(err)
		}

		err = updateTime(client, fmt.Sprintf("%v", chatID), getTime(client, chatID), 0)
		if err != nil {
			log.Fatal(err)
		}

	case "До 15.000р 💰":
		credit15, err := getValue(client, "credit15")
		if err != nil {
			log.Fatal("Not get value from credit15 key")
		}
		msg := tgbotapi.NewMessage(chatID, credit15)
		msg.ReplyMarkup = firstBtn
		if _, err := bot.Send(msg); err != nil {
			log.Fatal(err)
		}
		err = updateTime(client, fmt.Sprintf("%v", chatID), getTime(client, chatID), 0)
		if err != nil {
			log.Fatal(err)
		}

	case "До 30.000р 💰":
		credit30, err := getValue(client, "credit30")
		if err != nil {
			log.Fatal("Not get value from credit30 key")
		}
		msg := tgbotapi.NewMessage(chatID, credit30)
		msg.ReplyMarkup = firstBtn
		if _, err := bot.Send(msg); err != nil {
			log.Fatal(err)
		}
		err = updateTime(client, fmt.Sprintf("%v", chatID), getTime(client, chatID), 0)
		if err != nil {
			log.Fatal(err)
		}

	case "До 50.000р 💰":
		credit50, err := getValue(client, "credit50")
		if err != nil {
			log.Fatal("Not get value from credit50 key")
		}
		msg := tgbotapi.NewMessage(chatID, credit50)
		msg.ReplyMarkup = firstBtn
		if _, err := bot.Send(msg); err != nil {
			log.Fatal(err)
		}
		err = updateTime(client, fmt.Sprintf("%v", chatID), getTime(client, chatID), 0)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getTime(client *redis.Client, chatID int64) time.Time {
	strChatID := strconv.Itoa(int(chatID))
	lastActive, err := getValue(client, strChatID)
	if err != nil {
		log.Fatal(err)
	}
	t, err := time.Parse(layout, lastActive)
	if err != nil {
		log.Printf("can't parse time")
		t = time.Now()
	}
	return t
}

func wakeUp(bot *tgbotapi.BotAPI, client *redis.Client) {
	chatIds, err := client.LRange("chatIds", 0, -1).Result()
	if err != nil {
		log.Fatal(err)
	}
	for _, chatId := range chatIds {
		lastTime, err := getValue(client, chatId)
		if err != nil {
			log.Fatal(err)
		}
		timeNow := time.Now()
		t, err := time.Parse(layout, lastTime)
		if err != nil {
			log.Fatal(err)
		}

		diff := timeNow.Sub(t)
		if diff > time.Duration(4*time.Hour) {
			timerText, err := getValue(client, "timerText")
			if err != nil {
				log.Fatalf("not get value from timer text key with err %v", err)
			}
			intChatId, err := strconv.Atoi(chatId)
			if err != nil {
				log.Fatal(err)
			}
			msg := tgbotapi.NewMessage(int64(intChatId), timerText)
			if _, err := bot.Send(msg); err != nil {
				client.LRem("chatIds", 0, chatId)
			}
			if err := updateTime(client, chatId, timeNow, 24); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func updateTime(client *redis.Client, chatId string, lastTime time.Time, hour int) error {
	var newTime = time.Now().Format(layout)
	if hour != 0 {
		newTime = lastTime.Add(time.Duration(hour) * time.Hour).Format(layout)
	}
	err := setValue(client, chatId, newTime)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("update chatid %v from %v time to new %v time", chatId, lastTime, newTime)
	return nil
}

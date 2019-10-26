package app

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/impu1se/credit_bot/app/messages"

	"github.com/impu1se/credit_bot/app/buttons"

	tgbotapi "github.com/Syfaro/telegram-bot-api"

	"github.com/impu1se/credit_bot/app/db"

	"github.com/impu1se/credit_bot/app/config"
)

const (
	layout = "2006-01-02T15:04:05"
)

type CreditBot struct {
	Config   config.Config
	Redis    *db.MyRedis
	Updates  tgbotapi.UpdatesChannel
	Postgres interface{}
	Mutex    sync.Mutex
	Ticker   *time.Ticker
}

var valueFromRedis = map[string]string{
	"Получить займ 💸": "welcome",
	"До 15.000р 💰":    "credit15",
	"До 30.000р 💰":    "credit30",
	"До 50.000р 💰":    "credit50",
}

func NewCreditBot(conf config.Config, client *db.MyRedis, postgres interface{}, update tgbotapi.UpdatesChannel) *CreditBot {
	return &CreditBot{
		Config:   conf,
		Redis:    client,
		Postgres: postgres,
		Updates:  update,
		Ticker:   time.NewTicker(10 * time.Minute)}
}

func (c *CreditBot) init() {
	value, err := c.Redis.GetValue("counter")
	if err != nil {
		log.Println("can't get value counter with err:", err)
	}
	if value == "" {
		value = "0"
	}
	if err := c.Redis.SetValue("counter", value); err != nil {
		log.Println("can't set counter")
	}
}

func (c *CreditBot) Run(bot *tgbotapi.BotAPI) {
	for {
		select {
		case update := <-c.Updates:
			if update.Message == nil {
				return
			}
			if update.Message.IsCommand() {
				c.commandHandle(bot, &update)
				continue
			}
			if update.Message.Text != "" {
				c.handlingText(bot, &update)
				continue
			}
		case _ = <-c.Ticker.C:
			c.wakeUp(bot)
		}
	}
}

// Handle for command
func (c *CreditBot) commandHandle(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	switch update.Message.Command() {
	case "start":
		c.Redis.Client.LPush("chatIds", chatID)
		c.Mutex.Lock()

		counter, err := c.Redis.GetValue("counter")
		if err != nil {
			log.Print(err)
		}
		count, err := strconv.Atoi(counter)
		if err != nil {
			log.Print(err)
		}
		err = c.Redis.SetValue("counter", count+1)
		if err != nil {
			log.Print(err)
		}
		c.Mutex.Unlock()
		updateTime := time.Now()

		if err := c.Redis.SetValue(fmt.Sprintf("%v", chatID), updateTime.Format(layout)); err != nil {
			log.Print(err)
			return
		}
		afterStart, _ := c.Redis.GetValue("afterStart")
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(afterStart, update.Message.Chat.UserName))
		msg.ReplyMarkup = buttons.FirstBtn
		if _, err := bot.Send(msg); err != nil {
			log.Print(err)
			return
		}
	case "stat":
		if update.Message.Chat.UserName == "betferma" {
			counter, _ := c.Redis.GetValue("counter")
			msg := tgbotapi.NewMessage(chatID, counter)
			if _, err := bot.Send(msg); err != nil {
				log.Print(err)
			}
		}
		return
	}
}

// Handle for Text message
func (c CreditBot) handlingText(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	switch update.Message.Text {
	case "Получить займ 💸":
		if err := c.handleText(bot, valueFromRedis[update.Message.Text], chatID, buttons.SecondBtn); err != nil {
			fmt.Print(err)
		}
	case "До 15.000р 💰":
		if err := c.handleText(bot, valueFromRedis[update.Message.Text], chatID, buttons.FirstBtn); err != nil {
			fmt.Print(err)
		}
	case "До 30.000р 💰":
		if err := c.handleText(bot, valueFromRedis[update.Message.Text], chatID, buttons.FirstBtn); err != nil {
			fmt.Print(err)
		}
	case "До 50.000р 💰":
		if err := c.handleText(bot, valueFromRedis[update.Message.Text], chatID, buttons.FirstBtn); err != nil {
			fmt.Print(err)
		}
	default:
		c.updateText(bot, chatID, update.Message.Text)
	}
}

func (c *CreditBot) updateText(bot *tgbotapi.BotAPI, chatID int64, text string) {
	res := strings.Split(text, ";")
	if _, ok := messages.TextMsg[res[0]]; !ok {
		log.Println("can't get new text message for update, res :", res)
		return
	}
	if len(res) > 1 {
		key := res[0]
		message := res[1]
		if err := c.Redis.SetValue(key, message); err != nil {
			fmt.Printf("can't update new message for %v , because %v :\n", key, err)
		}
	}
}

func (c *CreditBot) wakeUp(bot *tgbotapi.BotAPI) {
	chatIds, err := c.Redis.Client.LRange("chatIds", 0, -1).Result()
	if err != nil {
		log.Println(err)
	}
	for _, chatId := range chatIds {
		lastTime, err := c.Redis.GetValue(chatId)
		if err != nil {
			log.Println(err)
		}
		timeNow := time.Now()
		t, err := time.Parse(layout, lastTime)
		if err != nil {
			log.Println(err)
		}
		diff := timeNow.Sub(t)
		if diff > 4*time.Hour {
			timerText, err := c.Redis.GetValue("timerText")
			if err != nil {
				log.Println("not get value from timer text key with err:", err)
			}
			intChatId, err := strconv.Atoi(chatId)
			if err != nil {
				log.Print(err)
			}
			msg := tgbotapi.NewMessage(int64(intChatId), timerText)
			if _, err := bot.Send(msg); err != nil {
				c.Redis.Client.LRem("chatIds", 0, chatId)
				counter, err := c.Redis.GetValue("counter")
				if err == nil {
					count, err := strconv.Atoi(counter)
					if err == nil {
						if err := c.Redis.SetValue("counter", count-1); err != nil {
							fmt.Println()
						}
					}
				}
				return
			}
			if err := c.updateTime(chatId, timeNow, 24); err != nil {
				log.Print(err)
			}
		}
	}
}

func (c CreditBot) getTime(chatID int64) time.Time {
	strChatID := strconv.Itoa(int(chatID))
	lastActive, err := c.Redis.GetValue(strChatID)
	if err != nil {
		log.Print(err)
	}
	t, err := time.Parse(layout, lastActive)
	if err != nil {
		log.Println("can't parse time")
		t = time.Now()
	}
	return t
}

func (c CreditBot) updateTime(chatId string, lastTime time.Time, hour int) error {
	var newTime = time.Now().Format(layout)
	if hour != 0 {
		newTime = lastTime.Add(time.Duration(hour) * time.Hour).Format(layout)
	}
	err := c.Redis.SetValue(chatId, newTime)
	if err != nil {
		return err
	}
	fmt.Printf("update chatid %v from %v time to new %v time\n", chatId, lastTime, newTime)
	return nil
}

func (c CreditBot) handleText(bot *tgbotapi.BotAPI, value string, chatID int64, button tgbotapi.ReplyKeyboardMarkup) error {
	text, err := c.Redis.GetValue(value)
	if err != nil {
		log.Printf("Not get value from %v key\n", value)
	}
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = button
	if _, err := bot.Send(msg); err != nil {
		return err
	}
	err = c.updateTime(fmt.Sprintf("%v", chatID), c.getTime(chatID), 0)
	if err != nil {
		log.Print(err)
	}
	return nil
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/impu1se/credit_bot/app"

	"github.com/impu1se/credit_bot/app/config"
	"github.com/impu1se/credit_bot/app/db"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func main() {

	client := db.NewClient("localhost:6379", 0)

	conf := config.NewConfig()

	fmt.Println("Running bot...")
	bot, err := tgbotapi.NewBotAPI(conf.ApiToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = conf.Debug
	log.Printf("Authorized on account %s", bot.Self.UserName)

	if conf.Tls {
		_, err = bot.SetWebhook(tgbotapi.NewWebhookWithCert(conf.Address+"/"+conf.ApiToken, os.Getenv("CREDIT_CERT")))
		if err != nil {
			log.Print(err)
		}
		go http.ListenAndServeTLS(":"+conf.Port, os.Getenv("CREDIT_CERT"), os.Getenv("CREDIT_KEY"), nil)

	} else {
		_, err = bot.SetWebhook(tgbotapi.NewWebhook(conf.Address + "/" + conf.ApiToken))
		if err != nil {
			log.Print(err)
		}
		go http.ListenAndServe(":"+conf.Port, nil)

	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Print(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)

	creditBot := app.NewCreditBot(conf, client, nil, updates)

	fmt.Printf("Start server on %v:%v ", conf.Address, conf.Port)

	creditBot.Run(bot)
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/impu1se/credit_bot/app/metrics"

	"github.com/impu1se/credit_bot/app"

	"github.com/impu1se/credit_bot/app/config"
	"github.com/impu1se/credit_bot/app/db"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func main() {

	conf := config.NewConfig()

	clientRedis := db.NewClientRedis(*conf)

	fmt.Println("Running bot...")
	bot, err := tgbotapi.NewBotAPI(conf.ApiToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = conf.Debug
	log.Printf("Authorized on account %s\n", bot.Self.UserName)

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
		log.Printf("Telegram callback failed: %s\n", info.LastErrorMessage)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)

	metric := metrics.New()

	creditBot := app.NewCreditBot(*conf, clientRedis, nil, updates, metric)

	prometheus.MustRegister(creditBot.Metrics.Collectors()...)
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":9010", nil)

	fmt.Printf("Start server on %v:%v \n", conf.Address, conf.Port)

	err = creditBot.InitCreditBot()
	if err != nil {
		log.Panicf("can't init credit bot with err: %v\n", err)
	}

	creditBot.Run(bot)
}

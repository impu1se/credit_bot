package db

import (
	"log"

	"github.com/impu1se/credit_bot/app/config"

	"github.com/impu1se/credit_bot/app/messages"

	"github.com/go-redis/redis"
)

type MyRedis struct {
	Client *redis.Client
}

func NewClientRedis(conf config.Config) *MyRedis {
	if conf.RedisPort != "" {
		conf.RedisHost = conf.RedisHost + ":" + conf.RedisPort
	}
	client := redis.NewClient(&redis.Options{
		Addr:     conf.RedisHost,
		Password: "",           // no password set
		DB:       conf.RedisDb, // use default DB
	})

	log.Println("trying connect to:", client.String())
	if _, err := client.Ping().Result(); err != nil {
		log.Fatal(err)
	}
	log.Println("connected successful...")
	cli := &MyRedis{client}
	cli.initText(messages.TextMsg)
	return cli
}

func (r MyRedis) GetValue(key string) (string, error) {
	val, err := r.Client.Get(key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func (r MyRedis) SetValue(key string, value interface{}) error {
	err := r.Client.Set(key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r MyRedis) initText(textMsg map[string]string) {
	for typeText, text := range textMsg {
		res, _ := r.GetValue(typeText)
		if res == "" {
			err := r.SetValue(typeText, text)
			if err != nil {
				log.Print(err)
			}
		}
	}
}

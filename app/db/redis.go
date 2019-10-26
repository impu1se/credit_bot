package db

import (
	"log"

	"github.com/impu1se/credit_bot/app/messages"

	"github.com/go-redis/redis"
)

type MyRedis struct {
	Client *redis.Client
}

func NewClient(addr string, db int) *MyRedis {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       db, // use default DB
	})

	if _, err := client.Ping().Result(); err != nil {
		log.Print(err)
	}
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

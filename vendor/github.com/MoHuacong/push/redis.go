package push

import (
	"github.com/MoHuacong/wic"
	"github.com/go-redis/redis"
)

type Redis struct {
	mq *wic.Mq
	client *redis.Client
}

func (redis *Redis) IsConnc() error {
	_, err := redis.client.Ping().Result()
	if err != nil {
		redis.mq.Topic("log-error").Send(err.Error())
		return err
	}
	return nil
}

package configs

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

func (r Redis) newClient() *redis.Client {
	host := fmt.Sprintf("%s:%s", r.Host, r.Port)
	password := r.Password

	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       0,
	})

	rdb = client

	return rdb
}

func (r Redis) GetRedisClient() *redis.Client {
	return r.newClient()
}

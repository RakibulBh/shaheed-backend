package redis

import (
	"github.com/redis/go-redis/v9"
)

func New(addr string, password string, db int, protocol int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		Protocol: protocol,
	})

	return client, nil
}

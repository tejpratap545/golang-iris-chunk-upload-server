package service

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

func (client *RedisClient) SetClient() {
	client.Client = NewRedisClient()
}
func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func (client *RedisClient) Set(key string, value interface{}) (*redis.StatusCmd, error) {
	p, err := json.Marshal(value)
	return client.Client.Set(client.Ctx, key, p, 0), err
}

func (client *RedisClient) Get(key string, dest interface{}) error {
	p, _ := client.Client.Get(client.Ctx, key).Bytes()

	return json.Unmarshal(p, dest)
}

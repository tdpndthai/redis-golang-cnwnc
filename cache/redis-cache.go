package cache

import (
	"encoding/json"
	"redis-golang-cnwnc/entities"
	"time"

	"github.com/go-redis/redis/v7"
)

type redisCache struct {
	host    string
	db      int
	exipres time.Duration
}

func NewRedisCache(host string, db int, exp time.Duration) ProductCache {
	return &redisCache{
		host:    host,
		db:      db,
		exipres: exp,
	}
}

func (cache *redisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: "",
		DB:       cache.db,
	})
}

func (cache *redisCache) Set(key string, product *[]entities.Product) {
	client := cache.getClient()

	json, err := json.Marshal(product)
	if err != nil {
		panic(err)
	}

	client.Set(key, json, cache.exipres*time.Second)
}

func (cache *redisCache) Get(key string) *[]entities.Product {
	client := cache.getClient()
	value, err := client.Get(key).Result()
	if err != nil {
		return nil
	}

	product := []entities.Product{}

	err = json.Unmarshal([]byte(value), &product)

	if err != nil {
		panic(err)
	}

	return &product
}

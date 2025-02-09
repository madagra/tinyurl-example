package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var longToShortH string = "long:to:short"
var shortToLongH string = "short:to:long"

// no key expiration since Redis is used as a database here
const CacheDuration = 0

// Redis distributed cache client
type RedisDbClient struct {
	client  *redis.Client
	context context.Context
}

func GetRedisDbClient() *RedisDbClient {

	var addr string = GetDbAddress()
	log.Debug().Msgf("Address Redis: %s", addr)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.Background()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Error initializing Redis: %v", err))
	}

	dbClient := new(RedisDbClient)
	dbClient.client = client
	dbClient.context = ctx

	return dbClient

}

func (client RedisDbClient) ExistLongUrl(url string) bool {

	exists, err := client.client.HExists(client.context, longToShortH, url).Result()
	log.Debug().Msgf("Result from Redis HEXIST with key %s: %s, %t, %v", longToShortH, url, exists, err)
	if err != nil {
		log.Error().Msgf("Error checking long URL existence: %s", err.Error())
		return false
	}

	return exists
}

func (client RedisDbClient) ExistShortUrl(url string) bool {

	exists, err := client.client.HExists(client.context, shortToLongH, url).Result()
	log.Debug().Msgf("Result from Redis HEXIST with key %s: %s, %t, %v", shortToLongH, url, exists, err)
	if err != nil {
		log.Error().Msgf("Error checking short URL existence: %s", err.Error())
		return false
	}

	return exists
}

func (client RedisDbClient) StoreShortUrl(shortUrl string, longUrl string) {
	// set both direct and inverse mapping with two different hash sets

	directVal := map[string]string{
		longUrl: shortUrl,
	}

	inverseVal := map[string]string{
		shortUrl: longUrl,
	}

	if _, err := client.client.HSet(client.context, longToShortH, directVal).Result(); err != nil {
		log.Error().Msgf("Error storing long to short URL mapping: %v", err)
		return
	}

	if _, err := client.client.HSet(client.context, shortToLongH, inverseVal).Result(); err != nil {
		log.Error().Msgf("Error storing long to short URL mapping: %v", err)
		return
	}

}

func (client RedisDbClient) RetrieveLongUrl(shortUrl string) string {

	res := client.client.HGet(client.context, shortToLongH, shortUrl)

	if res.Err() != nil {
		log.Info().Msgf("Cannot retrieve URL: %s", shortUrl)
	}

	return res.Val()

}

func (client RedisDbClient) RetrieveShortUrl(longUrl string) string {

	res := client.client.HGet(client.context, longToShortH, longUrl)

	if res.Err() != nil {
		log.Info().Msgf("Cannot retrieve URL: %s", longUrl)
	}

	return res.Val()
}

func (client RedisDbClient) HasKey(key string) bool {
	exists, err := client.client.Exists(client.context, key).Result()

	if err != nil {
		log.Error().Msgf("Error checking key %s", key)
		return false
	}

	return exists == 1
}

func (client RedisDbClient) SetKey(newKey string) string {

	res, err := client.client.Set(client.context, newKey, true, CacheDuration).Result()

	if err != nil {
		log.Error().Msgf("Error setting key %s", newKey)
	}

	return res
}

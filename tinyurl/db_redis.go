package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var longToShortH string = "long:to:short"
var shortToLongH string = "short:to:long"
var keyH = "url:key"


// no key expiration since Redis is used as a database here
const CacheDuration = 0

// Redis distributed cache client
type RedisDbClient struct {
	client  *redis.Client
	context context.Context
}

func GetRedisDbClient(isLocal bool) *RedisDbClient {

	var addr string
	if isLocal {
		addr = LocalDbAddr
	} else {
		addr = RemoteDbAddr
	}

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
	if err != nil {
		log.Error().Msgf("Error checking long URL existence: %s", err.Error())
		return false
	}

	return exists
}

func (client RedisDbClient) ExistShortUrl(url string) bool {

	exists, err := client.client.HExists(client.context, shortToLongH, url).Result()
	if err != nil {
		log.Error().Msgf("Error checking short URL existence: %s", err.Error())
		return false
	}

	return exists
}

func (client RedisDbClient) StoreShortUrl(shortUrl string, longUrl string) {

	// set both direct and inverse mapping
	// with two different hash sets

	directVal := map[string]interface{}{
		longUrl: shortUrl,
	}

	inverseVal := map[string]interface{}{
		shortUrl: longUrl,
	}

	err1 := client.client.HSet(client.context, longToShortH, directVal)
	err2 := client.client.HSet(client.context, shortToLongH, inverseVal)

	if err1 != nil || err2 != nil {
		log.Error().Msgf("Error storing short and/or long URLs: %s, %s", err1.Err().Error(), err2.Err().Error())
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

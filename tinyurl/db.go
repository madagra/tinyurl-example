package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/vishalkuo/bimap"
)

var LocalDbAddr = "localhost:6379"

// TODO
var RemoteDbAddr = "localhost:6379"

var HUrlToKey string = "url:to:key"
var HKeyToUrl string = "key:to:url"

type DbInterface interface {
	ExistLongUrl(url string) bool
	ExistShortUrl(url string) bool
	RetrieveShortUrl(longUrl string) string
	RetrieveLongUrl(shortUrl string) string
	StoreShortUrl(shortUrl string, longUrl string)
	HasKey(key string) bool
	SetKey(newKey string) string
}

// Redis distributed cache client
type RedisDbClient struct {
	client  *redis.Client
	context context.Context
}

// in-memory data structure
type LocalDbClient struct {
	urlDB     *bimap.BiMap[string, string]
	urlKeysDB map[string]bool
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

func (client LocalDbClient) ExistShortUrl(url string) bool {
	return client.urlDB.Exists(url)
}

func (client LocalDbClient) ExistLongUrl(url string) bool {
	return client.urlDB.ExistsInverse(url)
}

func (client LocalDbClient) StoreShortUrl(shortUrl string, longUrl string) {
	client.urlDB.Insert(shortUrl, longUrl)
}

func (client LocalDbClient) RetrieveLongUrl(shortUrl string) string {
	var longUrl, ok = client.urlDB.Get(shortUrl)
	if !ok {
		log.Error().Msgf("Short URL: %s not found in the DB", shortUrl)
	}
	return longUrl
}

func (client LocalDbClient) RetrieveShortUrl(longUrl string) string {
	var shortUrl, ok = client.urlDB.GetInverse(longUrl)
	if !ok {
		log.Error().Msgf("Long URL: %s not found in the DB", longUrl)
	}
	return shortUrl
}

func (client LocalDbClient) HasKey(key string) bool {
	var _, exist = client.urlKeysDB[key]
	return exist
}

func (client LocalDbClient) SetKey(newKey string) string {
	client.urlKeysDB[newKey] = true
	return newKey
}

func (client RedisDbClient) ExistShortUrl() {

}

func (client RedisDbClient) ExistLongUrl() {

}

func (client RedisDbClient) StoreShortUrl() {

}

func (client RedisDbClient) RetrieveLongUrl() {

}

func GetLocalDbClient() *LocalDbClient {

	var urlDB = bimap.NewBiMap[string, string]()
	var urlKeysDB map[string]bool = make(map[string]bool)

	dbClient := new(LocalDbClient)
	dbClient.urlDB = urlDB
	dbClient.urlKeysDB = urlKeysDB

	return dbClient
}

func PurgeLocalDb(dbClient *LocalDbClient) {
	for key := range dbClient.urlKeysDB {
		delete(dbClient.urlKeysDB, key)
	}
}

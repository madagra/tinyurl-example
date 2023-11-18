package main

import (
	"github.com/rs/zerolog/log"
	"github.com/vishalkuo/bimap"
)

// in-memory data structure
type LocalDbClient struct {
	urlDB     *bimap.BiMap[string, string]
	urlKeysDB map[string]bool
}

func (client LocalDbClient) ExistLongUrl(url string) bool {
	return client.urlDB.ExistsInverse(url)
}

func (client LocalDbClient) ExistShortUrl(url string) bool {
	return client.urlDB.Exists(url)
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

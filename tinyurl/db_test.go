package main

import (
	"testing"
)

var testLongUrl string = "https://learn.cantrill.io/courses/enrolled/1820301"
var testShortUrl string = "https://mdtiny.net/abcdef"


func initLocalTestDb(t *testing.T) LocalDbClient {
	var dbClient = GetLocalDbClient()

	// this function is registered to be called at the
	// of every test which uses this function via the
	// Cleanup() interface
	t.Cleanup(func() {
		for key := range dbClient.urlKeysDB {
			delete(dbClient.urlKeysDB, key)
		}
	})

	return *dbClient
}

func initRedisTestDb(t *testing.T) RedisDbClient {
	var dbClient = GetRedisDbClient(true)

	t.Cleanup(func() {
		dbClient.client.FlushAll(dbClient.context)
	})

	return *dbClient
}

func TestRedisStoreAndExistsKey(t *testing.T) {

	const testKey = "testKey"
	const wrongKey = "wrongKey"

	client := initRedisTestDb(t)
	
	client.SetKey(testKey)
	exists1 := client.HasKey(testKey)
	exists2 := client.HasKey(wrongKey)

	if !exists1 {
		t.Errorf("Key %s not found in the database", testKey)
	}
	
	if exists2 {
		t.Errorf("Wrong key %s found in the database", wrongKey)
	}
	
}

// TODO
func TestRedisStoreAndExistsUrls(t *testing.T) {
	
}
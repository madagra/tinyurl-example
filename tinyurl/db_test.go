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

func TestRedisStoreAndExistsUrls(t *testing.T) {

	const wrongUrl = "wrongUrl"

	client := initRedisTestDb(t)

	client.StoreShortUrl(testShortUrl, testLongUrl)
	existLong := client.ExistLongUrl(testLongUrl)
	existShort := client.ExistShortUrl(testShortUrl)

	if !existLong || !existShort {
		t.Errorf("Failed storing URLs %s, %s in the DB", testLongUrl, testShortUrl)
	}

	existWrong1 := client.ExistLongUrl(wrongUrl)
	existWrong2 := client.ExistShortUrl(wrongUrl)

	if existWrong1 || existWrong2 {
		t.Errorf("Wrong URL found in DB: %s", wrongUrl)
	}

}

func TestRedisRetrieveUrls(t *testing.T) {
	client := initRedisTestDb(t)

	client.StoreShortUrl(testShortUrl, testLongUrl)
	
	var long string = client.RetrieveLongUrl(testShortUrl) 
	var short string = client.RetrieveShortUrl(testLongUrl)

	if short != testShortUrl {
		t.Errorf("Error retrieving short URL %s, found %s", testShortUrl, short)
	}

	if long != testLongUrl {
		t.Errorf("Error retrieving long URL %s, found %s", testLongUrl, long)
	}
}
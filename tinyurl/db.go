package main

import "os"

var DbAddrEnv = "DB_ADDRESS"
var LocalDbAddr = "localhost:6379"

type DBManager interface {
	ExistLongUrl(url string) bool
	ExistShortUrl(url string) bool
	RetrieveShortUrl(longUrl string) string
	RetrieveLongUrl(shortUrl string) string
	StoreShortUrl(shortUrl string, longUrl string)
	HasKey(key string) bool
	SetKey(newKey string) string
}

func GetDbAddress() string {

	var dbAddress string

	value, exists := os.LookupEnv(DbAddrEnv)
	if exists {
		dbAddress = value
	} else {
		dbAddress = LocalDbAddr
	}
	return dbAddress
}

package main

var LocalDbAddr = "localhost:6379"

// TODO
var RemoteDbAddr = "localhost:6379"

type DbInterface interface {
	ExistLongUrl(url string) bool
	ExistShortUrl(url string) bool
	RetrieveShortUrl(longUrl string) string
	RetrieveLongUrl(shortUrl string) string
	StoreShortUrl(shortUrl string, longUrl string)
	HasKey(key string) bool
	SetKey(newKey string) string
}

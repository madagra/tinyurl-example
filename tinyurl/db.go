package main

import "github.com/vishalkuo/bimap"

var UrlDB = bimap.NewBiMap[string, string]()

var UrlKeysDB map[string]bool = make(map[string]bool)
var ShortUrlForwardDB map[string]string = make(map[string]string)
var ShortUrlInverseDB map[string]string = make(map[string]string)

func PurgeUrlDB() {
	for key := range UrlKeysDB {
		delete(UrlKeysDB, key)
	}
}

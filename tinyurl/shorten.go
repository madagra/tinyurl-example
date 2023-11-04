package main

import (
	b64 "encoding/base64"
	"log"
	"math/rand"
	"strings"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var LenShortUrl = 10

func keyGenerator(n int) string {

	// TODO: remove deprecated function
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// shorten the URL by generating a new key from a random set
// of characters. If a collision is found, the key is regenerated
func ShortenUrlKeygen(url string, prefix string, db DbInterface) (string, string) {

	var newKey string
	var not_found bool = true

	strippedUrl := strings.ReplaceAll(url, "https://", "")
	strippedUrl = strings.ReplaceAll(strippedUrl, "http://", "")

	for not_found {

		newKey = keyGenerator(LenShortUrl)
		var exists = db.HasKey(newKey)

		if !exists {
			db.SetKey(newKey)
			break
		}
		log.Println("Found a collision. Regenerating the key...")
	}

	var newUrl string = prefix + newKey
	return newUrl, newKey
}

// deprecated
// convert a url into its shortened version using
// base64 encoding with a rotating chunk of the encoding with the
// desired length
func ShortenUrlEncoding(url string, prefix string, db DbInterface) (string, string) {

	strippedUrl := strings.ReplaceAll(url, "https://", "")
	strippedUrl = strings.ReplaceAll(strippedUrl, "http://", "")
	var encoded string = b64.StdEncoding.EncodeToString([]byte(strippedUrl))

	var shortUrl string
	for i := 0; i < len(encoded)-LenShortUrl; i++ {

		shortUrl = encoded[i : i+LenShortUrl]
		var exists = db.HasKey(shortUrl)
		if !exists {
			db.SetKey(shortUrl)
			break
		} else {
			shortUrl = encoded
		}
	}

	if shortUrl == encoded {
		log.Println("Returning long URL, found collisions")
	}

	var newUrl string = prefix + shortUrl
	return newUrl, encoded
}

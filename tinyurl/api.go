package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type UrlToShorten struct {
	Url string `json:"url" xml:"url" form:"url"`
}

// TODO: replace all logs with either Logrus or Zerolog to use
// a level logging library

func CreateServer(urlPrefix string) *fiber.App {

	app := fiber.New()

	// health check endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("TinyURL is healthy")
	})

	// save a new URL and its corresponding shortened version
	app.Post("/shorten", func(c *fiber.Ctx) error {

		body := new(UrlToShorten)
		if err := c.BodyParser(body); err != nil {
			return err
		}

		if !UrlDB.ExistsInverse(body.Url) {
			var shortenedUrl, _ = ShortenUrlKeygen(body.Url, urlPrefix)
			log.Print("URL not found")
			UrlDB.Insert(shortenedUrl, body.Url)
		} else {
			url, _ := UrlDB.GetInverse(body.Url)
			log.Printf("Already in DB: %s", url)
		}

		return c.SendStatus(200)
	})

	app.Use("/", func(c *fiber.Ctx) error {

		var shortUrl string = c.Context().URI().String()
		// var urlKey string = strings.Split(strings.ReplaceAll(shortUrl, urlPrefix, ""), ".")[0]

		if UrlDB.Exists(shortUrl) {
			longUrl, _ := UrlDB.Get(shortUrl)
			log.Printf("Long URL: %s", longUrl)
			return c.Redirect(longUrl)
		} else {
			c.SendString("The requested URL does not exist!")
			return c.SendStatus(404)
		}
	})

	return app
}

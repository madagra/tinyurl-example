package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type UrlToShorten struct {
	Url string `json:"url" xml:"url" form:"url"`
}

func CreateServer(urlPrefix string, db DbInterface) *fiber.App {

	app := fiber.New()

	// health check endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("TinyURL is healthy")
	})

	// save a new URL and its corresponding shortened version
	app.Post("/shorten", func(c *fiber.Ctx) error {

		body := new(UrlToShorten)
		if err := c.BodyParser(body); err != nil {
			return err
		}

		var shortUrl string
		if !db.ExistLongUrl(body.Url) {
			shortUrl, _ = ShortenUrlKeygen(body.Url, urlPrefix, db)
			db.StoreShortUrl(shortUrl, body.Url)
			log.Debug().Msgf("URL not found: %s", body.Url)

		} else {
			shortUrl = db.RetrieveShortUrl(body.Url)
			log.Debug().Msgf("URL already in DB: %s", shortUrl)
		}

		return c.Status(fiber.StatusOK).SendString(shortUrl)
	})

	app.Use("/", func(c *fiber.Ctx) error {

		var shortUrl string = c.Context().URI().String()

		if db.ExistShortUrl(shortUrl) {
			var longUrl = db.RetrieveLongUrl(shortUrl)
			log.Debug().Msgf("Long URL: found: %s", longUrl)
			return c.Status(fiber.StatusOK).Redirect(longUrl)
		} else {
			return c.Status(fiber.StatusNotFound).SendString("The requested URL does not exist!")
		}
	})

	return app
}

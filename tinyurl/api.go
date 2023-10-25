package main

import (
	"github.com/rs/zerolog/log"
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
		return c.Status(fiber.StatusOK).SendString("TinyURL is healthy")
	})

	// save a new URL and its corresponding shortened version
	app.Post("/shorten", func(c *fiber.Ctx) error {

		body := new(UrlToShorten)
		if err := c.BodyParser(body); err != nil {
			return err
		}

		var shortUrl string
		if !UrlDB.ExistsInverse(body.Url) {
			shortUrl, _ = ShortenUrlKeygen(body.Url, urlPrefix)
			log.Debug().Msgf("URL % s not found", body.Url)
			UrlDB.Insert(shortUrl, body.Url)

		} else {
			shortUrl, _ = UrlDB.GetInverse(body.Url)
			log.Debug().Msgf("Already in DB: %s", shortUrl)
		}

		return c.Status(fiber.StatusOK).SendString(shortUrl)
	})

	app.Use("/", func(c *fiber.Ctx) error {

		var shortUrl string = c.Context().URI().String()

		if UrlDB.Exists(shortUrl) {
			longUrl, _ := UrlDB.Get(shortUrl)
				log.Debug().Msgf("Long URL: %s", longUrl)
			return c.Status(fiber.StatusOK).Redirect(longUrl)
		} else {
			return c.Status(fiber.StatusNotFound).SendString("The requested URL does not exist!")
		}
	})

	return app
}

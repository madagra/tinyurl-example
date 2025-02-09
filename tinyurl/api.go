package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type UrlToShorten struct {
	Url string `json:"url" xml:"url" form:"url"`
}

func shorten(app *fiber.App, urlPrefix string, dbManager DBManager) fiber.Router {

	internal := func(c *fiber.Ctx) error {
		body := new(UrlToShorten)
		if err := c.BodyParser(body); err != nil {
			return err
		}

		var shortUrl string
		if !dbManager.ExistLongUrl(body.Url) {
			shortUrl, _ = ShortenUrlKeygen(body.Url, urlPrefix, dbManager)
			dbManager.StoreShortUrl(shortUrl, body.Url)
			log.Debug().Msgf("URL not found in DB, create new shortened URL: %s", body.Url)

		} else {
			shortUrl = dbManager.RetrieveShortUrl(body.Url)
			log.Debug().Msgf("URL already in DB: %s", shortUrl)
		}

		return c.Status(fiber.StatusOK).SendString(shortUrl)
	}

	return app.Post("/shorten", internal)

}

func redirect(app *fiber.App, dbManager DBManager) fiber.Router {

	internal := func(c *fiber.Ctx) error {

		var shortUrl string = c.Context().URI().String()

		if dbManager.ExistShortUrl(shortUrl) {
			var longUrl = dbManager.RetrieveLongUrl(shortUrl)
			log.Debug().Msgf("Long URL: found: %s", longUrl)
			return c.Status(fiber.StatusOK).Redirect(longUrl)
		} else {
			return c.Status(fiber.StatusNotFound).SendString("The requested URL does not exist!")
		}
	}

	return app.Use("/", internal)
}

func health(app *fiber.App) fiber.Router {

	internal := func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("TinyURL is healthy")
	}

	return app.Get("/", internal)
}

func CreateServer(urlPrefix string, db DBManager) (*fiber.App, []fiber.Router) {

	app := fiber.New()

	// create routes
	healthRouter := health(app)
	shortenRouter := shorten(app, urlPrefix, db)
	redirectRouter := redirect(app, db)

	routers := []fiber.Router{healthRouter, shortenRouter, redirectRouter}
	return app, routers
}

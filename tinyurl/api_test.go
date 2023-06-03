package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

var urlToShorten string = "https://learn.cantrill.io/courses/enrolled/1820301"

var TestUrlPrefix = "http://example.com/"

func TestHealthEndpoint(t *testing.T) {
	app := CreateServer(TestUrlPrefix)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode == fiber.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if !(string(body) == "TinyURL is healthy") {
			t.Errorf("Wrong return string: %s", string(body))
		}
	} else {
		t.Errorf("Health endpoint not working")
	}
}

func TestShortenEndpoint(t *testing.T) {
	app := CreateServer(TestUrlPrefix)
	req := httptest.NewRequest(http.MethodPost, "/shorten", nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Form = make(url.Values)
	req.Form.Add("url", urlToShorten)
	resp, _ := app.Test(req)

	if resp.StatusCode != fiber.StatusOK || len(UrlKeysDB) != 1 {
		t.Errorf("Database has not been updated correctly")
	}

	t.Cleanup(PurgeUrlDB)
}

func TestRedirectEndpoint(t *testing.T) {

	var shortenedUrl, _ = ShortenUrlKeygen(urlToShorten, TestUrlPrefix)
	UrlDB.Insert(shortenedUrl, urlToShorten)

	app := CreateServer(TestUrlPrefix)

	var urlKey string = strings.Split(strings.ReplaceAll(shortenedUrl, TestUrlPrefix, ""), ".")[0]
	req := httptest.NewRequest(http.MethodGet, "/"+urlKey, nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != fiber.StatusFound {
		t.Errorf("Redirection has not been successful")
	}

	t.Cleanup(PurgeUrlDB)
}

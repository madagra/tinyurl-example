package main

import (
	"testing"
)

var exampleUrls = []string{
	"https://learn.cantrill.io/courses/enrolled/1820301",
	"https://portal.tutorialsdojo.com/courses/aws-certified-solutions-architect-associate-practice-exams/lessons/practice-exams-review-mode-4/quizzes/aws-certified-solutions-architect-associate-practice-exam-review-mode-set-2/",
	"https://aws.amazon.com/blogs/quantum-computing/quantum-chemistry-with-qucos-qubec-on-amazon-braket/",
	"https://github.com/golang/tools/blob/master/gopls/doc/workspace.md",
	"https://portal.tutorialsdojo.com/courses/aws-certified-solutions-architect-associate-practice-exams/lessons/practice-exams-review-mode-4/quizzes/aws-certified-solutions-architect-associate-practice-exam-review-mode-set-5/",
	"https://www.simonwenkel.com/notes/software_libraries/tinygrad/introducing-tinygrad.html",
}

func TestShortenEncoding(t *testing.T) {

	var dbClient = initLocalTestDb(t)

	for _, url := range exampleUrls {

		shortUrl, _ := ShortenUrlEncoding(url, LocalUrlPrefix, dbClient)

		if len(shortUrl) > len(LocalUrlPrefix)+lenShortUrl {
			t.Errorf("Encoding URL did not work: %s", shortUrl)
		}

	}

	if len(dbClient.urlKeysDB) != len(exampleUrls) {
		t.Errorf("Database has not been updated correctly: %d != %d", len(dbClient.urlKeysDB), len(exampleUrls))
	}

}

func TestShortenKeygen(t *testing.T) {

	var dbClient = initLocalTestDb(t)

	for _, url := range exampleUrls {

		shortUrl, _ := ShortenUrlKeygen(url, LocalUrlPrefix, dbClient)
		if len(shortUrl) > len(LocalUrlPrefix)+lenShortUrl {
			t.Errorf("Encoding URL did not work: %s", shortUrl)
		}
	}

	if len(dbClient.urlKeysDB) != len(exampleUrls) {
		t.Errorf("Database has not been updated correctly: %d != %d", len(dbClient.urlKeysDB), len(exampleUrls))
	}

}

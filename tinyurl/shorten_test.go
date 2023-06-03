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
	for _, url := range exampleUrls {

		shortUrl, _ := ShortenUrlEncoding(url, LocalUrlPrefix)

		if len(shortUrl) > len(LocalUrlPrefix)+LenShortUrl {
			t.Errorf("Encoding URL did not work: %s", shortUrl)
		}

	}

	if len(UrlKeysDB) != len(exampleUrls) {
		t.Errorf("Database has not been updated correctly: %d != %d", len(UrlKeysDB), len(exampleUrls))
	}

	t.Cleanup(PurgeUrlDB)
}

func TestShortenKeygen(t *testing.T) {
	for _, url := range exampleUrls {

		shortUrl, _ := ShortenUrlKeygen(url, LocalUrlPrefix)
		if len(shortUrl) > len(LocalUrlPrefix)+LenShortUrl {
			t.Errorf("Encoding URL did not work: %s", shortUrl)
		}
	}

	if len(UrlKeysDB) != len(exampleUrls) {
		t.Errorf("Database has not been updated correctly: %d != %d", len(UrlKeysDB), len(exampleUrls))
	}

	t.Cleanup(PurgeUrlDB)
}

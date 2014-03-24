package main

import (
	"log"
	"github.com/PuerkitoBio/goquery"
	"errors"
)

func parseTeam(url string) (int, error) {
	d, err := getDocument(url)
	if err != nil {
		log.Printf("could not parse team %s: %s", url, err)
	}

	title := d.Find("h1").Text()

	country := ""
	d.Find("div dl dt").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if s.Text() == "Country" {
			country = s.Next().Text()
			return false
		}

		return true
	})

	if country == "" {
		return 0, errors.New("could not find country of team")
	}

	// TODO: Add team
	_ = title

	return 1, nil
}

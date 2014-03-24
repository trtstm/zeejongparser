package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
)

func parseSeason(title, url string, competitionId int) {
	d, err := getDocument(url)
	if err != nil {
		log.Printf("could not parse season %s: %s", url, err)
		return
	}

	// TODO: Add season
	seasonId := 1
	_ = seasonId

	d.Find("td.score a").Each(func(i int, s *goquery.Selection) {
		url, ok := s.Attr("href")
		if !ok {
			log.Printf("could not get url for match")

			return
		}

		parseMatch(BASE + url, competitionId, seasonId)
	})
}

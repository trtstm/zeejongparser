package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
)

func getSeasons(d *goquery.Document) map[string]string {
	seasons := map[string]string{}

	d.Selection.Find(".season a").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		url, ok := s.Attr("href")
		if !ok {
			log.Printf("could not get url for season")

			return
		}

		seasons[title] = url
	})

	return seasons
}

func parseCompetition(url string) {
	d, err := getDocument(url)
	if err != nil {
		log.Printf("could not parse competition %s: %s", url, err)
		return
	}

	h1 := d.Find("h1")
	if len(h1.Children().Nodes) == 0 {
		log.Printf("could not find title for competition %s", url)
		return
	}

	competition := h1.Children().Get(0).NextSibling.Data
	competitionId := addCompetition(competition)

	seasons := getSeasons(d)
	for title, url := range seasons {
		parseSeason(title, BASE+url, competitionId)
	}

}

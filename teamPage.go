package main

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"log"
	"sync"
)

var counter = map[string]bool{}
var counterLock = sync.RWMutex{}


/*
Parse the team with the given url
*/
func parseTeam(url string) (int, error) {
	if id, ok := getUrlFromCache(url); ok {
		return id, nil
	}

	d, err := getDocument(url)
	if err != nil {
		log.Printf("could not parse team %s: %s", url, err)
	}

	name := d.Find("h1").Text()

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

	countryId := addCountry(country)
	id := addTeam(name, countryId, url)

	// Get image
	imgSrc, ok := d.Find(".block_team_info img").Attr("src")
	if ok {
		getImage(imgSrc, "Team", id)
	}

	counterLock.RLock()
	if _, ok := counter[url]; !ok {
		counterLock.RUnlock()

		counterLock.Lock()
		counter[url] = true

		d.Find(".table.squad div a").Each(func(i int, s *goquery.Selection) {
			playerUrl, ok := s.Attr("href")
			if !ok {
				return
			}

			playerId, err := parsePlayer(BASE + playerUrl)
			if err != nil {
				return
			}

			addPlaysIn(id, playerId)
		})
		counterLock.Unlock()

	} else {
		counterLock.RUnlock()
	}

	return id, nil
}

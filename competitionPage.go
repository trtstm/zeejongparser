package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"sync"
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
	if _, ok := getUrlFromCache(url); ok {
		return
	}

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
	competitionId := addCompetition(competition, url)

	seasons := getSeasons(d)
	
	wg := sync.WaitGroup{}
	
	for title, url := range seasons {
	
		wg.Add(1)
		go func(title string, url string, id int){
			defer wg.Done()
			parseSeason(title, url, competitionId)
		}(title, BASE + url, competitionId)
	
	}
	
	wg.Wait()

}

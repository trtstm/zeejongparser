package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"sync"
)


/*
Get all the season from the competition's archive page
*/
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

/*
Parse the competition with the given url
*/
func parseCompetition(url, name string) {
	if _, ok := getUrlFromCache(url); ok {
		return
	}

	d, err := getDocument(url)
	if err != nil {
		log.Printf("could not parse competition %s: %s", url, err)
		return
	}

	competitionId := addCompetition(name, url)

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

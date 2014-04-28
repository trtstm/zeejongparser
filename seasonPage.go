package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
)

func getSubSeasons(d *goquery.Document) []string {
	urls := []string{}

	d.Find(".left-tree li.expanded ul.level-1 .leaf a").Each(func(i int, s *goquery.Selection) {
		url, ok := s.Attr("href")
		if !ok {
			log.Printf("could not get url for match")

			return
		}

		subd, err := getDocument(BASE + url)
		if err != nil {
			log.Printf("could not parse sub season %s: %s", url, err)
			return
		}

		suburl, ok := subd.Find("#submenu ul li").Next().Find("a").Attr("href")
		if !ok {
			log.Printf("could not find sub season matches %s: %s", url, err)
			return
		}

		urls = append(urls, suburl)
	})

	return urls
}

func parseSeason(name, url string, competitionId int) {
	if _, ok := getUrlFromCache(url); ok {
		return
	}

	d, err := getDocument(url)
	if err != nil {
		log.Printf("could not parse season %s: %s", url, err)
		return
	}

	seasonId := addSeason(name, competitionId, url)

	for _, suburl := range getSubSeasons(d) {
		subd, err := getDocument(BASE + suburl)
		if err != nil {
			log.Printf("could not parse subseason %s: %s", suburl, err)
			continue
		}

		subd.Find("td.score a").Each(func(i int, s *goquery.Selection) {
			url, ok := s.Attr("href")
			if !ok {
				log.Printf("could not get url for match")
				return
			}

			parseMatch(BASE+url, competitionId, seasonId)
		})
	}
}

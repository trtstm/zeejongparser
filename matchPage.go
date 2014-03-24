package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"errors"
)

func getTeamsId(d *goquery.Document) ([2]int, error) {
	var ids [2]int

	url1, ok := d.Find("div.container.left h3 a").Attr("href")
	if !ok {
		return ids, errors.New("could not find team a")
	}

	idA, err := parseTeam(BASE + url1)
	if err != nil {
		return ids, err
	}

	url2, ok := d.Find("div.container.right h3 a").Attr("href")
	if !ok {
		return ids, errors.New("could not find team b")
	}

	idB, err := parseTeam(BASE + url2)
	if err != nil {
		return ids, err
	}

	ids[0] = idA
	ids[1] = idB
	return ids, nil
}

func parseMatch(url string, competitionId, seasonId int) {
	d, err := getDocument(url)
	if err != nil {
		log.Printf("could not parse match %s: %s", url, err)
		return
	}

	teamsId, err := getTeamsId(d)
	if err != nil {
		log.Printf("could not find teams in match")
		return
	}

	_ = teamsId
}

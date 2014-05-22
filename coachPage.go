package main

import (
	"log"
	"strings"
	"errors"
)


/*
Parse the coach with the given url
*/
func parseCoach(url string) (int, error) {
	if !strings.Contains(url, "/coaches/") {
		return 0, errors.New("Coach has none coach url")
	}

	if id, ok := getUrlFromCache(url); ok {
		return id, nil
	}

	d, err := getDocument(url)
	if err != nil {
		log.Printf("could not parse coach %s: %s", url, err)
		return 0, err
	}

	info := d.Selection.Find(".content .first dd")

	firstname, err := info.Html()
	if err != nil {
		log.Printf("Could not get firstname for coach %s\n", url)
		return 0, err
	}

	info = info.Next().Next()
	lastname, err := info.Html()
	if err != nil {
		log.Printf("Could not get lastname for coach %s\n", url)
		return 0, err
	}

	info = info.Next().Next()
	country, err := info.Html()
	if err != nil {
		log.Printf("Could not get country for coach %s\n", url)
		return 0, err
	}

	countryId := addCountry(country)
	id := addCoach(firstname, lastname, countryId, url)

	// Get image
	imgSrc, ok := d.Find(".block_player_passport img").Attr("src")
	if ok {
		getImage(imgSrc, "Coach", id)
	}

	return id, nil
}

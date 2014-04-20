package main

import (
	"log"
)

func parseReferee(url string) (int, error) {
	d, err := getDocument(url)
	if err != nil {
		log.Printf("could not parse referee %s: %s", url, err)
		return 0, err
	}

	info := d.Selection.Find(".content .first dd")

	firstname, err := info.Html()
	if err != nil {
		log.Printf("Could not get firstname for referee %s\n", url)
		return 0, err
	}

	info = info.Next().Next()
	lastname, err := info.Html()
	if err != nil {
		log.Printf("Could not get lastname for referee %s\n", url)
		return 0, err
	}

	info = info.Next().Next()
	country, err := info.Html()
	if err != nil {
		log.Printf("Could not get country for referee %s\n", url)
		return 0, err
	}

	countryId := addCountry(country)
	id := addReferee(firstname, lastname, countryId)

	// Get image
	imgSrc, ok := d.Find(".block_player_passport img").Attr("src")
	if ok {
		getImage(imgSrc, "Referee", id)
	}

	return id, nil
}

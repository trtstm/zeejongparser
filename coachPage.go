package main

import (
	"log"
)

func parseCoach(url string) (int, error) {
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

	_ = country
	id := addCoach(firstname, lastname, 0)

	return id, nil
}

package main

import (
	"log"
)

func parsePlayer(url string) (int, error) {
	d, err := getDocument(url)
	if err != nil {
		log.Printf("could not parse referee %s: %s", url, err)
		return 0, err
	}

	infoTag := d.Selection.Find(".content .first dd")

	firstname, err := infoTag.Html()
	if err != nil {
		log.Printf("Could not find firstname for player %s\n", url)
		return 0, err
	}

	infoTag = infoTag.Next().Next()
	lastname, err := infoTag.Html()
	if err != nil {
		log.Printf("Could not find lastname for player %s\n", url)
		return 0, err
	}

	infoTag = infoTag.Next().Next()
	country, err := infoTag.Html()
	if err != nil {
		log.Printf("Could not find country for player %s\n", url)
		return 0, err
	}

	// TODO: Add player, country
	_ = firstname
	_ = lastname
	_ = country

	return 1, nil
}

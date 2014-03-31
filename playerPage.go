package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"strings"
)

func removeNonAlpha(r rune) bool {
	if r >= '0' && r <= '9' {
		return false
	}

	return true
}

func parsePlayer(url string) (int, error) {
	d, err := getDocument(url)
	if err != nil {
		log.Printf("could not parse referee %s: %s", url, err)
		return 0, err
	}

	//Info contains the attributes of the player
	//e.g. info['Position'] = "Defender"
	info := map[string]string{}

	d.Find(".content .first dl").Children().EachWithBreak(func(i int, s *goquery.Selection) bool {

		//We only check the even elements (= dd elements)
		if i%2 == 0 {

			key := s.Text()
			key = strings.ToLower(strings.TrimSpace(key))

			valueSelection := s.Next()

			//Check if the given dd element exists for this dt element
			if valueSelection.Length() < 1 {
				log.Printf("Could not find metadata for player (unmatching dt/dd) %s\n", url)
				return false
			}

			value := valueSelection.Text()
			info[key] = value

		}

		return true
	})

	firstname := info["first name"]
	lastname := info["last name"]
	nationality := 0
	dateOfBirth := 0
	height, _ := strconv.Atoi(strings.TrimFunc(info["height"], removeNonAlpha))
	weight, _ := strconv.Atoi(strings.TrimFunc(info["weight"], removeNonAlpha))
	position := info["position"]

	id := addPlayer(firstname, lastname, nationality, dateOfBirth, height, weight, position)

	return id, nil
}

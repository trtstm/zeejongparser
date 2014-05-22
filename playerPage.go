package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"strings"
	"time"
)

func removeNonAlpha(r rune) bool {
	if r >= '0' && r <= '9' {
		return false
	}

	return true
}


/*
Parse a data of the format `2 January 2006`
Returns the unix time stamp
*/
func parseDate(date string) (int, error) {

	datetime, err := time.Parse("2 January 2006", date)
	if err != nil {
		return 0, err
	}
	return int(datetime.Unix()), nil
	
}


/*
Parse the player with the given url
*/
func parsePlayer(url string) (int, error) {
	if id, ok := getUrlFromCache(url); ok {
		return id, nil
	}

	d, err := getDocument(url)
	if err != nil {
		log.Printf("Could not parse player %s: %s", url, err)
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
	nationality := addCountry(info["nationality"])
	dateOfBirth, err := parseDate(info["date of birth"])
	if err != nil {
		//log.Printf("Could not parse date %s %s: %s", url, info["date of birth"])
	}
	height, _ := strconv.Atoi(strings.TrimFunc(info["height"], removeNonAlpha))
	weight, _ := strconv.Atoi(strings.TrimFunc(info["weight"], removeNonAlpha))
	position := info["position"]

	id := addPlayer(firstname, lastname, nationality, dateOfBirth, height, weight, position, url)

	// Get image
	imgSrc, ok := d.Find(".block_player_passport img").Attr("src")
	if ok {
		getImage(imgSrc, "Player", id)
	}

	return id, nil
}

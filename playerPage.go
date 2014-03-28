package main

import (
	"log"
	"github.com/PuerkitoBio/goquery"
)

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
		if(i % 2 == 0) {
		
			key := s.Text()
			valueSelection := s.Next()
			
			//Check if the given dd element exists for this dt element
			if(valueSelection.Length() < 1) {
				log.Printf("Could not find metadata for player (unmatching dt/dd) %s\n", url)
				return false;
			}
			
			value := valueSelection.Text()
			info[key] = value
			
		}
		
		return true;
	})
	
	
	_ = info
	
	return 1, nil
}

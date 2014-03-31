package main

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"strings"
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

func getScoreId(d *goquery.Document) (int, error) {
	score := d.Selection.Find("h1 .bidi").Text()

	scoreSplit := strings.Split(score, "-")
	if len(scoreSplit) != 2 {
		log.Println("Could not parse score for match")
		return 0, errors.New("could not find - symbol in score")
	}

	scoreSplit[0] = strings.TrimSpace(scoreSplit[0])
	scoreSplit[1] = strings.TrimSpace(scoreSplit[1])

	scoreA, err := strconv.Atoi(scoreSplit[0])
	if err != nil {
		log.Println("Could not parse score for match")
		return 0, err
	}

	scoreB, err := strconv.Atoi(scoreSplit[1])
	if err != nil {
		log.Println("Could not parse score for match")
		return 0, err
	}

	// TODO: Add score
	_ = scoreA
	_ = scoreB
	scoreId := 1

	return scoreId, nil
}

func getDate(d *goquery.Document) (int, error) {
	date := 0
	found := false

	d.Selection.Find(".block_match_info .details dt").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if s.Text() == "Date" {
			dateString, ok := s.Next().Find("span").Attr("data-value")
			if !ok {
				log.Printf("could not parse date for match")
				return false
			}

			var err error
			date, err = strconv.Atoi(dateString)
			if err != nil {
				log.Printf("could not parse date for match")
				return false
			}

			found = true

			return false
		}

		return true
	})

	if !found {
		return 0, errors.New("did not find date")
	}

	return date, nil
}

func getReferee(d *goquery.Document) (int, error) {
	refereeUrl, ok := d.Selection.Find(".referee").Attr("href")
	if !ok {
		return 0, errors.New("could not find referee for match")
	}

	refereeId, err := parseReferee(BASE + refereeUrl)
	if err != nil {
		return 0, err
	}

	return refereeId, nil
}

type Player struct {
	Url         string
	Shirtnumber int
}

func getPlayers(d *goquery.Document) []Player {
	players := []Player{}

	cols := []string{".left", ".right"}

	for _, col := range cols {
		playerTag := d.Selection.Find(".block_match_lineups " + col + " td.player")
		playerTag.Each(func(i int, s *goquery.Selection) {
			shirtnumber := -1
			shirtnumberString, err := s.Parent().Find(".shirtnumber").Html()
			if err != nil {
				log.Println("Could not find player shirtnumber for match")
			} else {
				shirtnumberString = strings.TrimSpace(shirtnumberString)

				shirtnumber, err = strconv.Atoi(shirtnumberString)
				if err != nil {
					shirtnumber = -1
				}
			}

			playerUrl, ok := s.Parent().Find(".player a").Attr("href")
			if !ok {
				log.Println("Could not find player url for match")
				return
			}

			p := Player{Url: playerUrl, Shirtnumber: shirtnumber}
			players = append(players, p)
		})
	}

	return players
}

type Coach struct {
	Url string
}

func getCoaches(d *goquery.Document) []Coach {
	coaches := []Coach{}

	cols := []string{".left", ".right"}

	for _, col := range cols {
		tmp := d.Selection.Find(".block_match_lineups " + col + " tbody").Find("tr a")
		coachUrl, ok := tmp.Eq(-1).Attr("href")
		if !ok {
			log.Println("Could not find coach url for match ")
			continue
		}

		c := Coach{Url: coachUrl}
		coaches = append(coaches, c)
	}

	return coaches
}

func parseMatch(url string, competitionId, seasonId int) {
	d, err := getDocument(url)
	if err != nil {
		log.Printf("could not parse match %s: %s", url, err)
		return
	}

	teamsId, err := getTeamsId(d)
	if err != nil {
		log.Printf("could not find teams in match: %s", url)
		return
	}

	scoreId, err := getScoreId(d)
	if err != nil {
		log.Printf("could not find score in match")
		return
	}

	date, err := getDate(d)
	if err != nil {
		log.Printf("could not find date in match")
		return
	}

	refereeId, err := getReferee(d)
	if err != nil {
		log.Printf("could not find referee in match")
		return
	}

	for _, p := range getPlayers(d) {
		playerId, err := parsePlayer(BASE + p.Url)
		if err != nil {
			log.Printf("could not get player id in match")
			continue
		}

		// TODO: Add player to match

		_ = playerId
	}

	for _, c := range getCoaches(d) {
		coachId, err := parseCoach(BASE + c.Url)
		if err != nil {
			log.Printf("could not find coach in match")
			continue
		}

		
		// TODO: Add coach to coaches team.
		_ = coachId
	}

	_ = date

	_ = scoreId

	_ = teamsId

	_ = refereeId
}

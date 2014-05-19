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

	scoreId := addScore(scoreA, scoreB)

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
	Team int
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

			team := 0
			if col == ".right" {
				team = 1
			}

			p := Player{Url: playerUrl, Team: team, Shirtnumber: shirtnumber}
			players = append(players, p)
		})
	}
	
	
	
	
	for _, col := range cols {
		playerTag := d.Selection.Find(".block_match_substitutes " + col + " td.player")
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
	
			team := 0
			if col == ".right" {
				team = 1
			}
	
			p := Player{Url: playerUrl, Team: team, Shirtnumber: shirtnumber}
			players = append(players, p)
		})
	}
	
	
	
	
	

	return players
}



func parseCards(d *goquery.Document, matchId int) {
	
	
	cols := []string{".left", ".right"}
	
	for _, col := range cols {
		playerTag := d.Selection.Find(".block_match_lineups " + col + " td.player")
		playerTag.Each(func(i int, s *goquery.Selection) {
			
			playerUrl, ok := s.Parent().Find(".player a").Attr("href")
			if !ok {
				return
			}
			playerId, err := parsePlayer(BASE + playerUrl)
			if err != nil {
				return
			}
			
			
			s.Parent().Find(".bookings span").Each(func(i int, s *goquery.Selection) {
			
				bookingUrl, ok := s.Find("img").Attr("src")
				cardType := 1
				
				if ok {
					
					if strings.Contains(bookingUrl, "YC.png") {
						cardType = 1
					} else if strings.Contains(bookingUrl, "RC.png") {
						cardType = 2
					} else if strings.Contains(bookingUrl, "Y2C.png") {
						cardType = 3
					} else {
						return
					}
					
					
					stringTime := strings.TrimSpace(s.Text())
					
					if err == nil {
						rawTime := strings.Split(strings.TrimSuffix(stringTime, "'"), "+")
						
						var time int64 = 0
						
						for _,item := range rawTime {
							num, err := strconv.ParseInt(item, 10, 16)
							if err == nil {
								time += num
							}
						}
						
						addCard(playerId, matchId, int(time), cardType)
					}
					
					
					
					
				}
			
			
				
			
			})
		
		})

	}
	
}



type Coach struct {
	Url string
	Team int
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

		team := 0
		if col == ".right" {
			team = 1
		}

		c := Coach{Url: coachUrl, Team: team}
		coaches = append(coaches, c)
	}

	return coaches
}



type Goal struct {
	Time int64
	PlayerUrl string
	MatchId int
}


func parseGoals(d *goquery.Document, matchId int) {
	
	d.Find(".block_match_goals-wrapper .event").EachWithBreak(func(i int, s *goquery.Selection) bool {
	
		goal := Goal{}
		goal.MatchId = matchId
	
		//Find the time of the goal
		rawTime := strings.Split(s.Find(".minute").Text(), "+")
		rawTime[0] = strings.TrimSuffix(rawTime[0], "'")
		
		var time int64 = 0
		
		for _,item := range rawTime {
			num, err := strconv.ParseInt(item, 10, 16)
			
			if err != nil {
				return false
			}
			
			time += num
		}
		
		goal.Time = time;

		
		
		
		//Find the player url for the goal
		//TODO handle error
		goal.PlayerUrl,_ = s.Find("a").Attr("href")
		goal.PlayerUrl = BASE + goal.PlayerUrl
		
		playerId, _ := parsePlayer(goal.PlayerUrl)
		
		
		
		//TODO add the goal to the database
		addGoal(playerId, matchId, int(time))
		
		return true
		
	})
	
}

func elemInArray(array []string, needle string) bool {
	for _, elem := range array {
		if needle == elem {
			return true
		}
	}

	return false
}

func parseMatch(url string, competitionId, seasonId int, finalType string) {
	if _, ok := getUrlFromCache(url); ok {
		return
	}

	finals := []string{"Final", "Finals", "Final replay", "Play-offs - Final", "Conference - Finals", "Europe League Play-offs - Finals", "Play Offs UEFA Final",
				 	"Semi-finals", "Play-offs - Semi-finals", "Conference - Semi-finals", "Europa League Play-offs - Semi-finals",
					"3rd Place Final",
					"Quarter-finals", "Quarter-finals Replays",
					"8th Finals",
					"16th Finals",
					"32th Finals",
					""}

	finalType = strings.TrimSpace(finalType)
	if !elemInArray(finals, finalType) {
		log.Printf("Found unknown final type: %s", finalType)
		finalType = "";
	}

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

	refereeId, err := getReferee(d)
	if err != nil {
		refereeId = 0
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

	matchId := addMatch(teamsId[0], teamsId[1], seasonId, refereeId, date, scoreId, url, finalType)

	for _, p := range getPlayers(d) {
		playerId, err := parsePlayer(BASE + p.Url)
		if err != nil {
			log.Printf("could not get player id in match")
			continue
		}

		addPlaysMatchInTeam(playerId, p.Shirtnumber, teamsId[p.Team], matchId)
	}

	for _, c := range getCoaches(d) {
		coachId, err := parseCoach(BASE + c.Url)
		if err != nil {
			log.Printf("could not find coach in match")
			continue
		}

		
		addCoacheses(coachId, teamsId[c.Team], matchId)
	}
	
	
	//Parse the goals
	parseGoals(d, matchId)
	
	//Parse the cards
	parseCards(d, matchId)
	
	_ = date

	_ = scoreId

	_ = teamsId
	
}

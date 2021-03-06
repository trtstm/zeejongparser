package main

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"strings"
)

/*
Get the two teams in a match
*/
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


/*
Get the score id for this match
*/
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

/*
Get the date for this match
*/
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

/*
Get the referee for this match
*/
func getReferee(d *goquery.Document) (int, error) {
	refereeUrl, ok := d.Selection.Find(".referee").Attr("href")
	if !ok {
		return 0, errors.New("could not find referee for match")
	}

	if !strings.HasPrefix(refereeUrl, "/referees/") {
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


/*
Get the players for this match
*/
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

			if !strings.HasPrefix(playerUrl, "/players/") {
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

			if !strings.HasPrefix(playerUrl, "/players/") {
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



/*
Parse the cards in the match,
and adds them to the database
*/
func parseCards(d *goquery.Document, matchId int) {
	
	
	cols := []string{".left", ".right"}
	
	for _, col := range cols {
		playerTag := d.Selection.Find(".block_match_lineups " + col + " td.player")
		playerTag.Each(func(i int, s *goquery.Selection) {
			
			playerUrl, ok := s.Parent().Find(".player a").Attr("href")
			if !ok {
				return
			}

			if !strings.HasPrefix(playerUrl, "/players/") {
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


/*
Get the coaches for this match
*/
func getCoaches(d *goquery.Document) []Coach {
	coaches := []Coach{}

	cols := []string{".left", ".right"}

	for _, col := range cols {
		tmp := d.Selection.Find(".block_match_lineups " + col + " tbody").Find("tr strong")
		if strings.Trim(tmp.Text(), " ") == "Coach:" {
			aTag := tmp.Next()
			if url, ok := aTag.Attr("href"); ok {
				team := 0
				if col == ".right" {
					team = 1
				}

				if !strings.HasPrefix(url, "/coaches/") {
					continue
				}	

				c := Coach{Url: url, Team: team}
				coaches = append(coaches, c)
			}
		}

	}

	return coaches
}



type Goal struct {
	Time int64
	PlayerUrl string
	MatchId int
}


/*
Parse the goals for this match,
and adds them to the database
*/
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
				return true
			}
			
			time += num
		}
		
		goal.Time = time;
	
		
		//Find the player url for the goal
		var ok bool
		goal.PlayerUrl, ok = s.Find("a").Attr("href")
		if !ok {
			return true
		}

		if !strings.HasPrefix(goal.PlayerUrl, "/players/") {
			return true
		}

		goal.PlayerUrl = BASE + goal.PlayerUrl
		
		playerId, err := parsePlayer(goal.PlayerUrl)
		if err != nil {
			log.Println(err)
			return true
		}
		
		
		//TODO add the goal to the database
		addGoal(playerId, matchId, int(time))
		
		return true
		
	})
	
}

/*
Parse the match with the given url
*/
func parseMatch(url string, competitionId, seasonId int, finalType string) {
	if _, ok := getUrlFromCache(url); ok {
		return
	}

	finalType = strings.TrimSpace(finalType)

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
		//log.Println(c.Url)
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

package main

import (
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"strconv"
	"strings"
	"encoding/json"
	"sync"
	"io/ioutil"
)

type player struct {
	Id          int
	Firstname   string
	Lastname    string
	Country     int
	DateOfBirth int
	Height      int
	Weight      int
	Position    string
}

type referee struct {
	Id          int
	Firstname   string
	Lastname    string
	Country     int
}

type coach struct {
	Id          int
	Firstname   string
	Lastname    string
	Country     int
}

type country struct {
	Id int
	Name string
}

type match struct {
	Id int
	TeamA int
	TeamB int
	Season int
	Referee int
	Date int
	Score int
}

type playsMatchInTeam struct {
	Id int
	PlayerId int
	Number int
	TeamId int
	MatchId int
}

type playsIn struct {
	Id int
	TeamId int
	PlayerId int
}

type score struct {
	Id int
	TeamA int
	TeamB int
}

type coacheses struct {
	Id int
	CoachId int
	TeamId int
	MatchId int
}

type team struct {
	Id int
	Name string
	CountryId int
}

type season struct {
	Id int
	Name string
	CompetitionId int
}

type competition struct {
	Id int
	Name string
}

type goal struct {
	Id int
	MatchId int
	PlayerId int
	Time int
}

type card struct {
	Id int
	MatchId int
	PlayerId int
	Time int
	Type int
}


var db = struct {
	
	
	dbLock sync.RWMutex

	Players     map[string]player

	Referees     map[string]referee

	Coaches     map[string]coach

	Countries     map[string]country

	Matches     map[string]match

	Scores     map[string]score

	PlaysMatchInTeams     map[string]playsMatchInTeam

	PlaysIn     map[string]playsIn

	Coacheses     map[string]coacheses

	Teams     map[string]team

	Seasons     map[string]season

	Competitions     map[string]competition
	
	Goals map[string]goal
	
	Cards map[string]card
	

	
}{
	Players: map[string]player{},
	Referees: map[string]referee{},
	Coaches: map[string]coach{},
	Countries: map[string]country{},
	Matches: map[string]match{},
	Scores: map[string]score{},
	PlaysMatchInTeams: map[string]playsMatchInTeam{},
	PlaysIn: map[string]playsIn{},
	Coacheses: map[string]coacheses{},
	Teams: map[string]team{},
	Seasons: map[string]season{},
	Competitions: map[string]competition{},
	Goals: map[string]goal{},
	Cards: map[string]card{},
}

func writeDb(file string) error {
	json, err := json.MarshalIndent(db, "", "    ") 
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, json, 0777)
}

func getDbSize() map[string]int {



	return map[string]int {
		"Players": len(db.Players),
		"Referees": len(db.Referees),
		"Coaches": len(db.Coaches),
		"Countries": len(db.Countries),
		"Matches": len(db.Matches),
		"Scores": len(db.Scores),
		"PlaysMatchInTeam": len(db.PlaysMatchInTeams),
		"PlaysIn": len(db.PlaysIn),
		"Coacheses": len(db.Coacheses),
		"Teams": len(db.Teams),
		"Seasons": len(db.Seasons),
		"Competitions": len(db.Competitions),
		"Goals": len(db.Goals),
		"Cards": len(db.Cards),
	}
}

func addPlayer(firstname, lastname string, countryId, dateOfBirth, height, weight int, position string) int {
	
	db.dbLock.Lock()
	defer db.dbLock.Unlock()
	
	firstname = strings.TrimSpace(firstname)
	lastname = strings.TrimSpace(lastname)
	position = strings.TrimSpace(position)

	hash := getHash(firstname, lastname, countryId, dateOfBirth)

	if player, ok := db.Players[hash]; ok {
		return player.Id
	}


	id := len(db.Players) + 1
	db.Players[hash] = player{Id: id, Firstname: firstname,
		Lastname: lastname, Country: countryId, DateOfBirth: dateOfBirth,
		Height: height, Weight: weight, Position: position}

	return id
}

func addReferee(firstname, lastname string, countryId int) int {
	
	db.dbLock.Lock()
	defer db.dbLock.Unlock()
	
	firstname = strings.TrimSpace(firstname)
	lastname = strings.TrimSpace(lastname)

	hash := getHash(firstname, lastname, countryId)

	if referee, ok := db.Referees[hash]; ok {
		return referee.Id
	}



	id := len(db.Referees) + 1
	db.Referees[hash] = referee{Id: id, Firstname: firstname, Lastname: lastname, Country: countryId}


	return id
}

func addCoach(firstname, lastname string, countryId int) int {
	
	db.dbLock.Lock()
	defer db.dbLock.Unlock()
	
	firstname = strings.TrimSpace(firstname)
	lastname = strings.TrimSpace(lastname)

	hash := getHash(firstname, lastname, countryId)


	if coach, ok := db.Coaches[hash]; ok {

		return coach.Id
	}



	id := len(db.Coaches) + 1
	db.Coaches[hash] = coach{Id: id, Firstname: firstname, Lastname: lastname, Country: countryId}


	return id
}

func addCountry(name string) int {
	
	db.dbLock.Lock()
	defer db.dbLock.Unlock()
	
	
	name = strings.TrimSpace(name)

	hash := getHash(name)


	if country, ok := db.Countries[hash]; ok {

		return country.Id
	}



	id := len(db.Countries) + 1
	db.Countries[hash] = country{Id: id, Name: name}


	return id
}

func addMatch(teamA, teamB, season, referee, date, score int) int {
	
	db.dbLock.Lock()
	defer db.dbLock.Unlock()
	
	
	hash := getHash(teamA, teamB, season, date)


	if match, ok := db.Matches[hash]; ok {

		return match.Id
	}



	id := len(db.Matches) + 1
	db.Matches[hash] = match{Id: id, TeamA: teamA, TeamB: teamB, Season: season,
							Referee: referee, Date: date, Score: score}



	return id
}

func addScore(teamA, teamB int) int {
	
	db.dbLock.Lock()
	defer db.dbLock.Unlock()
	
	
	hash := getHash(teamA, teamB)


	if score, ok := db.Scores[hash]; ok {

		return score.Id
	}



	id := len(db.Scores) + 1
	db.Scores[hash] = score{Id: id, TeamA: teamA, TeamB: teamB}



	return id
}

func addPlaysMatchInTeam(playerId, number, teamId, matchId int) int {
	
	db.dbLock.Lock()
	defer db.dbLock.Unlock()
	
	
	hash := getHash(playerId, teamId, matchId)


	if pmit, ok := db.PlaysMatchInTeams[hash]; ok {

		return pmit.Id
	}



	id := len(db.PlaysMatchInTeams) + 1
	db.PlaysMatchInTeams[hash] = playsMatchInTeam{Id: id, PlayerId: playerId, Number: number,
								TeamId: teamId, MatchId: matchId}



	return id
}

func addPlaysIn(teamId, playerId int) int {
	
	db.dbLock.Lock()
	defer db.dbLock.Unlock()
	
	
	hash := getHash(teamId, playerId)


	if pin, ok := db.PlaysIn[hash]; ok {

		return pin.Id
	}



	id := len(db.PlaysIn) + 1
	db.PlaysIn[hash] = playsIn{Id: id, TeamId: teamId, PlayerId: playerId}



	return id
}

func addCoacheses(coachId, teamId, matchId int) int {
	
	db.dbLock.Lock()
	defer db.dbLock.Unlock()
	
	
	hash := getHash(coachId, teamId, matchId)


	if coacheses, ok := db.Coacheses[hash]; ok {

		return coacheses.Id
	}



	id := len(db.Coacheses) + 1
	db.Coacheses[hash] = coacheses{Id: id, CoachId: coachId, TeamId: teamId, MatchId: matchId}



	return id
}

func addTeam(name string, countryId int) int {
	
	db.dbLock.Lock()
	defer db.dbLock.Unlock()
	
	hash := getHash(name, countryId)


	if team, ok := db.Teams[hash]; ok {

		return team.Id
	}



	id := len(db.Teams) + 1
	db.Teams[hash] = team{Id: id, Name: name, CountryId: countryId}



	return id
}

func addSeason(name string, competitionId int) int {
	
	db.dbLock.Lock()
	defer db.dbLock.Unlock()
	
	hash := getHash(name, competitionId)


	if season, ok := db.Seasons[hash]; ok {

		return season.Id
	}



	id := len(db.Seasons) + 1
	db.Seasons[hash] = season{Id: id, Name: name, CompetitionId: competitionId}



	return id
}

func addCompetition(name string) int {
	
	db.dbLock.Lock()
	defer db.dbLock.Unlock()
	
	hash := getHash(name)


	if competition, ok := db.Competitions[hash]; ok {

		return competition.Id
	}



	id := len(db.Competitions) + 1
	db.Competitions[hash] = competition{Id: id, Name: name}



	return id
}


func addGoal(playerId, matchId, time int) int {
	
	db.dbLock.Lock()
	defer db.dbLock.Unlock()
	
	
	hash := getHash(playerId, matchId, time)


	if goal, ok := db.Goals[hash]; ok {

		return goal.Id
	}



	id := len(db.Goals) + 1
	db.Goals[hash] = goal{Id: id, PlayerId: playerId, MatchId: matchId, Time: time}



	return id
}



func addCard(playerId, matchId, time, cardType int) int {

	db.dbLock.Lock()
	defer db.dbLock.Unlock()


	hash := getHash(playerId, matchId, time, cardType)


	if card, ok := db.Cards[hash]; ok {

		return card.Id
	}



	id := len(db.Cards) + 1
	db.Cards[hash] = card{Id: id, PlayerId: playerId, MatchId: matchId, Time: time, Type: cardType}



	return id
}





func getHash(params ...interface{}) string {
	encoding := ""
	for _, param := range params {
		switch param.(type) {
			case int:
				encoding += base32.StdEncoding.EncodeToString([]byte(strconv.Itoa(param.(int))))
				encoding += ":"

			case string:
				encoding += base32.StdEncoding.EncodeToString([]byte(strings.TrimSpace(strings.ToLower(param.(string)))))
				encoding += ":"
		}
	}

	checksum := sha256.Sum256([]byte(encoding))
	hash := hex.EncodeToString(checksum[:])
	return hash
}


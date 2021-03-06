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
	Url string
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
	FinalType string
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
	
	Urls map[string]int
	
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
	Urls: map[string]int{},
}


/*
Write the database to a file
*/
func writeDb(file string) error {
	json, err := json.MarshalIndent(db, "", "    ") 
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, json, 0777)
}


/*
Retrieve the database size
Returns a struct containing the number of each items
*/
func getDbSize() map[string]int {

	db.dbLock.Lock()
	defer db.dbLock.Unlock()

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


/*
Fetch an url from the cache
*/
func getUrlFromCache(url string) (int, bool) {
	db.dbLock.RLock()
	defer db.dbLock.RUnlock()

	id, ok := db.Urls[url]
	return id, ok
}


/*
Save an url int the cache
Should be called from a thread safe caller.
*/
func addUrlToCache(url string, id int) {
	db.Urls[url] = id
}


/*
Add a player record to the database
*/
func addPlayer(firstname, lastname string, countryId, dateOfBirth, height, weight int, position, url string) int {
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

	addUrlToCache(url, id);

	return id
}


/*
Add a referee record to the database
*/
func addReferee(firstname, lastname string, countryId int, url string) int {
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

	addUrlToCache(url, id);

	return id
}


/*
Add a coach record to the database
*/
func addCoach(firstname, lastname string, countryId int, url string) int {
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

	addUrlToCache(url, id);

	return id
}


/*
Add a country record to the database
*/
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


/*
Add a match record to the database
*/
func addMatch(teamA, teamB, season, referee, date, score int, url string, finalType string) int {
	
	db.dbLock.Lock()
	defer db.dbLock.Unlock()
	
	
	hash := getHash(teamA, teamB, season, date)


	if match, ok := db.Matches[hash]; ok {

		return match.Id
	}



	id := len(db.Matches) + 1
	db.Matches[hash] = match{Id: id, TeamA: teamA, TeamB: teamB, Season: season,
							Referee: referee, Date: date, Score: score, FinalType: finalType}

	addUrlToCache(url, id);

	return id
}


/*
Add a score record to the database
*/
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


/*
Add a playsmatchinteam record to the database
*/
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


/*
Add a playsIn record to the database
*/
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


/*
Add a coacheses record to the database
*/
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


/*
Add a team record to the database
*/
func addTeam(name string, countryId int, url string) int {
	
	db.dbLock.Lock()
	defer db.dbLock.Unlock()
	
	hash := getHash(name, countryId)


	if team, ok := db.Teams[hash]; ok {

		return team.Id
	}



	id := len(db.Teams) + 1
	db.Teams[hash] = team{Id: id, Name: name, CountryId: countryId}

	addUrlToCache(url, id);

	return id
}


/*
Add a season record to the database
*/
func addSeason(name string, competitionId int, url string) int {
	
	db.dbLock.Lock()
	defer db.dbLock.Unlock()
	
	hash := getHash(name, competitionId)


	if season, ok := db.Seasons[hash]; ok {

		return season.Id
	}



	id := len(db.Seasons) + 1
	db.Seasons[hash] = season{Id: id, Name: name, CompetitionId: competitionId}

	addUrlToCache(url, id);

	return id
}


/*
Add a competition record to the database
*/
func addCompetition(name string, url string) int {
	
	db.dbLock.Lock()
	defer db.dbLock.Unlock()
	
	hash := getHash(name)


	if competition, ok := db.Competitions[hash]; ok {

		return competition.Id
	}



	id := len(db.Competitions) + 1
	db.Competitions[hash] = competition{Id: id, Name: name}

	addUrlToCache(url, id);

	return id
}


/*
Add a goal record to the database
*/
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


/*
Add a card record to the database
*/
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


/*
Hash a set of parameters
*/
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


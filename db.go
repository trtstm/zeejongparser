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

var db = struct {
	playersLock sync.RWMutex
	Players     map[string]player

	refereesLock sync.RWMutex
	Referees     map[string]referee

	coachesLock sync.RWMutex
	Coaches     map[string]coach

	countriesLock sync.RWMutex
	Countries     map[string]country

	matchesLock sync.RWMutex
	Matches     map[string]match

	scoresLock sync.RWMutex
	Scores     map[string]score

	playsMatchInTeamsLock sync.RWMutex
	PlaysMatchInTeams     map[string]playsMatchInTeam

	coachesesLock sync.RWMutex
	Coacheses     map[string]coacheses

	teamsLock sync.RWMutex
	Teams     map[string]team

	seasonsLock sync.RWMutex
	Seasons     map[string]season

	competitionsLock sync.RWMutex
	Competitions     map[string]competition
}{
	Players: map[string]player{},
	Referees: map[string]referee{},
	Coaches: map[string]coach{},
	Countries: map[string]country{},
	Matches: map[string]match{},
	Scores: map[string]score{},
	PlaysMatchInTeams: map[string]playsMatchInTeam{},
	Coacheses: map[string]coacheses{},
	Teams: map[string]team{},
	Seasons: map[string]season{},
	Competitions: map[string]competition{},
}

func writeDb(file string) error {
	json, err := json.Marshal(db)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, json, 0777)
}

func getDbSize() map[string]int {
	db.playersLock.Lock()
	defer db.playersLock.Unlock()
	db.refereesLock.Lock()
	defer db.refereesLock.Unlock()
	db.coachesLock.Lock()
	defer db.coachesLock.Unlock()
	db.countriesLock.Lock()
	defer db.countriesLock.Unlock()
	db.matchesLock.Lock()
	defer db.matchesLock.Unlock()
	db.scoresLock.Lock()
	defer db.scoresLock.Unlock()
	db.playsMatchInTeamsLock.Lock()
	defer db.playsMatchInTeamsLock.Unlock()
	db.coachesesLock.Lock()
	defer db.coachesesLock.Unlock()
	db.teamsLock.Lock()
	defer db.teamsLock.Unlock()
	db.seasonsLock.Lock()
	defer db.seasonsLock.Unlock()
	db.competitionsLock.Lock()
	defer db.competitionsLock.Unlock()

	return map[string]int {
		"Players": len(db.Players),
		"Referees": len(db.Referees),
		"Coaches": len(db.Coaches),
		"Countries": len(db.Countries),
		"Matches": len(db.Matches),
		"Scores": len(db.Scores),
		"PlaysMatchInTeam": len(db.PlaysMatchInTeams),
		"Coacheses": len(db.Coacheses),
		"Teams": len(db.Teams),
		"Seasons": len(db.Seasons),
		"Competitions": len(db.Competitions),
	}
}

func addPlayer(firstname, lastname string, countryId, dateOfBirth, height, weight int, position string) int {
	firstname = strings.TrimSpace(firstname)
	lastname = strings.TrimSpace(lastname)
	position = strings.TrimSpace(position)

	hash := getHash(firstname, lastname, countryId, dateOfBirth)

	db.playersLock.RLock()
	if player, ok := db.Players[hash]; ok {
		db.playersLock.RUnlock()
		return player.Id
	}
	db.playersLock.RUnlock()

	db.playersLock.Lock()
	id := len(db.Players)
	db.Players[hash] = player{Id: id, Firstname: firstname,
		Lastname: lastname, Country: countryId, DateOfBirth: dateOfBirth,
		Height: height, Weight: weight, Position: position}
	db.playersLock.Unlock()

	return id
}

func addReferee(firstname, lastname string, countryId int) int {
	firstname = strings.TrimSpace(firstname)
	lastname = strings.TrimSpace(lastname)

	hash := getHash(firstname, lastname, countryId)

	db.refereesLock.RLock()
	if referee, ok := db.Referees[hash]; ok {
		db.refereesLock.RUnlock()
		return referee.Id
	}
	db.refereesLock.RUnlock()

	db.refereesLock.Lock()
	id := len(db.Referees)
	db.Referees[hash] = referee{Id: id, Firstname: firstname, Lastname: lastname, Country: countryId}
	db.refereesLock.Unlock()

	return id
}

func addCoach(firstname, lastname string, countryId int) int {
	firstname = strings.TrimSpace(firstname)
	lastname = strings.TrimSpace(lastname)

	hash := getHash(firstname, lastname, countryId)

	db.coachesLock.RLock()
	if coach, ok := db.Coaches[hash]; ok {
		db.coachesLock.RUnlock()
		return coach.Id
	}
	db.coachesLock.RUnlock()

	db.coachesLock.Lock()
	id := len(db.Coaches)
	db.Coaches[hash] = coach{Id: id, Firstname: firstname, Lastname: lastname, Country: countryId}
	db.coachesLock.Unlock()

	return id
}

func addCountry(name string) int {
	name = strings.TrimSpace(name)

	hash := getHash(name)

	db.countriesLock.RLock()
	if country, ok := db.Countries[hash]; ok {
		db.countriesLock.RUnlock()
		return country.Id
	}
	db.countriesLock.RUnlock()

	db.countriesLock.Lock()
	id := len(db.Countries)
	db.Countries[hash] = country{Id: id, Name: name}
	db.countriesLock.Unlock()

	return id
}

func addMatch(teamA, teamB, season, referee, date, score int) int {
	hash := getHash(teamA, teamB, season, date)

	db.matchesLock.RLock()
	if match, ok := db.Matches[hash]; ok {
		db.matchesLock.RUnlock()
		return match.Id
	}
	db.matchesLock.RUnlock()

	db.matchesLock.Lock()
	id := len(db.Matches)
	db.Matches[hash] = match{Id: id, TeamA: teamA, TeamB: teamB, Season: season,
							Referee: referee, Date: date, Score: score}

	db.matchesLock.Unlock()

	return id
}

func addScore(teamA, teamB int) int {
	hash := getHash(teamA, teamB)

	db.scoresLock.RLock()
	if score, ok := db.Scores[hash]; ok {
		db.scoresLock.RUnlock()
		return score.Id
	}
	db.scoresLock.RUnlock()

	db.scoresLock.Lock()
	id := len(db.Scores)
	db.Scores[hash] = score{Id: id, TeamA: teamA, TeamB: teamB}

	db.scoresLock.Unlock()

	return id
}

func addPlaysMatchInTeam(playerId, number, teamId, matchId int) int {
	hash := getHash(playerId, teamId, matchId)

	db.playsMatchInTeamsLock.RLock()
	if pmit, ok := db.PlaysMatchInTeams[hash]; ok {
		db.playsMatchInTeamsLock.RUnlock()
		return pmit.Id
	}
	db.playsMatchInTeamsLock.RUnlock()

	db.playsMatchInTeamsLock.Lock()
	id := len(db.PlaysMatchInTeams)
	db.PlaysMatchInTeams[hash] = playsMatchInTeam{Id: id, PlayerId: playerId, Number: number,
								TeamId: teamId, MatchId: matchId}

	db.playsMatchInTeamsLock.Unlock()

	return id
}

func addCoacheses(coachId, teamId, matchId int) int {
	hash := getHash(coachId, teamId, matchId)

	db.coachesesLock.RLock()
	if coacheses, ok := db.Coacheses[hash]; ok {
		db.coachesesLock.RUnlock()
		return coacheses.Id
	}
	db.coachesesLock.RUnlock()

	db.coachesesLock.Lock()
	id := len(db.Coacheses)
	db.Coacheses[hash] = coacheses{Id: id, CoachId: coachId, TeamId: teamId, MatchId: matchId}

	db.coachesesLock.Unlock()

	return id
}

func addTeam(name string, countryId int) int {
	hash := getHash(name, countryId)

	db.teamsLock.RLock()
	if team, ok := db.Teams[hash]; ok {
		db.teamsLock.RUnlock()
		return team.Id
	}
	db.teamsLock.RUnlock()

	db.teamsLock.Lock()
	id := len(db.Teams)
	db.Teams[hash] = team{Id: id, Name: name, CountryId: countryId}

	db.teamsLock.Unlock()

	return id
}

func addSeason(name string, competitionId int) int {
	hash := getHash(name, competitionId)

	db.seasonsLock.RLock()
	if season, ok := db.Seasons[hash]; ok {
		db.seasonsLock.RUnlock()
		return season.Id
	}
	db.seasonsLock.RUnlock()

	db.seasonsLock.Lock()
	id := len(db.Seasons)
	db.Seasons[hash] = season{Id: id, Name: name, CompetitionId: competitionId}

	db.seasonsLock.Unlock()

	return id
}

func addCompetition(name string) int {
	hash := getHash(name)

	db.competitionsLock.RLock()
	if competition, ok := db.Competitions[hash]; ok {
		db.competitionsLock.RUnlock()
		return competition.Id
	}
	db.competitionsLock.RUnlock()

	db.competitionsLock.Lock()
	id := len(db.Competitions)
	db.Competitions[hash] = competition{Id: id, Name: name}

	db.competitionsLock.Unlock()

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


package main

import (
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"

	"strconv"
	"strings"
	"sync"
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
}{
	Players: map[string]player{},
	Referees: map[string]referee{},
	Coaches: map[string]coach{},
	Countries: map[string]country{},
	Matches: map[string]match{},
}

func addPlayer(firstname, lastname string, countryId, dateOfBirth, height, weight int, position string) int {
	firstname = strings.TrimSpace(firstname)
	lastname = strings.TrimSpace(lastname)
	position = strings.TrimSpace(position)

	encoding := base32.StdEncoding.EncodeToString([]byte(firstname))
	encoding += ":"
	encoding += base32.StdEncoding.EncodeToString([]byte(lastname))
	encoding += ":"
	encoding += base32.StdEncoding.EncodeToString([]byte(strconv.Itoa(countryId)))
	encoding += ":"
	encoding += base32.StdEncoding.EncodeToString([]byte(strconv.Itoa(dateOfBirth)))

	checksum := sha256.Sum256([]byte(encoding))
	hash := hex.EncodeToString(checksum[:])

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

	encoding := base32.StdEncoding.EncodeToString([]byte(firstname))
	encoding += ":"
	encoding += base32.StdEncoding.EncodeToString([]byte(lastname))
	encoding += ":"
	encoding += base32.StdEncoding.EncodeToString([]byte(strconv.Itoa(countryId)))

	checksum := sha256.Sum256([]byte(encoding))
	hash := hex.EncodeToString(checksum[:])

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

	encoding := base32.StdEncoding.EncodeToString([]byte(firstname))
	encoding += ":"
	encoding += base32.StdEncoding.EncodeToString([]byte(lastname))
	encoding += ":"
	encoding += base32.StdEncoding.EncodeToString([]byte(strconv.Itoa(countryId)))

	checksum := sha256.Sum256([]byte(encoding))
	hash := hex.EncodeToString(checksum[:])

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

	encoding := base32.StdEncoding.EncodeToString([]byte(strings.ToLower(name)))

	checksum := sha256.Sum256([]byte(encoding))
	hash := hex.EncodeToString(checksum[:])

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
	encoding := base32.StdEncoding.EncodeToString([]byte(strconv.Itoa(teamA)))
	encoding += ":"
	encoding += base32.StdEncoding.EncodeToString([]byte(strconv.Itoa(teamB)))
	encoding += ":"
	encoding += base32.StdEncoding.EncodeToString([]byte(strconv.Itoa(season)))
	encoding += ":"
	encoding += base32.StdEncoding.EncodeToString([]byte(strconv.Itoa(date)))

	checksum := sha256.Sum256([]byte(encoding))
	hash := hex.EncodeToString(checksum[:])

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


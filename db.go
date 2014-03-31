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

var db = struct {
	playersLock sync.RWMutex
	players     map[string]player

	refereesLock sync.RWMutex
	referees     map[string]referee

	coachesLock sync.RWMutex
	coaches     map[string]coach

	countriesLock sync.RWMutex
	countries     map[string]country
}{
	players: map[string]player{},
	referees: map[string]referee{},
	coaches: map[string]coach{},
	countries: map[string]country{},
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
	if player, ok := db.players[hash]; ok {
		db.playersLock.RUnlock()
		return player.Id
	}
	db.playersLock.RUnlock()

	db.playersLock.Lock()
	id := len(db.players)
	db.players[hash] = player{Id: id, Firstname: firstname,
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
	if referee, ok := db.referees[hash]; ok {
		db.refereesLock.RUnlock()
		return referee.Id
	}
	db.refereesLock.RUnlock()

	db.refereesLock.Lock()
	id := len(db.referees)
	db.referees[hash] = referee{Id: id, Firstname: firstname, Lastname: lastname, Country: countryId}
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
	if coach, ok := db.coaches[hash]; ok {
		db.coachesLock.RUnlock()
		return coach.Id
	}
	db.coachesLock.RUnlock()

	db.coachesLock.Lock()
	id := len(db.coaches)
	db.coaches[hash] = coach{Id: id, Firstname: firstname, Lastname: lastname, Country: countryId}
	db.coachesLock.Unlock()

	return id
}

func addCountry(name string) int {
	name = strings.TrimSpace(name)

	encoding := base32.StdEncoding.EncodeToString([]byte(strings.ToLower(name)))

	checksum := sha256.Sum256([]byte(encoding))
	hash := hex.EncodeToString(checksum[:])

	db.countriesLock.RLock()
	if country, ok := db.countries[hash]; ok {
		db.countriesLock.RUnlock()
		return country.Id
	}
	db.countriesLock.RUnlock()

	db.countriesLock.Lock()
	id := len(db.countries)
	db.countries[hash] = country{Id: id, Name: name}
	db.countriesLock.Unlock()

	return id
}

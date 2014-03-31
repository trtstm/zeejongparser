package main

import (
	"sync"
	"strings"
	"encoding/base32"
	"crypto/sha256"
	"strconv"
	"encoding/hex"
	"log"
)

type player struct {
	Id int
	Firstname string
	Lastname string
	Country int
	DateOfBirth int
	Height int
	Weight int
	Position string
}

var db = struct {
	playersLock sync.RWMutex
	players map[string]player
	
}{
	players: map[string]player{},
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
	log.Println(db.players[hash])
	db.playersLock.Unlock()
	
	return id
}

package main

import (
	"github.com/PuerkitoBio/goquery"
	"sync"
	"encoding/hex"
	"crypto/md5"
	"net/http"
	"time"
	"errors"
	"os"
	"io/ioutil"
	"strings"
)

type HitRate struct {
	Accesses int
	Hits int
}

var hitRateLock sync.RWMutex
var hitRate HitRate

func getHitRate() HitRate {
	return hitRate
}

var diskLock sync.RWMutex

func getFromDisk(hash string) (*strings.Reader, error) {
	diskLock.RLock()
	defer diskLock.RUnlock()

	if _, err := os.Stat("cache/" + hash + ".html"); os.IsNotExist(err) {
		return nil, errors.New(hash + " not found in cache")
	}

	bytes, err := ioutil.ReadFile("cache/" + hash + ".html")
	if err != nil {
		return nil, errors.New("could not read " + hash + " from cache")
	}

	return strings.NewReader(string(bytes)), nil
}

func addToDisk(hash, contents string) error {
	diskLock.Lock()
	defer diskLock.Unlock()

	if _, err := os.Stat("cache/" + hash + ".html"); !os.IsNotExist(err) {
		return errors.New(hash + " is already in cache")
	}

	err := ioutil.WriteFile("cache/" + hash + ".html", []byte(contents), 0755)
	if err != nil {
		return errors.New("could not write to cache " + hash)
	}

	return nil
}

var documentCacheLock sync.RWMutex
var documentCache = map[string]*goquery.Document{}

func getDocument(url string) (*goquery.Document, error) {
	hitRateLock.Lock()
	hitRate.Accesses += 1
	hitRateLock.Unlock()
	
	bHash := md5.Sum([]byte(url))
	hash := hex.EncodeToString(bHash[:])

	documentCacheLock.RLock()
	document, ok := documentCache[hash]
	if ok {
		documentCacheLock.RUnlock()

		hitRateLock.Lock()
		hitRate.Hits += 1
		hitRateLock.Unlock()
		return document, nil
	}
	documentCacheLock.RUnlock()

	reader, err := getFromDisk(hash)
	if err != nil {
		var resp *http.Response
		for i := 0; i < 10; i++ {
			resp, err = http.Get(url)
			if err != nil {
				return document, err
			}
			defer resp.Body.Close()

			if resp.StatusCode == 200 {
				break
			} else if resp.StatusCode == 500 {
				time.Sleep(time.Second * 10)
			} else {
				return document, errors.New("response status is " + resp.Status)
			}
		}

		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return document, errors.New("could not read response")
		}

		addToDisk(hash, string(contents))
		reader = strings.NewReader(string(contents))
	}

	document, err = goquery.NewDocumentFromReader(reader)
	if err != nil {
		return document, err
	}

	documentCacheLock.Lock()
	documentCache[hash] = document
	documentCacheLock.Unlock()

	return document, nil
}

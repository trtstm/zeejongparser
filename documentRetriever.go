package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type CacheInfo struct {
	Accesses     int
	Items        int
	DiskAccesses int
	MemAccesses  int
	UrlAccesses  int
}

var cacheInfoLock sync.RWMutex
var cacheInfo CacheInfo

func getCacheInfo() CacheInfo {
	documentCacheLock.RLock()
	items := len(documentCache)
	documentCacheLock.RUnlock()

	cacheInfoLock.Lock()
	cacheInfo.Items = items
	cacheInfoLock.Unlock()

	return cacheInfo
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

	err := ioutil.WriteFile("cache/"+hash+".html", []byte(contents), 0755)
	if err != nil {
		return errors.New("could not write to cache " + hash)
	}

	return nil
}

var documentCacheLock sync.RWMutex
var documentCache = map[string]*goquery.Document{}

func getDocument(url string) (*goquery.Document, error) {
	documentCacheLock.RLock()
	if len(documentCache) > 200 {
		documentCacheLock.RUnlock()
		documentCacheLock.Lock()
		documentCache = map[string]*goquery.Document{}
		documentCacheLock.Unlock()

		runtime.GC()
	} else {
		documentCacheLock.RUnlock()
	}

	cacheInfoLock.Lock()
	cacheInfo.Accesses += 1
	cacheInfoLock.Unlock()

	bHash := md5.Sum([]byte(url))
	hash := hex.EncodeToString(bHash[:])

	documentCacheLock.RLock()
	document, ok := documentCache[hash]
	if ok {
		documentCacheLock.RUnlock()

		cacheInfoLock.Lock()
		cacheInfo.MemAccesses += 1
		cacheInfoLock.Unlock()
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

			if resp.StatusCode == 200 {
				cacheInfoLock.Lock()
				cacheInfo.UrlAccesses += 1
				cacheInfoLock.Unlock()
				break
			} else if resp.StatusCode/100 == 5 && i < 9 {
				resp.Body.Close()
				time.Sleep(time.Second * 10)
			} else {
				resp.Body.Close()
				return document, errors.New("response status is " + resp.Status)
			}
		}

		defer resp.Body.Close()

		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return document, errors.New("could not read response")
		}

		addToDisk(hash, string(contents))
		reader = strings.NewReader(string(contents))
	} else {
		cacheInfoLock.Lock()
		cacheInfo.DiskAccesses += 1
		cacheInfoLock.Unlock()
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

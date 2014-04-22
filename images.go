package main

import (
	"sync"
	"io/ioutil"
	"os"
	"net/http"
	"time"
	"strconv"
	"strings"
)

var imageLock = sync.RWMutex{}

func addImage(name string, data []byte) error {
	return ioutil.WriteFile("images/" + name, data, 0755)
}

func hasImage(name string) bool {
	if _, err := os.Stat("images/" + name); os.IsNotExist(err) {
		return false
	}

	return true
}

func getImage(url, ownerType string, id int) string {
	imageName := ownerType + "-" + strconv.Itoa(id) + ".png"
	genericName := ownerType + "-" + "generic.png"

	if strings.HasSuffix(url, "generic.png") {
		return genericName
	}

	imageLock.RLock()
	if hasImage(imageName) {
		imageLock.RUnlock()
		return imageName
	}
	imageLock.RUnlock()

	var resp *http.Response
	for i := 0; i < 10; i++ {
		var err error
		resp, err = cThrottler.Get(url)
		if err != nil {
			return genericName
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			break
		} else if resp.StatusCode/100 == 5 {
			time.Sleep(time.Second * 10)
		} else {
			return genericName
		}
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return genericName
	}

	imageLock.Lock()
	defer imageLock.Unlock()

	if hasImage(imageName) {
		return imageName
	}

	if addImage(imageName, contents) != nil {
		return genericName
	}

	return imageName
}

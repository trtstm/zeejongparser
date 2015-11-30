package main

import (
	"runtime"
	"sync"
	"flag"
	"log"
	"io/ioutil"
	"encoding/json"
)

const BASE = "http://SOME_BASE_URL.com"

var competitions = []map[string]string{}

func init() {
	competitions = append(competitions, map[string]string{"name": "World Cup", "archive": "http://SOME_WORLD_CUP_ARCHIVE.com"})

	competitions = append(competitions, map[string]string{"name": "European Championship", "archive": "http://SOME_EC_ARCHIVE.com"})
}

func main() {
	
	//Parse command line flags
	filename := flag.String("file", "output.json", "Filename for output")
	flag.Parse()
	
	jsonContents, err := ioutil.ReadFile(*filename)
	if err == nil {
		err = json.Unmarshal(jsonContents, &db)
		if err != nil {
			log.Println("Could not load json")
		} else {
			log.Println("Loaded json")
		}
	}
	
	//Run the parser
	runtime.GOMAXPROCS(runtime.NumCPU())

	go startWs()

	wg := sync.WaitGroup{}

	for _, competition := range competitions {
		wg.Add(1)
		go func(url, name string) {
			defer wg.Done()

			parseCompetition(url, name)
			
			_ = name
			
		}(competition["archive"], competition["name"])
	}

	wg.Wait()
	
	
	//Save output to json file
	err = writeDb(*filename)
	if(err != nil) {
		log.Println("Could not write file:", err)
	} else {
		log.Println("Output written to", *filename)
	}
}

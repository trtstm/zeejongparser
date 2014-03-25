package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"html/template"
	"os"
	"log"
	"fmt"
	"encoding/json"
)

var rootTpl *template.Template

func init() {
	var err error
	rootTpl, err = template.ParseFiles("templates/root.html")
	if err != nil {
		panic("Could not parse root.html")
	}
}

func root(w http.ResponseWriter, r *http.Request) {
	rootTpl.Execute(w, nil)
}

func cacheInfo(w http.ResponseWriter, r *http.Request) {
	json, err := json.Marshal(getCacheInfo())
	if err != nil {
		fmt.Fprintf(w, "{}")
		return
	}

	fmt.Fprintf(w, string(json))
}

func startWs() {
	address := ":8080"
	if len(os.Args) == 2 {
		address = os.Args[1]
	}

	r := mux.NewRouter()
	r.HandleFunc("/", root)
	r.HandleFunc("/cacheInfo", cacheInfo)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	server := http.Server{Addr: address, Handler: r}
	err := server.ListenAndServe()

	if err != nil {
		log.Println(err)
	}
}

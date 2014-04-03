package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
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

func cacheInfoHandler(w http.ResponseWriter, r *http.Request) {
	json, err := json.Marshal(getCacheInfo())
	if err != nil {
		fmt.Fprintf(w, "{}")
		return
	}

	fmt.Fprintf(w, string(json))
}

func dbInfoHandler(w http.ResponseWriter, r *http.Request) {
	json, err := json.Marshal(getDbSize())
	if err != nil {
		fmt.Fprintf(w, "{}")
		return
	}

	fmt.Fprintf(w, string(json))
}

func startWs() {
	address := ":8080"

	r := mux.NewRouter()
	r.HandleFunc("/", root)
	r.HandleFunc("/cacheInfo", cacheInfoHandler)
	r.HandleFunc("/dbInfo", dbInfoHandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	server := http.Server{Addr: address, Handler: r}
	err := server.ListenAndServe()

	if err != nil {
		log.Println(err)
	}
}

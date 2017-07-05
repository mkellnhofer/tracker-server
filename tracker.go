package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func handleLoc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		handleOptionsLoc(w, r)
	} else if r.Method == "GET" {
		hanldeGetLoc(w, r)
	} else if r.Method == "POST" {
		handlePostLoc(w, r)
	} else {
		log.Printf("Invalid method! (Method: '%s')", r.Method)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.Error(w, "Bad request! (Invalid method.)", http.StatusBadRequest)
	}
}

func handleOptionsLoc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func hanldeGetLoc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	locs, err := getLocations()
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error! (Error while reading locations.)",
			http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(locs)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error! (Error while serializing data.)",
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func handlePostLoc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var loc location

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loc)
	if err != nil {
		log.Printf("Invalid JSON! ('%s')", err)
		http.Error(w, "Bad request! (Invalid JSON)", http.StatusBadRequest)
		return
	}

	err = addLocation(&loc)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error! (Error while adding location.)",
			http.StatusInternalServerError)
	}
}

func main() {
	conf := loadConfig()

	log.Print("Starting Tracker server 1.0.0-alpha.")

	http.HandleFunc("/loc", handleLoc)

	log.Printf("Listen on port '%d'.", conf.port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", conf.port), nil)
	if err != nil {
		log.Fatalf("Could not start server! (Error: %s)", err)
	}
}

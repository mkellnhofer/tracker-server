package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	aModel "kellnhofer.com/tracker/api/model"
	"kellnhofer.com/tracker/config"
	"kellnhofer.com/tracker/constant"
	"kellnhofer.com/tracker/data"
	lModel "kellnhofer.com/tracker/model"
	"kellnhofer.com/tracker/repo"
)

var locRepo *repo.LocationRepo

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

	lLocs, err := locRepo.GetLocations()
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error! (Error while reading locations.)",
			http.StatusInternalServerError)
		return
	}

	aLocs := toApiLocs(lLocs)

	json, err := json.Marshal(aLocs)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error! (Error while serializing data.)",
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func toApiLocs(iLocs []*lModel.Location) []*aModel.Location {
	var oLocs []*aModel.Location
	for _, iLoc := range iLocs {
		oLocs = append(oLocs, toApiLoc(iLoc))
	}
	return oLocs
}

func toApiLoc(iLoc *lModel.Location) *aModel.Location {
	return &aModel.Location{iLoc.Name, iLoc.Time, iLoc.Lat, iLoc.Lng}
}

func handlePostLoc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var aLoc aModel.Location

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&aLoc)
	if err != nil {
		log.Printf("Invalid JSON! ('%s')", err)
		http.Error(w, "Bad request! (Invalid JSON)", http.StatusBadRequest)
		return
	}

	lLoc := toLogicLoc(&aLoc)

	err = locRepo.AddLocation(lLoc)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error! (Error while adding location.)",
			http.StatusInternalServerError)
	}
}

func toLogicLoc(iLoc *aModel.Location) *lModel.Location {
	return &lModel.Location{0, iLoc.Name, iLoc.Time, iLoc.Lat, iLoc.Lng}
}

func main() {
	conf := config.LoadConfig()

	log.Printf("Starting Tracker server %s.", constant.AppVersion)

	db := data.GetDb()
	data.UpdateDb(db)

	locRepo = repo.NewLocationRepo(db)

	http.HandleFunc("/loc", handleLoc)

	log.Printf("Listen on port '%d'.", conf.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
	if err != nil {
		log.Fatalf("Could not start server! (Error: %s)", err)
	}
}

package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"kellnhofer.com/tracker/api/mapper"
	aModel "kellnhofer.com/tracker/api/model"
	"kellnhofer.com/tracker/repo"
)

type locationController struct {
	lRepo *repo.LocationRepo
}

func NewLocationController(lRepo *repo.LocationRepo) *locationController {
	return &locationController{lRepo}
}

// --- Public methods ---

func (c locationController) GetLocationsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		c.handleGetLocations(w, r)
	}
}

func (c locationController) CreateLocationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		c.handleCreateLocation(w, r)
	}
}

func (c locationController) GetLocationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		c.handleGetLocation(w, r)
	}
}

func (c locationController) DeleteLocationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		c.handleDeleteLocation(w, r)
	}
}

// --- Private methods ---

func (c locationController) handleGetLocations(w http.ResponseWriter, r *http.Request) {
	lLocs, err := c.lRepo.GetLocations()
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error! (Error while reading locations.)",
			http.StatusInternalServerError)
		return
	}

	aLocs := mapper.ToApiLocs(lLocs)

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

func (c locationController) handleCreateLocation(w http.ResponseWriter, r *http.Request) {
	var aLoc aModel.Location

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&aLoc)
	if err != nil {
		log.Printf("Invalid JSON! ('%s')", err)
		http.Error(w, "Bad request! (Invalid JSON)", http.StatusBadRequest)
		return
	}

	lLoc := mapper.ToLogicLoc(&aLoc)

	id, err := c.lRepo.AddLocation(lLoc)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error! (Error while adding location.)",
			http.StatusInternalServerError)
	}

	aLoc.Id = id

	json, err := json.Marshal(aLoc)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error! (Error while serializing data.)",
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (c locationController) handleGetLocation(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromPath(r.URL.Path)
	if err != nil {
		log.Printf("Invalid location ID!")
		http.Error(w, "Bad request! (Invalid location ID.)", http.StatusBadRequest)
		return
	}

	lLoc, err := c.lRepo.GetLocation(id)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error! (Error while reading location.)",
			http.StatusInternalServerError)
		return
	}
	if lLoc == nil {
		http.Error(w, "Not found! (Unknown location ID.)", http.StatusNotFound)
		return
	}

	aLoc := mapper.ToApiLoc(lLoc)

	json, err := json.Marshal(aLoc)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error! (Error while serializing data.)",
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (c locationController) handleDeleteLocation(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromPath(r.URL.Path)
	if err != nil {
		log.Printf("Invalid location ID!")
		http.Error(w, "Bad request! (Invalid location ID.)", http.StatusBadRequest)
		return
	}

	exists, err := c.lRepo.ExistsLocation(id)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error! (Error while deleting location.)",
			http.StatusInternalServerError)
	}
	if !exists {
		http.Error(w, "Not found! (Unknown location ID.)", http.StatusNotFound)
		return
	}

	err = c.lRepo.DeleteLocation(id)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error! (Error while deleting location.)",
			http.StatusInternalServerError)
		return
	}
}

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

func (c locationController) CreateGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		c.handleGet(w, r)
	}
}

func (c locationController) CreatePostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		c.handlePost(w, r)
	}
}

// --- Private methods ---

func (c locationController) handleGet(w http.ResponseWriter, r *http.Request) {
	log.Println("Handle get.")

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

func (c locationController) handlePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var aLoc aModel.Location

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&aLoc)
	if err != nil {
		log.Printf("Invalid JSON! ('%s')", err)
		http.Error(w, "Bad request! (Invalid JSON)", http.StatusBadRequest)
		return
	}

	lLoc := mapper.ToLogicLoc(&aLoc)

	err = c.lRepo.AddLocation(lLoc)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error! (Error while adding location.)",
			http.StatusInternalServerError)
	}
}

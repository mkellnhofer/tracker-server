package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"kellnhofer.com/tracker/api/controller"
	"kellnhofer.com/tracker/config"
	"kellnhofer.com/tracker/constant"
	"kellnhofer.com/tracker/data"
	"kellnhofer.com/tracker/repo"
)

func handleOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "HEAD, GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func main() {
	conf := config.LoadConfig()

	log.Printf("Starting Tracker server %s.", constant.AppVersion)

	// Open and update database
	db := data.GetDb()
	data.UpdateDb(db)

	// Create repos
	locRepo := repo.NewLocationRepo(db)

	// Create controllers
	locCtrl := controller.NewLocationController(locRepo)

	// Create router
	router := mux.NewRouter().StrictSlash(true)
	// Add public routes
	router.HandleFunc("/loc", handleOptions).Methods("OPTIONS")
	router.Handle("/loc", locCtrl.GetLocationsHandler()).Methods("GET")
	router.Handle("/loc", locCtrl.CreateLocationHandler()).Methods("POST")
	router.Handle("/loc/deleted", locCtrl.GetDeletedLocationIdsHandler()).Methods("GET")
	router.Handle("/loc/{id}", locCtrl.GetLocationHandler()).Methods("GET")
	router.Handle("/loc/{id}", locCtrl.DeleteLocationHandler()).Methods("DELETE")

	// Register router
	http.Handle("/", router)

	// Start HTTP server
	log.Printf("Listen on port '%d'.", conf.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
	if err != nil {
		log.Fatalf("Could not start server! (Error: %s)", err)
	}
}

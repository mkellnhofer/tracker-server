package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"kellnhofer.com/tracker/api/controller"
	"kellnhofer.com/tracker/config"
	"kellnhofer.com/tracker/constant"
	"kellnhofer.com/tracker/data"
	"kellnhofer.com/tracker/middleware"
	"kellnhofer.com/tracker/repo"
)

func handleOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "HEAD, GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func createRoute(route *negroni.Negroni, handler http.HandlerFunc) http.Handler {
	newRoute := route.With()
	newRoute.UseHandlerFunc(handler)
	return newRoute
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

	// Create middlewares
	authMidw := middleware.NewAuthMiddleware(conf)

	// Create protected middleware route
	proRoute := negroni.New()
	proRoute.UseFunc(authMidw.GetAuthHandlerFunc())

	// Create router
	router := mux.NewRouter().StrictSlash(true)
	// Add public routes
	router.HandleFunc("/loc", handleOptions).Methods("OPTIONS")
	// Add protected routes
	router.Handle("/loc", createRoute(proRoute, locCtrl.GetLocationsHandler())).Methods("GET")
	router.Handle("/loc", createRoute(proRoute, locCtrl.CreateLocationHandler())).Methods("POST")
	router.Handle("/loc/deleted", createRoute(proRoute, locCtrl.GetDeletedLocationIdsHandler())).
		Methods("GET")
	router.Handle("/loc/{id}", createRoute(proRoute, locCtrl.GetLocationHandler())).
		Methods("GET")
	router.Handle("/loc/{id}", createRoute(proRoute, locCtrl.ChangeLocationHandler())).
		Methods("PUT")
	router.Handle("/loc/{id}", createRoute(proRoute, locCtrl.DeleteLocationHandler())).
		Methods("DELETE")

	// Register router
	http.Handle("/", router)

	// Start HTTP server
	log.Printf("Listen on port '%d'.", conf.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
	if err != nil {
		log.Fatalf("Could not start server! (Error: %s)", err)
	}
}

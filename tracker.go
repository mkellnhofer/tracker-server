package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"kellnhofer.com/tracker/api/controller"
	"kellnhofer.com/tracker/config"
	"kellnhofer.com/tracker/constant"
	"kellnhofer.com/tracker/data"
	"kellnhofer.com/tracker/middleware"
	"kellnhofer.com/tracker/repo"
)

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

	// Create router
	router := mux.NewRouter().StrictSlash(true)
	// Create API sub route
	apiRoute := router.PathPrefix("/api/v1").Subrouter()
	// Add endpoints
	// GET /loc
	apiRoute.Methods("GET").
		Path("/loc").
		Handler(locCtrl.GetLocationsHandler())
	// GET /loc?change_time={change_time}
	apiRoute.Methods("GET").
		Path("/loc").
		Queries("change_time", "{change_time}").
		Handler(locCtrl.GetLocationsHandler())
	// POST /loc
	apiRoute.Methods("POST").
		Path("/loc").
		Handler(locCtrl.CreateLocationHandler())
	// GET /loc/deleted
	apiRoute.Methods("GET").
		Path("/loc/deleted").
		Handler(locCtrl.GetDeletedLocationIdsHandler())
	// GET /loc/deleted?deletion_time={deletion_time}
	apiRoute.Methods("GET").
		Path("/loc/deleted").
		Queries("deletion_time", "{deletion_time}").
		Handler(locCtrl.GetDeletedLocationIdsHandler())
	// GET /loc/{id}
	apiRoute.Methods("GET").
		Path("/loc/{id}").
		Handler(locCtrl.GetLocationHandler())
	// PUT /loc/{id}
	apiRoute.Methods("PUT").
		Path("/loc/{id}").
		Handler(locCtrl.ChangeLocationHandler())
	// DELETE /loc/{id}
	apiRoute.Methods("DELETE").
		Path("/loc/{id}").
		Handler(locCtrl.DeleteLocationHandler())

	// Create middlewares
	authMidw := middleware.NewAuthMiddleware(conf)
	corsMidw := cors.AllowAll()

	// Create middleware
	midw := negroni.New()
	midw.Use(corsMidw)
	midw.Use(authMidw)
	midw.UseHandler(router)

	// Register handler
	http.Handle("/", midw)

	// Start HTTP server
	log.Printf("Listen on port '%d'.", conf.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
	if err != nil {
		log.Fatalf("Could not start server! (Error: %s)", err)
	}
}

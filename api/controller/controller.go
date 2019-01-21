package controller

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func getId(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	v := vars["id"]
	return strconv.ParseInt(v, 10, 64)
}

func getChangeTime(r *http.Request) (int64, error) {
	v := r.FormValue("change_time")
	if v == "" {
		return int64(0), nil
	}
	return strconv.ParseInt(v, 10, 64)
}

func getDeletionTime(r *http.Request) (int64, error) {
	v := r.FormValue("deletion_time")
	if v == "" {
		return int64(0), nil
	}
	return strconv.ParseInt(v, 10, 64)
}

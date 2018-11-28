package controller

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

func getIdFromPath(path string) (int64, error) {
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 5 {
		return 0, errors.New("invalid path")
	}
	return strconv.ParseInt(pathParts[4], 10, 64)
}

func getChangeTimeFromQuery(query url.Values) (int64, error) {
	return getIntParamFromQuery(query, "change_time")
}

func getDeletionTimeFromQuery(query url.Values) (int64, error) {
	return getIntParamFromQuery(query, "deletion_time")
}

func getIntParamFromQuery(query url.Values, name string) (int64, error) {
	queryValue := query.Get(name)
	if queryValue == "" {
		return 0, nil
	}
	return strconv.ParseInt(queryValue, 10, 64)
}

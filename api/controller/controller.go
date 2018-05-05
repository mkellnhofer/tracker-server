package controller

import (
	"errors"
	"strconv"
	"strings"
)

func getIdFromPath(path string) (int64, error) {
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 3 {
		return 0, errors.New("invalid path")
	}
	return strconv.ParseInt(pathParts[2], 10, 64)
}

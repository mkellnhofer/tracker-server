package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func getLocations() ([]*location, error) {
	file, err := os.Open("locations.txt")
	if err != nil {
		log.Print("Locations file does not exist!")
		return []*location{}, nil
	}
	defer file.Close()

	var locs []*location

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		record := scanner.Text()
		loc, err := parseLocationRecord(record)
		if err != nil {
			return nil, err
		}
		locs = append(locs, loc)
	}
	if err := scanner.Err(); err != nil {
		log.Print(err)
		e := fmt.Sprintf("Error at reading location file! ('%s')", err)
		return nil, errors.New(e)
	}

	return locs, nil
}

func parseLocationRecord(record string) (*location, error) {
	data := strings.Split(record, ",")
	if len(data) < 3 {
		e := "Invalid location record!"
		return nil, errors.New(e)
	}

	lat, err := strconv.ParseFloat(data[0], 32)
	if err != nil {
		log.Print(err)
		e := fmt.Sprintf("Invalid latitude '%s'!", data[0])
		return nil, errors.New(e)
	}

	lng, err := strconv.ParseFloat(data[1], 32)
	if err != nil {
		log.Print(err)
		e := fmt.Sprintf("Invalid longitude '%s'!", data[1])
		return nil, errors.New(e)
	}

	time, err := time.Parse(time.RFC3339, data[2])
	if err != nil {
		log.Print(err)
		e := fmt.Sprintf("Invalid time '%s'!", data[2])
		return nil, errors.New(e)
	}

	return &location{float32(lat), float32(lng), time}, nil
}

func addLocation(loc *location) error {
	file, err := os.OpenFile("locations.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Print(err)
		e := fmt.Sprintf("Error at opening location file! ('%s')", err)
		return errors.New(e)
	}
	defer file.Close()

	record := fmt.Sprintf("%f,%f,%s\n", loc.Lat, loc.Lng, loc.Time.Format(time.RFC3339))

	_, err = file.WriteString(record)
	if err != nil {
		log.Print(err)
		e := fmt.Sprintf("Error at writing location file! ('%s')", err)
		return errors.New(e)
	}

	return nil
}

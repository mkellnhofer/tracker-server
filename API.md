# API Documentation

This document describes the REST API located under `/api/v1`.

## Authorization

The API access control is very simple. The password which is set in the config file has to be send in the Authorization header without any encoding.

## Endpoints

### Create Location

    POST /api/v1/loc

Request body:

    {
      "name": string,
      "time": datetime,
      "lat": float,
      "lng": float,
      "description": string,
      "persons": {
        "firstName": string,
        "lastName": string
      }
    }

Response body:

    {
      "id": integer,
      "changeTime": integer,
      "name": string,
      "time": datetime,
      "lat": float,
      "lng": float,
      "description": string,
      "persons": {
        "firstName": string,
        "lastName": string
      }
    }

### Update Location

    PUT /api/v1/loc/{id}

Request body:

    {
      "name": string,
      "time": datetime,
      "lat": float,
      "lng": float,
      "description": string,
      "persons": {
        "firstName": string,
        "lastName": string
      }
    }

Response body:

    {
      "id": integer,
      "changeTime": integer,
      "name": string,
      "time": datetime,
      "lat": float,
      "lng": float,
      "description": string,
      "persons": {
        "firstName": string,
        "lastName": string
      }
    }

### Delete Location

    DELETE /api/v1/loc/{id}

### Get Locations

    GET /api/v1/loc

Request Parameters:

- change_time (integer, optional): The earliest change time. 

Response body:

    [
      {
        "id": integer,
        "changeTime": integer,
        "name": string,
        "time": datetime,
        "lat": float,
        "lng": float,
        "description": string,
        "persons": {
          "firstName": string,
          "lastName": string
        }
      }
    ]

### Get Deleted Location IDs

    GET /api/v1/loc/deleted

Request parameters:

- deletion_time (integer, optional): The earliest deletion time.

Response body:

    [integer]